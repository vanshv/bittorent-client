package main

import (
	"fmt"
	"io"
)

type Handshake struct{
	Protocol string
	InfoHash [20]byte
	PeerID [20]byte
}

func (hshake *Handshake) Serialize() []byte{
	hbuffer := make([]byte, len(hshake.Protocol) + 49)
	hbuffer[0] = byte(len(hshake.Protocol))
	curr := 1
    curr += copy(hbuffer[curr : ], hshake.Protocol)
	curr += copy(hbuffer[curr : ], make([]byte, 8))
	curr += copy(hbuffer[curr : ], hshake.InfoHash[:])
	curr += copy(hbuffer[curr : ], hshake.PeerID[:])
	return hbuffer
}

//parses a handshake from a stream
func ReadHandshake(r io.Reader) (*Handshake, error){
	lengthbuf := make([]byte, 1)
	_,  err := io.ReadFull(r, lengthbuf)
	if err != nil {
		return nil, err
	}

	plen := int(lengthbuf[0])
	if(plen == 0){
		err := fmt.Errorf("String len cannot be 0")
		return nil, err
	}

	handshakebuf := make([]byte, 48 + plen)
	_, err = io.ReadFull(r, handshakebuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakebuf[plen+8:plen+8+20])
	copy(peerID[:], handshakebuf[plen+8+20:])

	h := Handshake{
		Protocol:     string(handshakebuf[0:plen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	return &h, nil
}

func NewHandshake(infoHashnew [20]byte, peerID [20]byte) (*Handshake){
	return &Handshake{
		Protocol: "BitTorrent protocol",
		InfoHash: infoHashnew,
		PeerID: peerID,
	}
}