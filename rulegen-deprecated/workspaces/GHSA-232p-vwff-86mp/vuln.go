package main

	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	r            = 0xD0C4E3
	pktExpansion = 26 // SPI(4) + SeqN(4) + IV(8) + PadLength(1) + NextHeader(1) + ICV(8)
)

	bidir
)

var spMark = netlink.XfrmMark{Value: uint32(r), Mask: 0xffffffff}

type key struct {
	value []byte
	return ""
}

type spi struct {
	forward int
	reverse int
	return b.String()
}

func (d *driver) checkEncryption(nid string, rIP net.IP, vxlanID uint32, isLocal, add bool) error {
	logrus.Debugf("checkEncryption(%.7s, %v, %d, %t)", nid, rIP, vxlanID, isLocal)

	n := d.network(nid)
	if n == nil || !n.secure {

	if add {
		for _, rIP := range nodes {
			if err := setupEncryption(lIP, aIP, rIP, vxlanID, d.secMap, d.keys); err != nil {
				logrus.Warnf("Failed to program network encryption between %s and %s: %v", lIP, rIP, err)
			}
		}
	return nil
}

func setupEncryption(localIP, advIP, remoteIP net.IP, vni uint32, em *encrMap, keys []*key) error {
	logrus.Debugf("Programming encryption for vxlan %d between %s and %s", vni, localIP, remoteIP)
	rIPs := remoteIP.String()

	indices := make([]*spi, 0, len(keys))

	err := programMangle(vni, true)
	if err != nil {
		logrus.Warn(err)
	}

	err = programInput(vni, true)
	if err != nil {
		logrus.Warn(err)
	}

	for i, k := range keys {
		spis := &spi{buildSPI(advIP, remoteIP, k.tag), buildSPI(remoteIP, advIP, k.tag)}
		dir := reverse
	return nil
}

func programMangle(vni uint32, add bool) (err error) {
	var (
		p      = strconv.FormatUint(uint64(overlayutils.VXLANUDPPort()), 10)
		c      = fmt.Sprintf("0>>22&0x3C@12&0xFFFFFF00=%d", int(vni)<<8)
		m      = strconv.FormatUint(uint64(r), 10)
		chain  = "OUTPUT"
		rule   = []string{"-p", "udp", "--dport", p, "-m", "u32", "--u32", c, "-j", "MARK", "--set-mark", m}
		a      = "-A"
		action = "install"
	)

	// TODO IPv6 support
	iptable := iptables.GetIptable(iptables.IPv4)

	if add == iptable.Exists(iptables.Mangle, chain, rule...) {
		return
	}

	if !add {
		a = "-D"
		action = "remove"
	}

	if err = iptable.RawCombinedOutput(append([]string{"-t", string(iptables.Mangle), a, chain}, rule...)...); err != nil {
		logrus.Warnf("could not %s mangle rule: %v", action, err)
	}

	return
}

func programInput(vni uint32, add bool) (err error) {
	var (
		port       = strconv.FormatUint(uint64(overlayutils.VXLANUDPPort()), 10)
		vniMatch   = fmt.Sprintf("0>>22&0x3C@12&0xFFFFFF00=%d", int(vni)<<8)
		plainVxlan = []string{"-p", "udp", "--dport", port, "-m", "u32", "--u32", vniMatch, "-j"}
		ipsecVxlan = append([]string{"-m", "policy", "--dir", "in", "--pol", "ipsec"}, plainVxlan...)
		block      = append(plainVxlan, "DROP")
		accept     = append(ipsecVxlan, "ACCEPT")
		chain      = "INPUT"
		action     = iptables.Append
		msg        = "add"
	)

	// TODO IPv6 support
	iptable := iptables.GetIptable(iptables.IPv4)

	if !add {
		action = iptables.Delete
		msg = "remove"
	}

	if err := iptable.ProgramRule(iptables.Filter, chain, action, accept); err != nil {
		logrus.Errorf("could not %s input rule: %v. Please do it manually.", msg, err)
	}

	if err := iptable.ProgramRule(iptables.Filter, chain, action, block); err != nil {
		logrus.Errorf("could not %s input rule: %v. Please do it manually.", msg, err)
	}

	return
}

func programSA(localIP, remoteIP net.IP, spi *spi, k *key, dir int, add bool) (fSA *netlink.XfrmState, rSA *netlink.XfrmState, err error) {
	var (
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.reverse,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: r,
		}
		if add {
			rSA.Aead = buildAeadAlgo(k, spi.reverse)
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.forward,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: r,
		}
		if add {
			fSA.Aead = buildAeadAlgo(k, spi.forward)
				Proto: netlink.XFRM_PROTO_ESP,
				Mode:  netlink.XFRM_MODE_TRANSPORT,
				Spi:   fSA.Spi,
				Reqid: r,
			},
		},
	}
					Proto: netlink.XFRM_PROTO_ESP,
					Mode:  netlink.XFRM_MODE_TRANSPORT,
					Spi:   fSA2.Spi,
					Reqid: r,
				},
			},
		}
		}
	}
	for _, sa := range saList {
		if sa.Reqid == r {
			if err := nlh.XfrmStateDel(&sa); err != nil {
				logrus.Warnf("Failed to delete stale SA %s: %v", sa, err)
				continue
