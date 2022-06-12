package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
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
    curr += copy(hbuffer[curr:], hshake.Protocol)
	curr += copy(hbuffer[curr : ], make([]byte, 8))
	curr += copy(hbuffer[curr : ], hshake.InfoHash[:])
	curr += copy(hbuffer[curr : ], hshake.PeerID[:])
	return hbuffer
}

//parses a handshake from a stream
func ReadHandshake(r io.Reader) (*Handshake, error){
	hbuffer := make([]byte, 68)
	n, err := r.Read(hbuffer)
	if err != nil {
		panic(err)
	}
	if n != 68 {
		fmt.Println("size of recieved hanshake is unexpected")
	}
	h := Handshake{}
	lengthProtocol := hbuffer[0]
	if lengthProtocol != 19{
		fmt.Println("length of protocol field not as expected")
	}
	curr := 1
	extension := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	curr += copy([]byte(h.Protocol), hbuffer[curr:])
	curr += copy(extension[:], hbuffer[curr:])
	curr += copy(h.InfoHash[:], hbuffer[curr:])
	curr += copy(h.PeerID[:], hbuffer[curr:])
	return &h, nil
}

func ConnectToPeers(TR TrackerResponse){
	//idt this is complete
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

func NewHandshake(infoHashnew [20]byte, peerID [20]byte) (*Handshake){
	return &Handshake{
		Protocol: "Bittorent protocol",
		InfoHash: infoHashnew,
		PeerID: peerID,
	}
}