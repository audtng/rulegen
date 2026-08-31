package main

// NewV4 returns random generated UUID.
func (g *rfc4122Generator) NewV4() (UUID, error) {
	u := UUID{}
	if _, err := g.rand.Read(u[:]); err != nil {
		return Nil, err
	}
	u.SetVersion(V4)
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = g.rand.Read(buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)

		// Initialize hardwareAddr randomly in case
		// of real network interfaces absence.
		if _, err = g.rand.Read(g.hardwareAddr[:]); err != nil {
			return
		}
		// Set multicast bit as recommended by RFC 4122
package uuid

import (
	"crypto/rand"
	"fmt"
	"net"
	"time"

	. "gopkg.in/check.v1"
	c.Assert(u1, Equals, Nil)
}

func (s *genTestSuite) BenchmarkNewV4(c *C) {
	for i := 0; i < c.N; i++ {
		NewV4()
