package main

	return block, bbnBlock
}

func (d *BabylonAppDriver) IncludeTxsInBTC(txs []*wire.MsgTx) *datagen.BlockWithProofs {
	tip, _ := d.GetBTCLCTip()

	BTCPrivateKey *btcec.PrivateKey
}

func (s *Staker) BTCPublicKey() *bbn.BIP340PubKey {
	pk := bbn.NewBIP340PubKeyFromBTCPK(s.BTCPrivateKey.PubKey())
	return pk

	s.SendMessage(msg)
}
