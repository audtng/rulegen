package main

package etcdserver

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	pb "go.etcd.io/etcd/etcdserver/etcdserverpb"
	"go.etcd.io/etcd/mvcc"
	"go.etcd.io/etcd/pkg/types"

	"go.uber.org/zap"
		mismatch(uint64(s.ID()))
	}

	for _, p := range peers {
		if p.resp == nil {
			continue
		}
		id := p.resp.Header.MemberId

		// leader expects follower's latest revision less than or equal to leader's
			mismatch(id)
		}
	}
	return nil
}

type peerHashKVResp struct {
	id  types.ID
	eps []string

	resp *clientv3.HashKVResponse
	err  error
}

func (s *EtcdServer) getPeerHashKVs(rev int64) (resps []*peerHashKVResp) {
	// TODO: handle the case when "s.cluster.Members" have not
	// been populated (e.g. no snapshot to load from disk)
	mbs := s.cluster.Members()
	pss := make([]peerHashKVResp, len(mbs))
	for _, m := range mbs {
		if m.ID == s.ID() {
			continue
		}
		pss = append(pss, peerHashKVResp{id: m.ID, eps: m.PeerURLs})
	}

	lg := s.getLogger()

	for _, p := range pss {
		if len(p.eps) == 0 {
			continue
		}
		cli, cerr := clientv3.New(clientv3.Config{
			DialTimeout: s.Cfg.ReqTimeout(),
			Endpoints:   p.eps,
		})
		if cerr != nil {
			if lg != nil {
				lg.Warn(
					"failed to create client to peer URL",
					zap.String("local-member-id", s.ID().String()),
					zap.String("remote-peer-id", p.id.String()),
					zap.Strings("remote-peer-endpoints", p.eps),
					zap.Error(cerr),
				)
			} else {
				plog.Warningf("%s failed to create client to peer %q for hash checking (%q)", s.ID(), p.eps, cerr.Error())
			}
			continue
		}

		respsLen := len(resps)
		for _, c := range cli.Endpoints() {
			ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.ReqTimeout())
			var resp *clientv3.HashKVResponse
			resp, cerr = cli.HashKV(ctx, c, rev)
			cancel()
			if cerr == nil {
				resps = append(resps, &peerHashKVResp{id: p.id, eps: p.eps, resp: resp, err: nil})
				break
			}
			if lg != nil {
					"failed hash kv request",
					zap.String("local-member-id", s.ID().String()),
					zap.Int64("requested-revision", rev),
					zap.String("remote-peer-endpoint", c),
					zap.Error(cerr),
				)
			} else {
				plog.Warningf("%s hash-kv error %q on peer %q with revision %d", s.ID(), cerr.Error(), c, rev)
			}
		}
		cli.Close()

		if respsLen == len(resps) {
			resps = append(resps, &peerHashKVResp{id: p.id, eps: p.eps, resp: nil, err: cerr})
		}
	}
	return resps

func newApplierV3Corrupt(a applierV3) *applierV3Corrupt { return &applierV3Corrupt{a} }

func (a *applierV3Corrupt) Put(txn mvcc.TxnWrite, p *pb.PutRequest) (*pb.PutResponse, error) {
	return nil, ErrCorrupt
}

func (a *applierV3Corrupt) Range(txn mvcc.TxnRead, p *pb.RangeRequest) (*pb.RangeResponse, error) {
	return nil, ErrCorrupt
}

	return nil, ErrCorrupt
}

func (a *applierV3Corrupt) Compaction(compaction *pb.CompactionRequest) (*pb.CompactionResponse, <-chan struct{}, error) {
	return nil, nil, ErrCorrupt
}

func (a *applierV3Corrupt) LeaseGrant(lc *pb.LeaseGrantRequest) (*pb.LeaseGrantResponse, error) {
func (a *applierV3Corrupt) LeaseRevoke(lc *pb.LeaseRevokeRequest) (*pb.LeaseRevokeResponse, error) {
	return nil, ErrCorrupt
}

	plog = capnslog.NewPackageLogger("go.etcd.io/etcd", "wal")

	ErrMetadataConflict = errors.New("wal: conflicting metadata found")
	ErrFileNotFound     = errors.New("wal: file not found")
	ErrCRCMismatch      = errors.New("wal: crc mismatch")
	ErrSnapshotMismatch = errors.New("wal: snapshot mismatch")
	ErrSnapshotNotFound = errors.New("wal: snapshot not found")
	crcTable            = crc32.MakeTable(crc32.Castagnoli)
)

// WAL is a logical representation of the stable storage.
	decoder   *decoder       // decoder to decode records
	readClose func() error   // closer for decode reader

	mu      sync.Mutex
	enti    uint64   // index of the last entry saved to the wal
	encoder *encoder // encoder to encode records
}

// Create creates a WAL ready for appending records. The given metadata is
// recorded at the head of each WAL file, and can be retrieved with ReadAll.
func Create(lg *zap.Logger, dirpath string, metadata []byte) (*WAL, error) {
	if Exist(dirpath) {
		return nil, os.ErrExist
		}
		return nil, perr
	}
	if perr = fileutil.Fsync(pdir); perr != nil {
		if lg != nil {
			lg.Warn(
		}
		return nil, perr
	}
	if perr = pdir.Close(); perr != nil {
		if lg != nil {
			lg.Warn(
	return w, nil
}

func (w *WAL) cleanupWAL(lg *zap.Logger) {
	var err error
	if err = w.Close(); err != nil {
	defer w.mu.Unlock()

	rec := &walpb.Record{}
	decoder := w.decoder

	var match bool
		switch rec.Type {
		case entryType:
			e := mustUnmarshalEntry(rec.Data)
			if e.Index > w.start.Index {
				ents = append(ents[:e.Index-w.start.Index-1], e)
			}
			w.enti = e.Index

	return metadata, state, ents, err
}

// Verify reads through the given WAL and verifies that it is not corrupted.
// It creates a new decoder to read through the records of the given WAL.
// It does not conflict with any open WAL, but it is recommended not to
	if err = os.Rename(newTail.Name(), fpath); err != nil {
		return err
	}
	if err = fileutil.Fsync(w.dirFile); err != nil {
		return err
	}

	// reopen newTail with its new path so calls to Name() match the wal filename format
	newTail.Close()
}

func (w *WAL) sync() error {
	if w.encoder != nil {
		if err := w.encoder.flush(); err != nil {
			return err
	return err
}

// ReleaseLockTo releases the locks, which has smaller index than the given index
// except the largest one among them.
// For example, if WAL is holding lock 1,2,3,4,5,6, ReleaseLockTo(4) will release
