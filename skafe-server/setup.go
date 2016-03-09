package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
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
	tlsStrict   bool
	tlsConf     *tls.Config
	tlsKeyPath  string
	tlsCertPath string
	tlsCaPath   string
}

var (
	errTlsMissingCert error = fmt.Errorf("Missing tlscert entry in config")
	errTlsMissingKey  error = fmt.Errorf("Missing tlskey entry in config")
)

func setupConfig(cfgPath string) (*ServerConfig, error) {

	// create the default server config object
	cfg := &ServerConfig{
		serverLogPath: "/var/log/skafe/skafe-server.log",
		eventLogPath:  "/var/log/skafe/events.log",
		alertLogPath:  "/var/log/skafe/alerts.log",

		rulesDirPath: "/etc/skafe/rules",

		port: uint16(DEFAULT_PORT),

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

	if defSec.HasKey("port") {
		port, err := defSec.Key("port").Uint()
		if err != nil {
			return nil, err
		}

		if port > 65534 {
			return nil, fmt.Errorf("Port out of range")
		}

		cfg.port = uint16(port)
	}

	// if not specified, don't use TLS
	if !defSec.HasKey("tls") {
		cfg.tls = false
	} else {

		// get the bool value from the config
		useTls, err := defSec.Key("tls").Bool()
		if err != nil {
			return nil, err
		}

		// if false, no tls
		if !useTls {
			cfg.tls = false
		} else {

			cfg.tls = true

			// verify other params are set
			if !defSec.HasKey("tlscert") {
				return nil, errTlsMissingCert
			}
			if !defSec.HasKey("tlskey") {
				return nil, errTlsMissingKey
			}

			cfg.tlsCertPath = defSec.Key("tlscert").Value()
			cfg.tlsKeyPath = defSec.Key("tlskey").Value()

			// check for struct checking
			if key, err := defSec.GetKey("tlsca"); err == nil {
				cfg.tlsCaPath = key.Value()
				cfg.tlsStrict = true
			}

		}
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

func setupTLS(conf *ServerConfig) error {

	if !conf.tls {
		return nil
	}

	// attempt to load keypair
	cert, err := tls.LoadX509KeyPair(conf.tlsCertPath, conf.tlsKeyPath)
	if err != nil {
		return err
	}

	conf.tlsConf = &tls.Config{}

	conf.tlsConf.Certificates = append(conf.tlsConf.Certificates, cert)

	// check if we're strict checking
	if !conf.tlsStrict {
		return nil
	}

	// set the server to strict checking
	conf.tlsConf.ClientAuth = tls.RequireAndVerifyClientCert

	// create the empty CA pool
	conf.tlsConf.ClientCAs = x509.NewCertPool()

	// Load the CA's PEM encoded cert
	caPemData, err := ioutil.ReadFile(conf.tlsCaPath)
	if err != nil {
		return err
	}

	// Parse the PEM into readable data
	header, caData := pem.Decode(caPemData)
	if header == nil {
		return fmt.Errorf("Unable to parse CA pem data")
	}

	// Parse the certs
	caCerts, err := x509.ParseCertificates(caData)
	if err != nil {
		return err
	}

	// Load the certs into the trusted CA pool
	for _, caCert := range caCerts {
		conf.tlsConf.ClientCAs.AddCert(caCert)
	}

	return nil
}
