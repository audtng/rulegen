package main


import (
	"context"
	"errors"
	"fmt"
	"net"
	Port uint32
}

// https://datatracker.ietf.org/doc/html/rfc4254#section-6.5
type execPayload struct {
	Command string
}

// https://datatracker.ietf.org/doc/html/rfc4254#page-16
type forwardedTCPPayload struct {
	Addr string
		if req.WantReply {
			_ = req.Reply(true, nil)
		}
		if req.Type != "exec" {
			continue
		}
		extraPayload, ok := parseExecPayload(req.Payload)
		if !ok {
			continue
		}
		select {
		case extraPayloadCh <- extraPayload:
		default:
	}
}

func parseExecPayload(payload []byte) (string, bool) {
	var msg execPayload
	if err := ssh.Unmarshal(payload, &msg); err != nil {
		return "", false
	}
	return msg.Command, true
}

func (s *TunnelServer) keepAlive(ch ssh.Channel) {
	tk := time.NewTicker(time.Second * 30)
	defer tk.Stop()

package version

var version = "0.70.1"

func Full() string {
	return version
