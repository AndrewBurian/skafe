package main

import (
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"log/syslog"
	"os"
)

type AuditEvent map[string]string

type ServerConfig struct {

	// paths
	serverLogPath string
	eventLogPath  string
	alertLogPath  string
	rulesDirPath  string

	// loggers
	serverLog *log.Logger
	eventLog  *log.Logger
	alertLog  *log.Logger

	// TCP params
	port uint16

	// TLS params
	tls         bool
	tlsCertPath string
	tlsKeyPath  string
}

func main() {

	// flag for specifying a non-default config directory
	flagConfFile := flag.String("conf", "/etc/skafe/skafe-server.conf", "Set the directory containing config and rules")

	flag.Parse()

	// Setup server configuration
	conf, err := setupConfig(*flagConfFile)
	if err != nil {
		log.Fatalf("Unable to setup server config: %s\n", err.Error())
	}

	// Setup loggers
	err = setupLoggers(conf)
	if err != nil {
		log.Fatalf("Unable to set up log files: %s\n", err.Error())
	}

	conf.serverLog.Println("SKAFE Server started!")
	defer conf.serverLog.Println("SKAFE Server terminated.")

	evChan := make(chan *AuditEvent, 32)
	ruleChan := make(chan *AuditEvent, 32)

	go ClientLink(conf, evChan)

	go RuleEngine(conf, ruleChan)

	QueueEvents(conf, evChan, ruleChan)

}

func setupConfig(cfgPath string) (*ServerConfig, error) {

	// create the default server config object
	cfg := &ServerConfig{
		serverLogPath: "/var/log/skafe/skafe-server.log",
		eventLogPath:  "/var/log/skafe/events.log",
		alertLogPath:  "/var/log/skafe/alerts.log",

		rulesDirPath: "/etc/skafe/rules",

		port: uint16(6969),

		tls: false,
	}

	// read the config file in
	cfgFile, err := ini.Load(cfgPath)
	if err != nil {
		return nil, err
	}

	// get the default section containing the setup
	defSec := cfgFile.Section(ini.DEFAULT_SECTION)

	if key, err := defSec.GetKey("serverlog"); err == nil {
		cfg.serverLogPath = key.Value()
	}

	if key, err := defSec.GetKey("eventlog"); err == nil {
		cfg.eventLogPath = key.Value()
	}

	if key, err := defSec.GetKey("alertlog"); err == nil {
		cfg.alertLogPath = key.Value()
	}

	if key, err := defSec.GetKey("rulesdir"); err == nil {
		cfg.rulesDirPath = key.Value()
	}

	return cfg, nil
}

func setupLoggers(conf *ServerConfig) error {

	var err error

	// System logger
	switch conf.serverLogPath {
	case "stdout", "stderr":
		conf.serverLog, err = setupStdioLogger(conf.serverLogPath)
	case "syslog":
		conf.serverLog, err = setupSysLogger()
	default:
		conf.serverLog, err = setupFileLogger(conf.serverLogPath)
	}

	if err != nil {
		return err
	}

	// Alert logger
	switch conf.alertLogPath {
	case "stdout", "stderr":
		conf.alertLog, err = setupStdioLogger(conf.alertLogPath)
	case "syslog":
		conf.alertLog, err = setupSysLogger()
	default:
		conf.alertLog, err = setupFileLogger(conf.alertLogPath)
	}

	if err != nil {
		return err
	}

	// Event logger
	switch conf.eventLogPath {
	case "stdout", "stderr":
		conf.eventLog, err = setupStdioLogger(conf.eventLogPath)
	case "syslog":
		conf.eventLog, err = setupSysLogger()
	default:
		conf.eventLog, err = setupFileLogger(conf.eventLogPath)
	}

	if err != nil {
		return err
	}

	return nil
}

func setupFileLogger(path string) (*log.Logger, error) {

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	return log.New(NewSyncWriter(logFile), "", log.LstdFlags), nil
}

func setupStdioLogger(path string) (*log.Logger, error) {
	switch path {
	case "stdout":
		return log.New(os.Stdout, "", log.LstdFlags), nil

	case "stderr":
		return log.New(os.Stderr, "", log.LstdFlags), nil
	default:
		return nil, fmt.Errorf("Why did you call setupStdioLogger with a non-stdio...")

	}
}

func setupSysLogger() (*log.Logger, error) {
	return syslog.NewLogger(syslog.LOG_ALERT|syslog.LOG_DAEMON, 0)
}
