package main

	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/transport/internal"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/tap"
)

const logLevel = 2

func init() {
	internal.TimeNowFunc = func() int64 { return time.Now().UnixNano() }
	c       chan recvMsg
	mu      sync.Mutex
	backlog []recvMsg
	err     error
}

// init allows a recvBuffer to be initialized in-place, which is useful
// for resetting a buffer or for avoiding a heap allocation when the buffer
// is embedded in another struct.
func (b *recvBuffer) init() {
	b.c = make(chan recvMsg, 1)
}

func (b *recvBuffer) put(r recvMsg) {
	b.mu.Lock()
	if b.err != nil {
		// drop the buffer on the floor. Since b.err is not nil, any subsequent reads
		// will always return an error, making this buffer inaccessible.
		r.buffer.Free()
		b.mu.Unlock()
		// An error had occurred earlier, don't accept more
		// data or errors.
		return
	if len(b.backlog) == 0 {
		select {
		case b.c <- r:
			b.mu.Unlock()
			return
		default:
		}
	}
	b.backlog = append(b.backlog, r)
	b.mu.Unlock()
}

func (b *recvBuffer) load() {
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = recvMsg{}
			b.backlog = b.backlog[1:]
		default:
	"slices"
	"sort"
	"sync"
)

const (
	goPageSize = 4 * 1024 // 4KiB. N.B. this must be a power of 2.
)

var uintSize = bits.UintSize // use a variable for mocking during tests.

// bufferPool is a copy of the public bufferPool interface used to avoid
	// throttling limit if unforeseen issues arise, and it will be removed in a
	// future release.
	//
	// TODO: Remove this env var once v1.83.0 is release.
	ControlBufferThrottleLimit = uint64FromEnv("GRPC_GO_EXPERIMENTAL_CONTROL_BUFFER_THROTTLE_LIMIT", 100, 1, 10000)
)

func boolFromEnv(envVar string, def bool) bool {
