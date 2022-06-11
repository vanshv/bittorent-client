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
	peerID   [20]byte
}

func NewClient(peer Peer, PeerID, InfoHash [20]byte) (*Client, error){
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

func (c *Client) SendUnchoke() error {
	msg := Message{ID :	MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := Message{ID :	MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	msg := FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
