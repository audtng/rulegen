package main

	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/envconfig"
	imem "google.golang.org/grpc/internal/mem"
	"google.golang.org/grpc/internal/transport/internal"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/tap"
)

const (
	logLevel = 2
	// recvMsgSize estimates the memory overhead of a recvMsg in the backlog.
	// It accounts for the recvMsg struct itself and the slice header of the
	// underlying buffer's data.
	recvMsgSize = int(unsafe.Sizeof(recvMsg{}) + unsafe.Sizeof([]byte{}))

	// utilizationFactor controls when we consider memory utilization acceptable.
	// When backlogHeapSize / payloadSize <= utilizationFactor (meaning at least
	// 50% of the heap memory is actual payload data), compaction is skipped.
	utilizationFactor = 2
)

var (
	// compactionThreshold is approx 57KB (on 64-bit systems). It allows
	// accumulating up to 1024 1-byte payloads before triggering compaction.
	//
	// Because individual payloads <= 1024 bytes are allocated on the heap
	// outside mem.BufferPool, waiting for at least 1024 bytes to accumulate
	// ensures that compaction coalesces those small heap allocations into a
	// single large buffer from mem.BufferPool, enabling buffer reuse while
	// avoiding frequent copying for small bursts of frames.
	compactionThreshold = imem.BufferPoolingThreshold * (recvMsgSize + 1)
)

func init() {
	internal.TimeNowFunc = func() int64 { return time.Now().UnixNano() }
	c       chan recvMsg
	mu      sync.Mutex
	backlog []recvMsg
	// uncompactedSuffixLen tracks the number of consecutive data messages at
	// the tail of backlog that have not been compacted.
	uncompactedSuffixLen int
	// uncompactedBytes tracks the total payload bytes across the trailing
	// uncompactedSuffixLen messages.
	uncompactedBytes int
	err              error
	bufPool          mem.BufferPool
}

// init allows a recvBuffer to be initialized in-place, which is useful
// for resetting a buffer or for avoiding a heap allocation when the buffer
// is embedded in another struct.
func (b *recvBuffer) init(pool mem.BufferPool) {
	b.c = make(chan recvMsg, 1)
	b.bufPool = pool
}

func (b *recvBuffer) put(r recvMsg) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.err != nil {
		// drop the buffer on the floor. Since b.err is not nil, any subsequent reads
		// will always return an error, making this buffer inaccessible.
		r.buffer.Free()
		// An error had occurred earlier, don't accept more
		// data or errors.
		return
	if len(b.backlog) == 0 {
		select {
		case b.c <- r:
			return
		default:
		}
	}
	b.backlog = append(b.backlog, r)
	b.compactBacklogLocked(r)
}

func (b *recvBuffer) compactBacklogLocked(r recvMsg) {
	if !envconfig.EnableReceiveBufferCompaction {
		return
	}
	if r.buffer == nil {
		b.uncompactedBytes = 0
		b.uncompactedSuffixLen = 0
		return
	}

	b.uncompactedSuffixLen++
	b.uncompactedBytes += r.buffer.Len()
	backlogHeapSize := b.uncompactedSuffixLen*recvMsgSize + b.uncompactedBytes

	// If the memory overhead is less than 50% of the heap usage (e.g., because
	// a large DATA frame arrived), the average message size in the suffix is
	// large enough that memory bloat is not a concern. Reset suffix tracking.
	if backlogHeapSize <= utilizationFactor*b.uncompactedBytes {
		b.uncompactedBytes = 0
		b.uncompactedSuffixLen = 0
		return
	}
	// Avoid compacting too frequently for short bursts of small frames.
	// Wait until we have accumulated at least ~1024 small messages (~57 KB).
	if backlogHeapSize <= compactionThreshold {
		// Still can accumulate more payloads.
		return
	}

	// Since the memory utilization is less than 50%, the average payload size
	// of each recvMsg must be less than recvMsgSize (approx 56 bytes).
	// In the worst case for bytes copied (where the average payload is just
	// below recvMsgSize), compaction will occur once every:
	//   compactionThreshold / (recvMsgSize + avg_payload) = ~520 messages,
	// copying ~29KB of data.

	start := 0
	newBuf := b.bufPool.Get(b.uncompactedBytes)
	startIdx := len(b.backlog) - b.uncompactedSuffixLen

	for i := startIdx; i < len(b.backlog); i++ {
		m := b.backlog[i]
		b.backlog[i] = recvMsg{}
		start += copy((*newBuf)[start:], m.buffer.ReadOnlyData())
		m.buffer.Free()
	}
	b.backlog[startIdx] = recvMsg{
		buffer: mem.NewBuffer(newBuf, b.bufPool),
	}
	b.backlog = b.backlog[:startIdx+1]
	// After compaction, the suffix is replaced with a single message containing
	// the combined payload. The new utilization is close to 1.0 (overhead of
	// one recvMsg relative to the large compacted payload), which is well
	// below the utilization factor of 2.
	b.uncompactedBytes = 0
	b.uncompactedSuffixLen = 0
}

func (b *recvBuffer) load() {
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			// backlog[0] is only part of the tracked uncompacted suffix if the
			// entire backlog currently consists of the suffix. If an earlier
			// compaction or reset occurred, backlog[0] is already compacted.
			if envconfig.EnableReceiveBufferCompaction && b.uncompactedSuffixLen == len(b.backlog) {
				b.uncompactedSuffixLen--
				b.uncompactedBytes -= b.backlog[0].buffer.Len()
			}
			b.backlog[0] = recvMsg{}
			b.backlog = b.backlog[1:]
		default:
	"slices"
	"sort"
	"sync"

	"google.golang.org/grpc/internal"
)

const (
	goPageSize = 4 * 1024 // 4KiB. N.B. this must be a power of 2.
)

var (
	// BufferPoolingThreshold is the minimum size of a buffer that can be pooled.
	// This is used to determine whether to pool buffers or allocate them directly.
	BufferPoolingThreshold = 1 << 10
)

func init() {
	internal.SetBufferPoolingThresholdForTesting = func(threshold int) {
		BufferPoolingThreshold = threshold
	}
}

var uintSize = bits.UintSize // use a variable for mocking during tests.

// bufferPool is a copy of the public bufferPool interface used to avoid
	// throttling limit if unforeseen issues arise, and it will be removed in a
	// future release.
	//
	// TODO: Remove this env var once v1.83.0 is released.
	ControlBufferThrottleLimit = uint64FromEnv("GRPC_GO_EXPERIMENTAL_CONTROL_BUFFER_THROTTLE_LIMIT", 100, 1, 10000)

	// EnableReceiveBufferCompaction enables the compaction of data buffers
	// to reduce the number of buffers in the receive buffer.
	//
	// This environment variable serves as an escape hatch to disable the
	// feature if unforeseen issues arise, and it will be removed in a future
	// release.
	//
	// TODO: Remove this env var once v1.85.0 is released.
	EnableReceiveBufferCompaction = boolFromEnv("GRPC_GO_EXPERIMENTAL_ENABLE_RECEIVE_BUFFER_COMPACTION", true)
)

func boolFromEnv(envVar string, def bool) bool {
