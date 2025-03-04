//go:build linux
// +build linux

package aghnet

import (
	"io"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const arpAOutputWrt = `
IP address    HW type     Flags       HW address            Mask     Device
192.168.1.2   0x1         0x2         ab:cd:ef:ab:cd:ef     *        wan
::ffff:ffff   0x1         0x2         ef:cd:ab:ef:cd:ab     *        wan`

const arpAOutput = `
? (192.168.1.2) at ab:cd:ef:ab:cd:ef on en0 ifscope [ethernet]
? (::ffff:ffff) at ef:cd:ab:ef:cd:ab on em0 expires in 100 seconds [ethernet]`

const ipNeighOutput = `
192.168.1.2 dev enp0s3 lladdr ab:cd:ef:ab:cd:ef DELAY
::ffff:ffff dev enp0s3 lladdr ef:cd:ab:ef:cd:ab router STALE`

var wantNeighs = []Neighbor{{
	IP:  net.IPv4(192, 168, 1, 2),
	MAC: net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	IP:  net.ParseIP("::ffff:ffff"),
	MAC: net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}}

func TestFSysARPDB(t *testing.T) {
	a := &fsysARPDB{
		ns: &neighs{
			mu: &sync.RWMutex{},
			ns: make([]Neighbor, 0),
		},
		fsys:     testdata,
		filename: "proc_net_arp",
	}

	err := a.Refresh()
	require.NoError(t, err)

	ns := a.Neighbors()
	assert.Equal(t, wantNeighs, ns)
}

func TestCmdARPDB_arpawrt(t *testing.T) {
	a := &cmdARPDB{
		parse:  parseArpAWrt,
		runcmd: func() (r io.Reader, err error) { return strings.NewReader(arpAOutputWrt), nil },
		ns: &neighs{
			mu: &sync.RWMutex{},
			ns: make([]Neighbor, 0),
		},
	}

	err := a.Refresh()
	require.NoError(t, err)

	assert.Equal(t, wantNeighs, a.Neighbors())
}

func TestCmdARPDB_ipneigh(t *testing.T) {
	a := &cmdARPDB{
		parse:  parseIPNeigh,
		runcmd: func() (r io.Reader, err error) { return strings.NewReader(ipNeighOutput), nil },
		ns: &neighs{
			mu: &sync.RWMutex{},
			ns: make([]Neighbor, 0),
		},
	}
	err := a.Refresh()
	require.NoError(t, err)

	assert.Equal(t, wantNeighs, a.Neighbors())
}
