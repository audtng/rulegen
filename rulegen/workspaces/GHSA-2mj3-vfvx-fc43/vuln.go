package main

}

func (s *snapshotter) EnsureLayer(ctx context.Context, key string) ([]layer.DiffID, error) {
	diffIDs, err := s.GetDiffIDs(ctx, key)
	if err != nil {
		return nil, err
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/snapshot"
	"github.com/moby/buildkit/util/leaseutil"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
type snapshotter struct {
	opt Opt

	refs map[string]layer.Layer
	db   *bolt.DB
	mu   sync.Mutex
	reg  graphIDRegistrar
}

// NewSnapshotter creates a new snapshotter
	}

	s := &snapshotter{
		opt:  opt,
		db:   db,
		refs: map[string]layer.Layer{},
		reg:  reg,
	}

	slm := newLeaseManager(s, prevLM)
