package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	MATCH  int = 1
	SCRIPT int = 2
)

const (
	LOG   int = 1
	ALERT int = 2
	BOTH  int = 3
)

type RuleNode struct {
	name    string
	action  int
	watch   string
	matches map[string]*regexp.Regexp
	script  string
	lang    string
	nodes   []*RuleNode
	trigger int
}

func RuleEngine(conf *ServerConfig, base *RuleNode, events <-chan *AuditEvent, scripts *ScriptPool) {

	conf.serverLog.Println("Rule Engine started")

	// the channel to dispatch events to workers over
	workers := make(chan *AuditEvent)

	go RuleEngineWorker(conf, base, workers, scripts)

	for ev := range events {
		workers <- ev
	}

	close(workers)

}

func RuleEngineWorker(conf *ServerConfig, base *RuleNode, events <-chan *AuditEvent, scripts *ScriptPool) {

	// so long as there are events to process
	for ev := range events {

		// recurively follow the descision tree
		RunNode(conf, base, ev, scripts)
	}
}

func RunNode(conf *ServerConfig, node *RuleNode, ev *AuditEvent, scripts *ScriptPool) {

	if node.action == SCRIPT {

		// get a script worker in the right language
		worker, err := scripts.GetWorker(node.lang)
		if err != nil {
			conf.serverLog.Println("Error retrieving worker: %s", err)
			return
		}

		// run the specified function
		match, err := worker.Run(node.script, ev)
		if err != nil {
			conf.serverLog.Println("Error running script %s-%s: %s", node.lang, node.script, err)
			return
		}

		// return if no match
		if !match {
			return
		}
	}

	if node.action == MATCH {

		// for each requested match
		for key, regex := range node.matches {

			// check that the key exists
			data, ok := (*ev)[key]
			if !ok {
				return
			}

			// check the match
			if !regex.MatchString(data) {
				return
			}
		}
	}

	// Arriving here indicates no conditions failed

	// run any triggers
	switch node.trigger {
	case LOG:
		conf.eventLog.Printf("[%s] - %s\n", node.name, *ev)
	case ALERT:
		conf.alertLog.Printf("[%s] - %s\n", node.name, *ev)
	case BOTH:
		conf.eventLog.Printf("[%s] - %s\n", node.name, *ev)
		conf.alertLog.Printf("[%s] - %s\n", node.name, *ev)
	}

	// recursively call all watching nodes
	for _, n := range node.nodes {
		RunNode(conf, n, ev, scripts)
	}
}

func SetupRuleTreeConfig(conf *ServerConfig) (*ini.File, error) {

	// get all files in the rules directory
	fileInfos, err := ioutil.ReadDir(conf.rulesDirPath)
	if err != nil {
		return nil, err
	}

	rulesConf := ini.Empty()
	// loop through and add them to the list
	for _, fileInfo := range fileInfos {

		// skip if not a .rules file
		if !strings.HasSuffix(fileInfo.Name(), ".rules") {
			continue
		}

		conf.serverLog.Println("Loading rules file " + conf.rulesDirPath + "/" + fileInfo.Name())

		// append to the array to be opened
		err := rulesConf.Append(conf.rulesDirPath + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}
	}

	return rulesConf, nil
}

func SetupRuleTree(conf *ServerConfig, rulesConf *ini.File) (*RuleNode, error) {

	// Create the base rule. A matching rule with no matches will match all
	baseNode := &RuleNode{
		action: MATCH,
	}

	// create a map for storing the rules
	ruleTree := make(map[string]*RuleNode)
	ruleTree["base"] = baseNode

	// for each rule section
	for _, rule := range rulesConf.Sections() {

		err := createRule(rule, conf, ruleTree)
		if err != nil {
			return nil, err
		}

	}

	return baseNode, nil
}

func createRule(rule *ini.Section, conf *ServerConfig, ruleTree map[string]*RuleNode) error {

	// skip the default section
	if rule.Name() == ini.DEFAULT_SECTION {
		return nil
	}

	// skip any rule that doesn't watch anything
	if !rule.HasKey("watch") {
		conf.serverLog.Printf("Rule [%s] has no watch, skipping", rule.Name())
		return nil
	}

	// check to ensure the watched rule exists
	if _, ok := ruleTree[rule.Key("watch").Value()]; !ok {
		conf.serverLog.Printf("Rule [%s] watching non-existant rule [%s], skipping", rule.Name(), rule.Key("watch").Value())
		return nil
	}

	// skip any rule with no action
	if !rule.HasKey("action") {
		conf.serverLog.Printf("Rule [%s] has no action, skipping", rule.Name())
		return nil
	}

	// create the rule
	var newRule *RuleNode
	var err error

	switch rule.Key("action").Value() {
	case "match":
		newRule, err = createMatchRule(rule)
	case "script":
		newRule, err = createScriptRule(rule)
	default:
		return fmt.Errorf("Rule [%s] has invalid action: %s", rule.Name(), rule.Key("action").Value())

	}

	// check if the rule created sucessfully
	if err != nil {
		conf.serverLog.Printf("Failed to create rule [%s]: %s", rule.Name(), err.Error())
		return err
	}

	// register it to watch the target rule
	watchedRule := ruleTree[newRule.watch]
	watchedRule.nodes = append(watchedRule.nodes, newRule)

	// add this rule to the tree
	ruleTree[newRule.name] = newRule

	return nil

}

func createMatchRule(conf *ini.Section) (*RuleNode, error) {
	rule := &RuleNode{
		name:    conf.Name(),
		action:  MATCH,
		watch:   conf.Key("watch").Value(),
		matches: make(map[string]*regexp.Regexp),
	}

	// setup rule triggers
	if err := addRuleTrigger(conf, rule); err != nil {
		return nil, err
	}

	// set the regex type
	regexType := conf.Key("regextype").Value()

	for _, key := range conf.Keys() {

		if !strings.HasPrefix(key.Name(), "match_") {
			continue
		}

		match := strings.TrimLeft(key.Name(), "match_")
		var regex *regexp.Regexp
		var err error

		switch regexType {
		case "posix":
			regex, err = regexp.CompilePOSIX(key.Value())
		case "default", "perl", "normal":
			regex, err = regexp.Compile(key.Value())
		default:
			return nil, fmt.Errorf("Unknown regex type %s", regexType)
		}

		if err != nil {
			return nil, err
		}

		rule.matches[match] = regex

	}

	return rule, nil
}

func createScriptRule(conf *ini.Section) (*RuleNode, error) {
	return nil, fmt.Errorf("Not implemented")
	//TODO
}

func addRuleTrigger(conf *ini.Section, rule *RuleNode) error {

	if conf.HasKey("trigger") {
		switch conf.Key("trigger").Value() {
		case "log":
			rule.trigger = LOG
		case "alert":
			rule.trigger = ALERT
		case "both":
			rule.trigger = BOTH
		default:
			return fmt.Errorf("Unknown trigger type %s", conf.Key("trigger").Value())
		}
	}
	return nil
}
