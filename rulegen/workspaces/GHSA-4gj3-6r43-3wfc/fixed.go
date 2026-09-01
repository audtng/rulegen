package main

}

// bitmap calculates the bitmap of which links in the shard are set.
func (s *shard) bitmap() ([]byte, error) {
	bm, err := bitfield.NewBitfield(s.size)
	if err != nil {
		return nil, err
	}
	for i := 0; i < s.size; i++ {
		if _, ok := s.children[i]; ok {
			bm.SetBit(i)
		}
	}
	return bm.Bytes(), nil
}

// serialize stores the concrete representation of this shard in the link system and
// returns a link to it.
func (s *shard) serialize(ls *ipld.LinkSystem) (ipld.Link, uint64, error) {
	bm, err := s.bitmap()
	if err != nil {
		return nil, 0, err
	}
	ufd, err := BuildUnixFS(func(b *Builder) {
		DataType(b, data.Data_HAMTShard)
		HashType(b, s.hasher)
		Data(b, bm)
		Fanout(b, uint64(s.size))
	})
	if err != nil {
	return len(fmt.Sprintf("%X", nd.FieldFanout().Must().Int()-1))
}

const maximumHamtWidth = 1 << 10

func bitField(nd data.UnixFSData) (bitfield.Bitfield, error) {
	fanout := int(nd.FieldFanout().Must().Int())
	if fanout > maximumHamtWidth {
		return nil, fmt.Errorf("hamt witdh (%d) exceed maximum allowed (%d)", fanout, maximumHamtWidth)
	}
	bf := bitfield.NewBitfield(fanout)
	bf.SetBytes(nd.FieldData().Must().Bytes())
	return bf, nil
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
	bf, err := bitField(data)
	if err != nil {
		return nil, err
	}
	return &_UnixFSHAMTShard{
		ctx:          ctx,
		_substrate:   substrate,
