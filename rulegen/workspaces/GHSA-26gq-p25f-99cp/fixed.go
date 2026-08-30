package main

    - legacy
    - std-error-handling
    rules:
    # Keep h2c and HTTP/1.1 Upgrade:h2c compatibility until an explicit migration.
    - linters:
      - staticcheck
      path: ^pkg/util/vhost/http\.go$
      text: "^SA1019:"
    - linters:
      - errcheck
      - maligned

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
// Copyright 2026 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	cryptossh "golang.org/x/crypto/ssh"
)

func TestParseExecPayload(t *testing.T) {
	payload := cryptossh.Marshal(&execPayload{Command: "tcp --remote_port 6000"})

	got, ok := parseExecPayload(payload)

	require.True(t, ok)
	require.Equal(t, "tcp --remote_port 6000", got)
}

func TestParseExecPayloadRejectsMalformedPayloads(t *testing.T) {
	overflowLength := make([]byte, 5)
	binary.BigEndian.PutUint32(overflowLength[:4], ^uint32(0))

	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{
			name:    "empty",
			payload: nil,
		},
		{
			name:    "short length prefix",
			payload: []byte{0, 0, 0},
		},
		{
			name:    "declared length exceeds remaining payload",
			payload: []byte{0, 0, 0, 2, 'x'},
		},
		{
			name:    "overflow length",
			payload: overflowLength,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				got string
				ok  bool
			)
			require.NotPanics(t, func() {
				got, ok = parseExecPayload(tc.payload)
			})
			require.False(t, ok)
			require.Empty(t, got)
		})
	}
}

package version

var version = "0.70.1"

func Full() string {
	return version
