package main

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
    MsgChoke         messageID = 0
    MsgUnchoke       messageID = 1
    MsgInterested    messageID = 2
    MsgNotInterested messageID = 3
    MsgHave          messageID = 4
    MsgBitfield      messageID = 5
    MsgRequest       messageID = 6
    MsgPiece         messageID = 7
    MsgCancel        messageID = 8
)

type Message struct {
    ID      messageID
    Payload []byte
}

type Bitfield []byte

func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index/8
	offset := index%8
	return bf[byteIndex]>>(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index/8
	offset := index%8
	bf[byteIndex] |= 1 << (7 - offset)
}//each byte has 8 bits, this expression turns the specified bit in that byte on

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)//have to + 1 for messageID
	buf := make([]byte, 4 + length)// have to include + 4 for length length
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func ReadMessage(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil{
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil, nil
	}
	//keep-alive message
	//what is a keep-alive message?

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil{
		return nil, err
	}

	m := Message{
		ID:			messageID(messageBuf[0]),
		Payload:	messageBuf[1:],	
	}

	return &m, nil
}

func FormatHave(index int) *Message{
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}