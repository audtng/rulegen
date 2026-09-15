package main

	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

/*
Encrypted overlay networks use IPsec in transport mode to encrypt and
authenticate the VXLAN UDP datagrams. This driver implements a bespoke control
plane which negotiates the security parameters for each peer-to-peer tunnel.

IPsec Terminology

 - ESP: IPSec Encapsulating Security Payload
 - SPI: Security Parameter Index
 - ICV: Integrity Check Value
 - SA: Security Association https://en.wikipedia.org/wiki/IPsec#Security_association


Developer documentation for Linux IPsec is rather sparse online. The following
slide deck provides a decent overview.
https://libreswan.org/wiki/images/e/e0/Netdev-0x12-ipsec-flow.pdf

The Linux IPsec stack is part of XFRM, the netlink packet transformation
interface.
https://man7.org/linux/man-pages/man8/ip-xfrm.8.html
*/

const (
	// Value used to mark outgoing packets which should have our IPsec
	// processing applied. It is also used as a label to identify XFRM
	// states (Security Associations) and policies (Security Policies)
	// programmed by us so we know which ones we can clean up without
	// disrupting other VPN connections on the system.
	mark = 0xD0C4E3

	pktExpansion = 26 // SPI(4) + SeqN(4) + IV(8) + PadLength(1) + NextHeader(1) + ICV(8)
)

	bidir
)

// Mark value for matching packets which should have our IPsec security policy
// applied.
var spMark = netlink.XfrmMark{Value: mark, Mask: 0xffffffff}

type key struct {
	value []byte
	return ""
}

// Security Parameter Indices for the IPsec flows between local node and a
// remote peer, which identify the Security Associations (XFRM states) to be
// applied when encrypting and decrypting packets.
type spi struct {
	forward int
	reverse int
	return b.String()
}

