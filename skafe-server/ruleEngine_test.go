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

	RunNode(conf, node, event, nil)

	if counter.Count != 1 {
		t.Errorf("Node did not trigger correct logs\nExpecting [1]\nGot [%d]", counter.Count)
	}
}

func TestNoMatchMissingKey(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key": nil,
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key2": "val2",
	}

	RunNode(conf, node, event, nil)

	if counter.Count != 0 {
		t.Errorf("Node triggered log event where it shouldn't have")
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

	RunNode(conf, node, event, nil)

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

	RunNode(conf, node, event, nil)

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

	RunNode(conf, node, event, nil)

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

	RunNode(conf, node, event, nil)

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

	RunNode(conf, node, event, nil)

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
	conf := &ServerConfig{
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
		serverLog:    logger,
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
		serverLog:    logger,
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
		serverLog:    logger,
		rulesDirPath: "tests/rulesets/rules_permission",
	}

	_, err := SetupRuleTreeConfig(conf)

	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}

func TestCreateRuleSkipDefault(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		watch = base
		action = match
		trigger = llama
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section(""), conf, ruleTree, nil)
	checkErr(nil, err, t)

	if counter.Count != 0 {
		t.Errorf("Something else happened")
	}
}

func TestCreateRuleSkipNoWatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		action = match
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), conf, ruleTree, nil)
	checkErr(nil, err, t)

	if counter.Count != 1 {
		t.Errorf("Rule not skipped")
	}
}

func TestCreateRuleSkipNoParent(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		action = match
		watch = llama
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), conf, ruleTree, nil)
	checkErr(nil, err, t)

	if counter.Count != 1 {
		t.Errorf("Rule not skipped")
	}
}

func TestCreateRuleSkipNoAction(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		watch = base
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), conf, ruleTree, nil)
	checkErr(nil, err, t)

	if counter.Count != 1 {
		t.Errorf("Rule not skipped")
	}
}

func TestCreateRuleBadAction(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		action = llama
		watch = base
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), conf, ruleTree, nil)
	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}

func TestCreateRuleBadTrigger(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	conf := &ServerConfig{
		serverLog: logger,
	}

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		action = match
		watch = base
		trigger = llama
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), conf, ruleTree, nil)
	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}

func TestCreateRuleMatch(t *testing.T) {

	ruleTree := map[string]*RuleNode{
		"base": &RuleNode{},
	}

	file, err := ini.Load([]byte(`
		[Rule]
		action = match
		watch = base
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	err = createRule(file.Section("Rule"), nil, ruleTree, nil)
	checkErr(nil, err, t)

	rule, ok := ruleTree["Rule"]
	if !ok {
		t.Errorf("Rule failed to get added to tree")
	}

	if ruleTree["base"].nodes[0] != rule {
		t.Errorf("Rule not registerd with base")
	}

	if rule.action != MATCH {
		t.Errorf("Rule got assigned wrong action")
	}

	if rule.trigger != LOG {
		t.Errorf("Rule got assigned wrong trigger")
	}

}

func TestSetupScriptRule(t *testing.T) {

	// script rule entry
	file, err := ini.Load([]byte(`
		[Rule]
		action = script
		watch = base
		trigger = log
		script = testscript
		lang = ruby
	`))
	checkErr(nil, err, t)

	// Server Config
	conf := &ServerConfig{
		execRuby: "ruby",
	}

	// script pool
	pool, err := SetupScriptPool(conf)
	checkErr(nil, err, t)

	// Script rule
	rule, err := createScriptRule(file.Section("Rule"), pool)
	checkErr(nil, err, t)

	// verify script parts
	if rule.name != "Rule" {
		t.Errorf("Rule name mismatch\nExpected [%s]\nGot [%s]", "Rule", rule.name)
	}

	if rule.action != SCRIPT {
		t.Errorf("Rule got wrong action")
	}

	if rule.watch != "base" {
		t.Errorf("Rule got wrong watch\nExpected [%s]\nGot [%s]", "base", rule.watch)
	}

	if rule.trigger != LOG {
		t.Errorf("Rule got wrong trigger")
	}

	if rule.lang != SCRIPT_LANG_RUBY {
		t.Errorf("Rule got wrong lang\nExpected [%s]\nGot [%s]", SCRIPT_LANG_RUBY, rule.lang)
	}

	if rule.script != "testscript" {
		t.Errorf("Rule got wrong script\nExpected [%s]\nGot [%s]", "testscript", rule.script)
	}
}

func TestSetupRuleTreePass(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		action = match
		watch = base
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	base, err := SetupRuleTree(nil, file, nil)
	checkErr(nil, err, t)

	if base == nil {
		t.Errorf("Base node retruned as nil")
	}

	if len(base.nodes) != 1 {
		t.Errorf("Child nodes not created")
	}

}

func TestSetupRuleTreeFail(t *testing.T) {

	file, err := ini.Load([]byte(`
		[Rule]
		action = llama
		watch = base
		trigger = log
		regextype = perl
		match_key = ^val$
	`))
	checkErr(nil, err, t)

	base, err := SetupRuleTree(nil, file, nil)

	if err == nil || base != nil {
		t.Errorf("No error returned or base node initialized when it shouldn't have been")
	}
}
