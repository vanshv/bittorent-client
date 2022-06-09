package main

import (
	"net"
	"strconv"
	"time"
)

// A Client is a TCP connection with a peer
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield Bitfield
	peer     Peer
	infoHash [20]byte
	peerID   int//changed this from [20]byte to int, check if any problems
}

func NewClient(peer Peer, PeerID int, InfoHash [20]byte) (*Client, error){
	connect :=  peer.IP + ":" + strconv.Itoa(int(peer.Port))
	conn, err := net.DialTimeout("tcp", connect, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: conn,
		Choked : true,
		Bitfield: getbitfield(conn),
		peer: peer,
		infoHash: InfoHash,
		peerID: PeerID,
	}, nil
}
