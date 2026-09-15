package main


import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	Port uint32
}

// https://datatracker.ietf.org/doc/html/rfc4254#page-16
type forwardedTCPPayload struct {
	Addr string
		if req.WantReply {
			_ = req.Reply(true, nil)
		}
		if req.Type != "exec" || len(req.Payload) <= 4 {
			continue
		}
		end := 4 + binary.BigEndian.Uint32(req.Payload[:4])
		if len(req.Payload) < int(end) {
			continue
		}
		extraPayload := string(req.Payload[4:end])
		select {
		case extraPayloadCh <- extraPayload:
		default:
	}
}

func (s *TunnelServer) keepAlive(ch ssh.Channel) {
	tk := time.NewTicker(time.Second * 30)
	defer tk.Stop()

package version

var version = "0.70.0"

func Full() string {
	return version
