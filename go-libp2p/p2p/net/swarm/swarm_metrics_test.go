package swarm

import (
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"

	"github.com/stretchr/testify/require"
)

func BenchmarkMetricsConnOpen(b *testing.B) {
	b.ReportAllocs()
	quicConnState := network.ConnectionState{Transport: "quic"}
	tcpConnState := network.ConnectionState{
		StreamMultiplexer: "yamux",
		Security:          "tls",
		Transport:         "tcp",
	}
	_, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(b, err)
	tr := NewMetricsTracer()
	for i := 0; i < b.N; i++ {
		switch i % 2 {
		case 0:
			tr.OpenedConnection(network.DirInbound, pub, quicConnState)
		case 1:
			tr.OpenedConnection(network.DirInbound, pub, tcpConnState)
		}
	}
}
