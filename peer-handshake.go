package main

import (
	"fmt"
	"net"
	"time"
	"strconv"
)

type Handshake struct{
	lenProtocolID int
	Protocol string
	Extensions [8]byte
	InfoHash [20]byte
	PeerID [20]byte
}

func ExtendHand(hshake *Handshake) []byte{
	buf := make([]byte, len(hshake.Protocol) + 49)
	buf[0] = byte(len(hshake.Protocol))
	curr := 1
	curr += copy(buf[curr : ], make([]byte, 8))
	curr += copy(buf[curr : ], hshake.InfoHash[:])
	curr += copy(buf[curr : ], hshake.PeerID[:])
	return buf
}

// func RecieveHand(buf []byte) (*Handshake){

// }

func ConnectToPeers(TR *TrackerResponse){
	connect :=  TR.Peers[0].IP + ":" + strconv.Itoa(int(TR.Peers[0].Port))
	conn, err := net.DialTimeout("tcp", connect, 5*time.Second)
    if err != nil {
    panic(err)
    }
	var str []byte

	conn.Read(str)
	fmt.Println(str)
	fmt.Println(conn.LocalAddr())
}