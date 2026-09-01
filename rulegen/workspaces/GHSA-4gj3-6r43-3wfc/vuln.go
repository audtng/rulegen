package main

}

// bitmap calculates the bitmap of which links in the shard are set.
func (s *shard) bitmap() []byte {
	bm := bitfield.NewBitfield(s.size)
	for i := 0; i < s.size; i++ {
		if _, ok := s.children[i]; ok {
			bm.SetBit(i)
		}
	}
	return bm.Bytes()
}

// serialize stores the concrete representation of this shard in the link system and
// returns a link to it.
func (s *shard) serialize(ls *ipld.LinkSystem) (ipld.Link, uint64, error) {
	ufd, err := BuildUnixFS(func(b *Builder) {
		DataType(b, data.Data_HAMTShard)
		HashType(b, s.hasher)
		Data(b, s.bitmap())
		Fanout(b, uint64(s.size))
	})
	if err != nil {
	return len(fmt.Sprintf("%X", nd.FieldFanout().Must().Int()-1))
}

func bitField(nd data.UnixFSData) bitfield.Bitfield {
	bf := bitfield.NewBitfield(int(nd.FieldFanout().Must().Int()))
	bf.SetBytes(nd.FieldData().Must().Bytes())
	return bf
}

func checkLogTwo(v int) error {
data/builder/dirshard.go |  15 +++-
go.mod                   |  54 ++++++-------
go.sum                   | 162 +++++++++++++++++++++++----------------
hamt/shardeddir_test.go  |  14 +++-
hamt/util.go             |   5 +-
5 files changed, 149 insertions(+), 101 deletions(-)
		return nil, err
	}
	shardCache := make(map[ipld.Link]*_UnixFSHAMTShard, substrate.FieldLinks().Length())
	bf := bitField(data)
	return &_UnixFSHAMTShard{
		ctx:          ctx,
		_substrate:   substrate,
