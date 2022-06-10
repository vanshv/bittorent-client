package main

// import (
// )

type pieceProgress struct {
	index	int
	client	*Client
	buf 	[]byte //are these the files we've downloaded?
	downloaded int
	requested int
	backlog	int //what is backlog?
}

//need to understand this later, when i get more clarity
func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read() //this call blocks
	switch msg.ID{
	case message.MsgUnchoke:
		state.client.Choked = false
	case message.MsgChoke:
		state.client.Choked = true
    case message.MsgHave:
        index, err := message.ParseHave(msg)
        state.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		state.downloaded += n
		state.backlog--
	}
}