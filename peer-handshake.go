package main

import (
	"fmt"

)

type Handshake struct{
	lenProtocolID int
	Protocol string
	Extensions [8]byte
	InfoHash [20]byte
	PeerID [20]byte
}

func(hshake *Handshake) ExtendHand() []byte{
	buf := make([]byte, len(hshake.Protocol) + 49)
	buf[0] = byte(len(hshake.Protocol))
	curr := 1
	curr += copy(buf[curr : ], make([]byte, 8))
	curr += copy(buf[curr : ], hshake.InfoHash[:])
	curr += copy(buf[curr : ], hshake.PeerID[:])
	return buf
}

func(buf []byte) RecieveHand() (*Handshake){

}