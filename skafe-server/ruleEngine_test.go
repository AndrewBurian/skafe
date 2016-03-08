package main

import (
	"github.com/go-ini/ini"
	"log"
	"regexp"
	"testing"
)

type LogCounter struct {
	Count int
	t     *testing.T
}

func (l *LogCounter) Write(b []byte) (int, error) {
	l.Count += 1
	l.t.Log(string(b))
	return len(b), nil
}

func TestMatchNodeBlankEvent(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 1 {
		t.Errorf("Node did not trigger correct logs\nExpecting [1]\nGot [%d]", counter.Count)
	}
}

func TestMatchNodeBlankBoth(t *testing.T) {
	eventCounter := &LogCounter{
		t: t,
	}
	alertCounter := &LogCounter{t: t}

	eventLogger := log.New(eventCounter, "", 0)
	alertLogger := log.New(alertCounter, "", 0)

	node := &RuleNode{
		trigger: BOTH,
		action:  MATCH,
	}

	conf := &ServerConfig{
		eventLog: eventLogger,
		alertLog: alertLogger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if eventCounter.Count != 1 {
		t.Errorf("Node did not trigger correct event logs\nExpecting [1]\nGot [%d]", eventCounter.Count)
	}
	if alertCounter.Count != 1 {
		t.Errorf("Node did not trigger correct alert logs\nExpecting [1]\nGot [%d]", alertCounter.Count)
	}
}

func TestMatchNodeBlankAlert(t *testing.T) {
	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: ALERT,
		action:  MATCH,
	}

	conf := &ServerConfig{
		alertLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 1 {
		t.Errorf("Node did not trigger correct logs\nExpecting [1]\nGot [%d]", counter.Count)
	}
}

func TestMatchNodeNoMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key": regexp.MustCompile("other val"),
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 0 {
		t.Errorf("Node did not trigger correct logs\nExpecting [0]\nGot [%d]", counter.Count)
	}
}

func TestMatchNodeHalfMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key":  regexp.MustCompile("val"),
			"key2": regexp.MustCompile("other val2"),
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key":  "val",
		"key2": "val2",
	}

	RunNode(conf, node, event)

	if counter.Count != 0 {
		t.Errorf("Node did not trigger correct logs\nExpecting [0]\nGot [%d]", counter.Count)
	}
}

func TestRecursiveMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		name:    "first",
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key":  regexp.MustCompile("val"),
			"key2": regexp.MustCompile("val2"),
		},
		nodes: []*RuleNode{
			&RuleNode{
				name:    "second",
				trigger: LOG,
				action:  MATCH,
				matches: map[string]*regexp.Regexp{
					"key":  regexp.MustCompile("val"),
					"key2": regexp.MustCompile("val2"),
				},
			},
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key":  "val",
		"key2": "val2",
	}

	RunNode(conf, node, event)

	if counter.Count != 2 {
		t.Errorf("Node did not trigger correct logs\nExpecting [2]\nGot [%d]", counter.Count)
	}
}

func TestCreateMatchRuleGoodSingleEvent(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	rule, err := createMatchRule(file.Section("Rule"))
	checkErr(nil, err, t)

	if rule.name != "Rule" {
		t.Errorf("Rule name incorrect\nExpected [Rule]\nGot [%s]", rule.name)
	}
	if rule.trigger != LOG {
		t.Errorf("Rule received wrong trigger")
	}

	regex, ok := rule.matches["key"]

	if !ok {
		t.Errorf("Rule failed to get match key")
	}
	if regex == nil {
		t.Errorf("Failed to compile match regex")
	}
	if !regex.MatchString("val") {
		t.Errorf("Regex not working as intended. Missing match")
	}
	if regex.MatchString("other val") {
		t.Errorf("Regex not working as intended. Unexpected match")
	}
}

func TestCreateMatchRuleGoodSingleAlert(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		trigger = alert
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	rule, err := createMatchRule(file.Section("Rule"))
	checkErr(nil, err, t)

	if rule.name != "Rule" {
		t.Errorf("Rule name incorrect\nExpected [Rule]\nGot [%s]", rule.name)
	}
	if rule.trigger != ALERT {
		t.Errorf("Rule received wrong trigger")
	}

	regex, ok := rule.matches["key"]

	if !ok {
		t.Errorf("Rule failed to get match key")
	}
	if regex == nil {
		t.Errorf("Failed to compile match regex")
	}
	if !regex.MatchString("val") {
		t.Errorf("Regex not working as intended. Missing match")
	}
	if regex.MatchString("other val") {
		t.Errorf("Regex not working as intended. Unexpected match")
	}
}

func TestCreateMatchRuleGoodSingleBoth(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		trigger = both
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	rule, err := createMatchRule(file.Section("Rule"))
	checkErr(nil, err, t)

	if rule.name != "Rule" {
		t.Errorf("Rule name incorrect\nExpected [Rule]\nGot [%s]", rule.name)
	}
	if rule.trigger != BOTH {
		t.Errorf("Rule received wrong trigger")
	}

	regex, ok := rule.matches["key"]

	if !ok {
		t.Errorf("Rule failed to get match key")
	}
	if regex == nil {
		t.Errorf("Failed to compile match regex")
	}
	if !regex.MatchString("val") {
		t.Errorf("Regex not working as intended. Missing match")
	}
	if regex.MatchString("other val") {
		t.Errorf("Regex not working as intended. Unexpected match")
	}
}

func TestCreateMatchRuleBadTrigger(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		trigger = llama
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	_, err = createMatchRule(file.Section("Rule"))
	if err == nil {
		t.Errorf("No error was returned where one should have been")
	}
}

func TestCreateMatchRuleBadRegexType(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		trigger = log
		regextype = llama
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	_, err = createMatchRule(file.Section("Rule"))
	if err == nil {
		t.Errorf("No error was returned where one should have been")
	}
}

func TestCreateMatchRuleBadRegex(t *testing.T) {
	file, err := ini.Load([]byte(`
		[Rule]
		trigger = log
		regextype = posix
		match_key = [a-z+
	`))
	checkErr(nil, err, t)

	_, err = createMatchRule(file.Section("Rule"))
	if err == nil {
		t.Errorf("No error was returned where one should have been")
	}
}


func TestRulesConfInvalidDir(t *testing.T) {
	conf := &ServerConfig {
		rulesDirPath: "/dev/null/broken",
	}

	_, err := SetupRuleTreeConfig(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestRulesConfCountRules(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
		rulesDirPath: "tests/rulesets/rules_empty",
	}

	_, err := SetupRuleTreeConfig(conf)
	checkErr(nil, err, t)

	if counter.Count != 5 {
		t.Errorf("Wrong number of rule files loaded\nExpected [5]\nGot [%d]", counter.Count)
	}
}

func TestRulesConfIgnoreOther(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
		rulesDirPath: "tests/rulesets/rules_others",
	}

	_, err := SetupRuleTreeConfig(conf)
	checkErr(nil, err, t)

	if counter.Count != 3 {
		t.Errorf("Wrong number of rule files loaded\nExpected [3]\nGot [%d]", counter.Count)
	}
}

func TestRulesConfPermissionErr(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
		rulesDirPath: "tests/rulesets/rules_permission",
	}

	_, err := SetupRuleTreeConfig(conf)

	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}
