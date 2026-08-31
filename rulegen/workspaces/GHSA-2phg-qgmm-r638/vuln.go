package main

)

const (
	DefaultMaxBodyLength   = 2 * 1024 * 1024 * 1024 // 2Gb
	DefaultMaxUnauthBodyLength = 8 * 1024 * 1024   // 8Mb
	DefaultHTTPTimeout     = time.Minute
	DefaultLongPollTimeout = time.Second
	DefaultLongPollJitter  = time.Second
	minPollTimeout         = time.Second
)

var (
		s.defaultHandler(resp, req)
		return
	}
	data, err := encoder.Decode(body)
	if err != nil {
		httpLog.Errorf("Failed to decode body %s", err)
		s.defaultHandler(resp, req)
		return nil, err
	}

	data, err := encoder.Decode(body)
	if err != nil {
		httpLog.Warnf("Failed to decode body %s", err)
		s.defaultHandler(resp, req)
	return plaintext, err
}

func (s *SliverHTTPC2) getServerPollTimeout() time.Duration {
	min := s.ServerConf.LongPollTimeout
	max := s.ServerConf.LongPollTimeout + s.ServerConf.LongPollJitter
	Decode([]byte) ([]byte, error)
}

// EncoderFS - Generic interface to read wasm encoders from a filesystem
type EncoderFS interface {
	Open(name string) (fs.File, error)
import (
	"bytes"
	"compress/gzip"
	"sync"
)

// Gzip - Gzip compression encoder
type Gzip struct{}

var gzipWriterPools = &sync.Pool{}

func init() {

// Decode - Uncompressed data with gzip
func (g Gzip) Decode(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
		}
	}
}
