package main

import (
	"crypto/tls"
	"testing"
)

func checkErr(expected, actual error, t *testing.T) {
	if expected == nil && actual != nil {
		t.Fatalf("Error returned where none expected: \n[%s]", actual)
		return
	}

	if expected != nil && actual == nil {
		t.Fatalf("Error not returned where one should have been.")
		return
	}

	if expected != actual {
		t.Fatalf("Wrong error returned.\nExpected: [%s]\nGot: [%s]", expected, actual)
		return
	}
}

// Test that all appropriate tls* keys in the conf
// parse correctly
func TestTLSConfig(t *testing.T) {

	_, err := setupConfig("tests/tls_good.conf")
	checkErr(nil, err, t)
}

func TestTLSConfigMissingCert(t *testing.T) {
	_, err := setupConfig("tests/tls_missing_cert.conf")
	checkErr(errTlsMissingCert, err, t)
}

func TestTLSConfigMissingKey(t *testing.T) {
	_, err := setupConfig("tests/tls_missing_key.conf")
	checkErr(errTlsMissingKey, err, t)
}

// Ensure a valid cert and key are accepted
func TestTLSSetupGood(t *testing.T) {

	conf, err := setupConfig("tests/tls_good.conf")
	checkErr(nil, err, t)
	err = setupTLS(conf)
	checkErr(nil, err, t)

	if conf.tls != true {
		t.Errorf("TLS didn't enable when it should have")
	}
	if conf.tlsConf == nil {
		t.Errorf("TLS Conf not created")
	}
}

func TestTLSSetupGoodWithCA(t *testing.T) {
	conf, err := setupConfig("tests/tls_goodca.conf")
	checkErr(nil, err, t)
	err = setupTLS(conf)
	checkErr(nil, err, t)

	if conf.tlsConf.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Errorf("TLS didn't enable strict client checking")
	}

	if conf.tlsConf.ClientCAs == nil {
		t.Errorf("TLS didn't setup client CA pool")
	}
}

func TestTLSSetupBadCA(t *testing.T) {
	conf, err := setupConfig("tests/tls_badca.conf")
	checkErr(nil, err, t)
	err = setupTLS(conf)

	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

// Encure a bad cert crashes
func TestTLSSetupBadCert(t *testing.T) {

	conf := &ServerConfig{
		tls:         true,
		tlsCertPath: "tls/nonexistant.pem",
		tlsKeyPath:  "tls/demoserver.key",
	}

	err := setupTLS(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestTLSSetupBadKey(t *testing.T) {

	conf := &ServerConfig{
		tls:         true,
		tlsCertPath: "tls/demoserver.pem",
		tlsKeyPath:  "tls/nonexistant.key",
	}

	err := setupTLS(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestFileLoggerSetupGood(t *testing.T) {
	_, err := setupFileLogger("tests/filelogger.log")
	checkErr(nil, err, t)
}

func TestFileLoggerSetupFailed(t *testing.T) {
	_, err := setupFileLogger("/dev/null/impossible/why/ohgod/log.log")
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestFileLoggerSetup(t *testing.T) {
	// this test cannot fail
	t.Skip("Syslog can't fail, will panic instead")
}

func TestStdioLoggerNonsense(t *testing.T) {
	_, err := setupStdioLogger("llama")
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestSetupLoggersFailServer(t *testing.T) {
	conf, err := setupConfig("tests/loggers_fail_server.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestSetupLoggersFailEvent(t *testing.T) {
	conf, err := setupConfig("tests/loggers_fail_event.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestSetupLoggersFailAlert(t *testing.T) {
	conf, err := setupConfig("tests/loggers_fail_alert.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}
func TestSetupConfigFail(t *testing.T) {
	_, err := setupConfig("/dev/null/impossible/gooby/y/u/do.conf")
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestTlsConfigNonbool(t *testing.T) {
	_, err := setupConfig("tests/tls_nonbool.conf")
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
}

func TestTlsDisable(t *testing.T) {
	conf, err := setupConfig("tests/tls_disable.conf")
	checkErr(nil, err, t)
	if conf.tls {
		t.Errorf("TLS failed to be disabled")
	}
}

func TestTlsSetupWithDisabled(t *testing.T) {
	conf := &ServerConfig{
		tls: false,
	}

	err := setupTLS(conf)
	checkErr(nil, err, t)
}

func TestStdoutLoggers(t *testing.T) {
	conf, err := setupConfig("tests/loggers_stdout.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	checkErr(nil, err, t)
}

func TestStderrLoggers(t *testing.T) {
	conf, err := setupConfig("tests/loggers_stderr.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	checkErr(nil, err, t)
}

func TestSysLoggers(t *testing.T) {
	conf, err := setupConfig("tests/loggers_syslog.conf")
	checkErr(nil, err, t)
	err = setupLoggers(conf)
	checkErr(nil, err, t)
}

func TestFileLoggers(t *testing.T) {
	conf := &ServerConfig{
		serverLogPath: "tests/server.log",
		alertLogPath:  "tests/alert.log",
		eventLogPath:  "tests/event.log",
	}

	err := setupLoggers(conf)
	checkErr(nil, err, t)
}

func TestListenPortDefault(t *testing.T) {
	conf, err := setupConfig("tests/blank.conf")
	checkErr(nil, err, t)
	if conf.port != DEFAULT_PORT {
		t.Errorf("Default port expected\nGot port [%d]", conf.port)
	}
}

func TestListenPortSet(t *testing.T) {
	conf, err := setupConfig("tests/port.conf")
	checkErr(nil, err, t)
	if conf.port != 6666 {
		t.Errorf("Got incorrect listen port.\nExpected [%d]\nGot [%d]", 6666, conf.port)
	}
}

func TestListenPortTooHigh(t *testing.T) {
	_, err := setupConfig("tests/port_high.conf")
	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}

func TestListenPortNan(t *testing.T) {
	_, err := setupConfig("tests/port_nan.conf")
	if err == nil {
		t.Errorf("No error returned where one was expected")
	}
}
