package main

import (
	audit "github.com/andrewburian/libaudit-go"
	"io"
)

func Audisp(receivedEvents chan<- AuditEvent, in io.Reader) {

}

func ParseEvent(data string) (AuditEvent, error) {
	_, _, m, err := audit.ParseAuditEvent(data)
	if err != nil {
		return nil, err
	}

	return AuditEvent(m), nil
}
