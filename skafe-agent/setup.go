package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"log/syslog"
	"os"
)

type AgentConfig struct {
	addr string
	port uint16

	log *log.Logger

	tlsConf *tls.Config
}

func setupConfig(confPath string) (*AgentConfig, error) {

	// create default config
	conf := &AgentConfig{

		port: DEFAULT_PORT,
	}

	// read config file
	confFile, err := ini.Load(confPath)
	if err != nil {
		return nil, err
	}

	// all entries in the default section
	sect := confFile.Section(ini.DEFAULT_SECTION)

	// get addr
	if addr, err := sect.GetKey("server"); err == nil {
		conf.addr = addr.Value()
	} else {
		return nil, fmt.Errorf("No server address specified")
	}

	// get port
	if sect.HasKey("port") {
		port, err := sect.Key("port").Uint()
		if err != nil {
			return nil, err
		}
		if port > 65534 {
			return nil, fmt.Errorf("Port out of range")
		}

		conf.port = uint16(port)
	}

	// get logger

	// default log path
	var logPath string = "/var/log/skafe/skafe-agent.log"

	// check for override
	if logKey, err := sect.GetKey("log"); err == nil {
		logPath = logKey.Value()
	}

	// create the logger
	switch logPath {

	// std outputs
	case "stdout":
		conf.log = log.New(os.Stdout, "", log.LstdFlags)

	case "stderr":
		conf.log = log.New(os.Stderr, "", log.LstdFlags)

	// syslog
	case "syslog":
		conf.log, err = syslog.NewLogger(syslog.LOG_ALERT|syslog.LOG_DAEMON, 0)
		if err != nil {
			return nil, err
		}

	// log file
	default:
		logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			return nil, err
		}
		conf.log = log.New(logFile, "", log.LstdFlags)
	}

	return conf, nil

}
