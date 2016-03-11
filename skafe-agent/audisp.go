package main

import (
	"bufio"
	audit "github.com/andrewburian/libaudit-go"
	"io"
)

func Audisp(receivedEvents chan<- AuditEvent, in io.Reader) error {

	bufReader := bufio.NewReader(in)

	for {
		eventString, readErr := bufReader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return readErr
		}

		event, parseErr := ParseEvent(eventString)
		if parseErr != nil {
			return parseErr
		}

		receivedEvents <- event

		if readErr == io.EOF {
			break
		}
	}

	close(receivedEvents)

	return nil

}

func ParseEvent(data string) (AuditEvent, error) {
	_, _, m, err := audit.ParseAuditEvent(data)
	if err != nil {
		return nil, err
	}

	return AuditEvent(m), nil
}
