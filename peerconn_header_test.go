package torrent

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/anacrolix/torrent/internal/alloclim"
	pp "github.com/anacrolix/torrent/peer_protocol"
)

// Test that sendChunk prefixes piece data with a 16 byte hex header
// containing the piece index and offset.
func TestSendChunkHexHeader(t *testing.T) {
	c := &PeerConn{}
	c.t = &Torrent{cl: &Client{config: &ClientConfig{SendHexPieceHeader: true}}}
	l := alloclim.Limiter{Max: 1}
	state := &peerRequestState{data: []byte{1, 2, 3}, allocReservation: l.Reserve(1)}
	r := Request{Index: 5, ChunkSpec: ChunkSpec{Begin: 7, Length: 3}}
	var got pp.Message
	c.sendChunk(r, func(m pp.Message) bool { got = m; return true }, state)

	var raw [8]byte
	binary.BigEndian.PutUint32(raw[:4], uint32(r.Index))
	binary.BigEndian.PutUint32(raw[4:], uint32(r.Begin))
	var expect [16]byte
	hex.Encode(expect[:], raw[:])
	expectedPiece := append(expect[:], state.data...)
	if string(got.Piece) != string(expectedPiece) {
		t.Fatalf("unexpected piece: %x", got.Piece)
	}
}
