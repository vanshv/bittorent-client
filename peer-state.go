package main

// import (
// )

type pieceProgress struct {
	index	int
	client	*Client
	buf 	[]byte //are these the files we've downloaded?
	downloaded int
	requested int //size of all the pieces requested from the peer
	backlog	int //number of pieces requested from the peer
}

//need to understand this later, when i get more clarity
func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read() //this call blocks
	if err != nil{
		return err
	}

	//keep alive message
	if msg == nil{
		return nil
	}
	
	switch msg.ID {
	case MsgUnchoke:
		state.client.Choked = false
	case MsgChoke:
		state.client.Choked = true
    case MsgHave: // peer is telling us we have the piece (?)
				//why do we need to decode the payload for this case?
        index, err := ParseHave(*msg)
		if(err != nil){
			return err
		}
        state.client.Bitfield.SetPiece(index)
	case MsgPiece: // they send the fucking piece for once
		n, err := ParsePiece(state.index, state.buf, msg)
		if err != nil{
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}