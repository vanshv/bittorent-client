package main

import (
	"log"
)

type piecetoDl struct{
	index	int
	hash	[20]byte
	length	int
}

type pieceDled struct{
	index	int
	buf		[]byte
}//maybe size of this needs to be fixed?

func (t *TorrentFile) run(){
	workQueue := make(chan * piecetoDl, len(t.PieceHashes))
	results := make(chan *pieceDled)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		//isn't piecesize calulated at start
		//piece size is piece size or + remainder
		workQueue <- &piecetoDl{index, hash, length}
	}

	for _, peer := range t.Length{//list of peers, this is a placeholder
		go t.startDlWorker(peer, workQueue, results)
	}

	
	buf := make([]byte, t.Length)
	dledPieces := 0
	for dledPieces < len(t.PieceHashes){
		res := <-results//read from results(?)
		begin, end := t.calculateBounds(res.index)
		copy(buf[begin : end], res.buf)
		dledPieces++
	}

	close(workQueue)
}

func (t *TorrentFile) startDlWorker(/*wtf is this*/peer Peer, workQueue chan *piecetoDl, results chan *pieceDled){
	c, err := client.New(peer, t, t.InfoHash)
	if err != nil{
		panic(err)
		return
	}
	defer c.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	for p2dl := range workQueue{
		if !c.Bitfield.HasPiece(p2dl.index){
			workQueue <- p2dl//go back lmao
			continue
		}
	
	
		buf, err := attemptPieceDownload(c, p2dl)
		if err != nil{
			log.Println("Exiting", err)
			workQueue <- p2dl
			return
		}

		err = checkHash(p2dl, buf)
		if err != nil{
			log.Printf("Piece #%d is corrupted", p2dl.index)
			workQueue <- p2dl
			continue
		}

		c.SendHave(p2dl.index)
		results <- &pieceDled{p2dl.index, buf}
	}
}

type pieceProgress struct {
	index	int
	client	*client.Client
	buf 	[]byte
	downloaded int
	requested int
	backlog	int
}

func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read()
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