func (d *driver) checkEncryption(nid string, rIP net.IP, isLocal, add bool) error {
	logrus.Debugf("checkEncryption(%.7s, %v, %t)", nid, rIP, isLocal)

	n := d.network(nid)
	if n == nil || !n.secure {

	if add {
		for _, rIP := range nodes {
			if err := setupEncryption(lIP, aIP, rIP, d.secMap, d.keys); err != nil {
				logrus.Warnf("Failed to program network encryption between %s and %s: %v", lIP, rIP, err)
			}
		}
	return nil
}

// setupEncryption programs the encryption parameters for secure communication
// between the local node and a remote node.
func setupEncryption(localIP, advIP, remoteIP net.IP, em *encrMap, keys []*key) error {
	logrus.Debugf("Programming encryption between %s and %s", localIP, remoteIP)
	rIPs := remoteIP.String()

	indices := make([]*spi, 0, len(keys))

	for i, k := range keys {
		spis := &spi{buildSPI(advIP, remoteIP, k.tag), buildSPI(remoteIP, advIP, k.tag)}
		dir := reverse
	return nil
}

type matchVXLANFunc func(port, vni uint32) []string

// programVXLANRuleFunc returns a function which tries calling programWithMatch
// with the u32 match, falling back to the BPF match if installing u32 variant
// of the rules fails.
func programVXLANRuleFunc(programWithMatch func(matchVXLAN matchVXLANFunc, vni uint32, add bool) error) func(vni uint32, add bool) error {
	return func(vni uint32, add bool) error {
		if add {
			if err := programWithMatch(matchVXLANWithU32, vni, add); err != nil {
				// That didn't work. Maybe the xt_u32 module isn't available? Try again with xt_bpf.
				err2 := programWithMatch(matchVXLANWithBPF, vni, add)
				if err2 != nil {
					return multierror.Append(err, err2)
				}
			}
		} else {
			// Delete both flavours.
			err := programWithMatch(matchVXLANWithU32, vni, add)
			return multierror.Append(err, programWithMatch(matchVXLANWithBPF, vni, add)).ErrorOrNil()
		}
		return nil
	}
}

var programMangle = programVXLANRuleFunc(func(matchVXLAN matchVXLANFunc, vni uint32, add bool) error {
	var (
		m      = strconv.FormatUint(mark, 10)
		chain  = "OUTPUT"
		rule   = append(matchVXLAN(overlayutils.VXLANUDPPort(), vni), "-j", "MARK", "--set-mark", m)
		a      = iptables.Append
		action = "install"
	)

	// TODO IPv6 support
	iptable := iptables.GetIptable(iptables.IPv4)

	if !add {
		a = iptables.Delete
		action = "remove"
	}

	if err := iptable.ProgramRule(iptables.Mangle, chain, a, rule); err != nil {
		return fmt.Errorf("could not %s mangle rule: %w", action, err)
	}

	return nil
})

var programInput = programVXLANRuleFunc(func(matchVXLAN matchVXLANFunc, vni uint32, add bool) error {
	var (
		plainVxlan = matchVXLAN(overlayutils.VXLANUDPPort(), vni)
		chain      = "INPUT"
		msg        = "add"
	)

	rule := func(policy, jump string) []string {
		args := append([]string{"-m", "policy", "--dir", "in", "--pol", policy}, plainVxlan...)
		return append(args, "-j", jump)
	}

	// TODO IPv6 support
	iptable := iptables.GetIptable(iptables.IPv4)

	if !add {
		msg = "remove"
	}

	action := func(a iptables.Action) iptables.Action {
		if !add {
			return iptables.Delete
		}
		return a
	}

	// Accept incoming VXLAN datagrams for the VNI which were subjected to IPSec processing.
	// Append to the bottom of the chain to give administrator-configured rules precedence.
	if err := iptable.ProgramRule(iptables.Filter, chain, action(iptables.Append), rule("ipsec", "ACCEPT")); err != nil {
		return fmt.Errorf("could not %s input accept rule: %w", msg, err)
	}

	// Drop incoming VXLAN datagrams for the VNI which were received in cleartext.
	// Insert at the top of the chain so the packets are dropped even if an
	// administrator-configured rule exists which would otherwise unconditionally
	// accept incoming VXLAN traffic.
	if err := iptable.ProgramRule(iptables.Filter, chain, action(iptables.Insert), rule("none", "DROP")); err != nil {
		return fmt.Errorf("could not %s input drop rule: %w", msg, err)
	}

	return nil
})

func programSA(localIP, remoteIP net.IP, spi *spi, k *key, dir int, add bool) (fSA *netlink.XfrmState, rSA *netlink.XfrmState, err error) {
	var (
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.reverse,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: mark,
		}
		if add {
			rSA.Aead = buildAeadAlgo(k, spi.reverse)
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.forward,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: mark,
		}
		if add {
			fSA.Aead = buildAeadAlgo(k, spi.forward)
				Proto: netlink.XFRM_PROTO_ESP,
				Mode:  netlink.XFRM_MODE_TRANSPORT,
				Spi:   fSA.Spi,
				Reqid: mark,
			},
		},
	}
					Proto: netlink.XFRM_PROTO_ESP,
					Mode:  netlink.XFRM_MODE_TRANSPORT,
					Spi:   fSA2.Spi,
					Reqid: mark,
				},
			},
		}
		}
	}
	for _, sa := range saList {
		if sa.Reqid == mark {
			if err := nlh.XfrmStateDel(&sa); err != nil {
				logrus.Warnf("Failed to delete stale SA %s: %v", sa, err)
				continue
package overlay

import (
	"fmt"
	"strings"

	"golang.org/x/net/bpf"
)

// vniMatchBPF returns a BPF program suitable for passing to the iptables bpf
// match which matches on the VXAN Network ID of encapsulated packets. The
// program assumes that it will be used in a rule which only matches UDP
// datagrams.
func vniMatchBPF(vni uint32) []bpf.RawInstruction {
	asm, err := bpf.Assemble([]bpf.Instruction{
		bpf.LoadMemShift{Off: 0},                                    // ldx 4*([0] & 0xf) ; Load length of IPv4 header into X
		bpf.LoadIndirect{Off: 12, Size: 4},                          // ld [x + 12]       ; Load VXLAN ID (UDP header + 4 bytes) into A
		bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0xffffff00},        // and #0xffffff00   ; VXLAN ID is in top 24 bits
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: vni << 8, SkipTrue: 1}, // jeq ($vni << 8), match
		bpf.RetConstant{Val: 0},                                     // ret #0
		bpf.RetConstant{Val: ^uint32(0)},                            // match: ret #-1
	})
	// bpf.Assemble() only errors if an instruction is invalid. As the only variable
	// part of the program is an instruction value for which the entire range is
	// valid, whether the program can be successfully assembled is independent of
	// the input. Given that the only recourse is to fix this function and
	// recompile, there's little value in bubbling the error up to the caller.
	if err != nil {
		panic(err)
	}
	return asm
}

// marshalXTBPF marshals a BPF program into the "decimal" byte code format
// which is suitable for passing to the [iptables bpf match].
//
//	iptables -m bpf --bytecode
//
// [iptables bpf match]: https://ipset.netfilter.org/iptables-extensions.man.html#lbAH
func marshalXTBPF(prog []bpf.RawInstruction) string { //nolint:unused
	var b strings.Builder
	fmt.Fprintf(&b, "%d", len(prog))
	for _, ins := range prog {
		fmt.Fprintf(&b, ",%d %d %d %d", ins.Op, ins.Jt, ins.Jf, ins.K)
	}
	return b.String()
}
package overlay

import (
	"fmt"
	"strconv"
)

// matchVXLANWithU32 returns an iptables rule fragment which matches VXLAN
// datagrams with the given destination port and VXLAN Network ID utilizing the
// xt_u32 netfilter kernel module. The returned slice's backing array is
// guaranteed not to alias any other slice's.
func matchVXLANWithU32(port, vni uint32) []string {
	dport := strconv.FormatUint(uint64(port), 10)

	// The u32 expression language is documented in iptables-extensions(8).
	// https://ipset.netfilter.org/iptables-extensions.man.html#lbCK
	//
	// 0>>22&0x3C                ; Compute number of octets in IPv4 header
	//           @               ; Make this the new offset into the packet
	//                           ; (jump to start of UDP header)
	//            12&0xFFFFFF00  ; Read 32-bit value at offset 12 and mask off the bottom octet
	//                         = ; Test whether the value is equal to a constant
	//
	// A UDP header is eight octets long so offset 12 from the start of the
	// UDP header is four octets into the payload: the VNI field of the
	// VXLAN header.
	vniMatch := fmt.Sprintf("0>>22&0x3C@12&0xFFFFFF00=%d", int(vni)<<8)

	return []string{"-p", "udp", "--dport", dport, "-m", "u32", "--u32", vniMatch}
}
