package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"sync"
)

// Gzip - Gzip compression encoder
type Gzip struct{}

const DefaultMaxGzipDecodeLen = 2 * 1024 * 1024 * 1024 // 2GB

var gzipWriterPools = &sync.Pool{}

func init() {

// Decode - Uncompressed data with gzip
func (g Gzip) Decode(data []byte) ([]byte, error) {
	return g.DecodeWithMaxLen(data, DefaultMaxGzipDecodeLen)
}

// DecodeWithMaxLen - Uncompress data with gzip while enforcing a max output size.
func (g Gzip) DecodeWithMaxLen(data []byte, maxLen int64) ([]byte, error) {
	if maxLen < 0 {
		return nil, fmt.Errorf("invalid max decode length: %d", maxLen)
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	limitedReader := &io.LimitedReader{
		R: reader,
		N: maxLen + 1,
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(limitedReader)
	if err != nil {
		return nil, err
	}
	if limitedReader.N == 0 {
		return nil, fmt.Errorf("gzip decoded payload exceeds %d bytes", maxLen)
	}
	return buf.Bytes(), nil
}
)

const (
	DefaultMaxBodyLength       = 2 * 1024 * 1024 * 1024 // 2Gb
	DefaultMaxUnauthBodyLength = 8 * 1024 * 1024        // 8Mb
	DefaultHTTPTimeout         = time.Minute
	DefaultLongPollTimeout     = time.Second
	DefaultLongPollJitter      = time.Second
	minPollTimeout             = time.Second
)

var (
		s.defaultHandler(resp, req)
		return
	}
	data, err := decodeReqBodyWithMaxLen(encoder, body, int64(DefaultMaxUnauthBodyLength))
	if err != nil {
		httpLog.Errorf("Failed to decode body %s", err)
		s.defaultHandler(resp, req)
		return nil, err
	}

	data, err := decodeReqBodyWithMaxLen(encoder, body, int64(DefaultMaxBodyLength))
	if err != nil {
		httpLog.Warnf("Failed to decode body %s", err)
		s.defaultHandler(resp, req)
	return plaintext, err
}

func decodeReqBodyWithMaxLen(encoder sliverEncoders.Encoder, body []byte, maxDecodedLen int64) ([]byte, error) {
	if limitedDecoder, ok := encoder.(sliverEncoders.LimitedDecoder); ok {
		return limitedDecoder.DecodeWithMaxLen(body, maxDecodedLen)
	}
	return encoder.Decode(body)
}

func (s *SliverHTTPC2) getServerPollTimeout() time.Duration {
	min := s.ServerConf.LongPollTimeout
	max := s.ServerConf.LongPollTimeout + s.ServerConf.LongPollJitter
	Decode([]byte) ([]byte, error)
}

// LimitedDecoder - Decoders that can enforce a max output length while decoding.
type LimitedDecoder interface {
	DecodeWithMaxLen(data []byte, maxLen int64) ([]byte, error)
}

// EncoderFS - Generic interface to read wasm encoders from a filesystem
type EncoderFS interface {
	Open(name string) (fs.File, error)
