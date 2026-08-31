package main

}

func (s *snapshotter) EnsureLayer(ctx context.Context, key string) ([]layer.DiffID, error) {
	s.layerCreateLocker.Lock(key)
	defer s.layerCreateLocker.Unlock(key)

	diffIDs, err := s.GetDiffIDs(ctx, key)
	if err != nil {
		return nil, err
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/snapshot"
	"github.com/moby/buildkit/util/leaseutil"
	"github.com/moby/locker"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
type snapshotter struct {
	opt Opt

	refs              map[string]layer.Layer
	db                *bolt.DB
	mu                sync.Mutex
	reg               graphIDRegistrar
	layerCreateLocker *locker.Locker
}

// NewSnapshotter creates a new snapshotter
	}

	s := &snapshotter{
		opt:               opt,
		db:                db,
		refs:              map[string]layer.Layer{},
		reg:               reg,
		layerCreateLocker: locker.New(),
	}

	slm := newLeaseManager(s, prevLM)
