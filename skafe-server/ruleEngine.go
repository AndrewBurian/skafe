package main

import (
	"github.com/robertkrimen/otto"
	"regexp"
)

const (
	MATCH  int = 1
	SCRIPT int = 2
)

type RuleNode struct {
	name    string
	action  int
	matches map[string]*regexp.Regexp
	script  *otto.Script
	nodes   []*RuleNode
}

func RuleEngine(conf *ServerConfig, events <-chan *AuditEvent) {

	// the channel to dispatch events to workers over
	//workers := make(chan *AuditEvent)

}

func RuleEngineWorker(base *RuleNode, events <-chan *AuditEvent) {

	// so long as there are events to process
	for ev := range events {

		// recurively follow the descision tree
		RunNode(base, ev)
	}
}

func RunNode(node *RuleNode, ev *AuditEvent) {

	if node.action == SCRIPT {
		// TODO Scripts
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

	// recursively call all watching nodes
	for _, n := range node.nodes {
		RunNode(n, ev)
	}
}
