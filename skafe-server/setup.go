package main

import (
	"flag"
	"github.com/go-ini/ini"
	"log"
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

	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND

	// setup log files
	logFile, err := os.OpenFile(conf.serverLogPath, flags, 0640)
	if err != nil {
		return err
	}

	eventFile, err := os.OpenFile(conf.eventLogPath, flags, 0640)
	if err != nil {
		return err
	}

	alertFile, err := os.OpenFile(conf.alertLogPath, flags, 0640)
	if err != nil {
		return err
	}

	// create a SyncWriter so multiple goroutines can safely use a logger
	syncLog := NewSyncWriter(logFile)
	syncEvent := NewSyncWriter(eventFile)
	syncAlert := NewSyncWriter(alertFile)

	// target default logger to this file
	conf.serverLog = log.New(syncLog, "", log.LstdFlags)
	conf.eventLog = log.New(syncEvent, "", log.LstdFlags)
	conf.alertLog = log.New(syncAlert, "", log.LstdFlags)

	return nil
}
