package main

import (
	"bytes"
	"fmt"
	//"log"
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

func NewClient(peer Peer, PeerID [20]byte, InfoHash [20]byte) (*Client, error){
	peerstr :=  net.JoinHostPort(peer.IP, strconv.Itoa(int(peer.Port)))
	conn, err := net.DialTimeout("tcp", peerstr, 30*time.Second)
	if err != nil{
		return nil, err
	}

	_, err = completeHandshake(conn, PeerID, InfoHash)
	if err != nil {
		conn.Close()
		return nil, err
	}

	//if we are not asking for bf, peer must be sending it on its own, meaning it sends its bitfield first thing
	//after the connection is set up. I don't understand how we don't keep a bitfield, since exhsanging bfs 
	//it feels like an important in the bittorent algorithm
	bf, err := recvBitField(conn)
	if err != nil{
		conn.Close()
		return nil, err
	}
	// log.Printf("Recieved bitfield from peer %q %q", peer.IP, bf)

	return &Client{
		Conn: conn,
		Choked : true,
		Bitfield:	bf,
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

func (c *Client) SendRequest(index, begin, length int) error {
	msg := FormatRequest(index, begin, length)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func recvBitField(conn net.Conn) (Bitfield, error){
	conn.SetDeadline(time.Now().Add(5*time.Second))
	defer conn.SetDeadline(time.Time{})
	
	msg, err := ReadMessage(conn)
	if err != nil{
		return nil, err
	}
	if(msg.ID != MsgBitfield){
		return nil, fmt.Errorf("expected msg.ID to be MsgBitfield and recieved %d", msg.ID)
	}

	return msg.Payload, nil
}

func completeHandshake(conn net.Conn, PeerID [20]byte, InfoHash [20]byte) (*Handshake, error) {
	conn.SetDeadline(time.Now().Add(5*time.Second))
	defer conn.SetDeadline(time.Time{})

	req := NewHandshake(InfoHash, PeerID)
	_, err := conn.Write(req.Serialize())
	if err != nil{
		return nil, err
	}

	req, err = ReadHandshake(conn)
	if err != nil{
		return nil, err
	}
	if (!bytes.Equal(req.InfoHash[:], InfoHash[:])){
		return nil, fmt.Errorf("Expected infohash %s but recieved %s", req.InfoHash[:], InfoHash[:])
	}

	return req, nil
}

func (c *Client) Read() (*Message, error) {
	msg, err := ReadMessage(c.Conn)
	return msg, err
}