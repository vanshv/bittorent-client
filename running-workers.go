package main

import (
	"bytes"
	"log"
	"time"
	"errors"
)

//how to remove errors with t. or how to remove t from everywhere?

type piecetoDl struct{
	index	int
	hash	[20]byte
	length	int// what does this length value mean?
}

type pieceDled struct{
	index	int
	buf		[]byte
}//maybe size of this needs to be fixed?
//no, its not size
//maybe length of number of blocks?

func(t *TorrentData) Download(){
	workQueue := make(chan * piecetoDl, len(t.PieceHashes))
	results := make(chan *pieceDled)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &piecetoDl{index, hash, length}
	}
	for _, peer := range t.TrackerResp.Peers{//list of all peers
		go t.startDlWorker(peer, workQueue, results)
	}
	//basically, the whole file is stored in this random buffer
	buf := make([]byte, t.Length)
	dledPieces := 0
	for dledPieces < len(t.PieceHashes){
		res := <-results
		begin, end := t.calculateBounds(res.index)//this can use calculatePieceSize ig
		copy(buf[begin : end], res.buf)
		dledPieces++
	}

	close(workQueue)
}

//makes connection with peer, pushes dled piece to results buffer
func (t *TorrentData) startDlWorker(peer Peer, workQueue chan *piecetoDl, results chan *pieceDled){
	c, err := NewClient(peer, t.MyPeerID, t.InfoHash)
	if err != nil{
		panic(err)
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

		err = checkHash(*p2dl, buf)
		if err != nil{
			log.Printf("Piece #%d is corrupted", p2dl.index)
			workQueue <- p2dl
			continue
		}

		c.SendHave(p2dl.index)
		results <- &pieceDled{p2dl.index, buf}
	}
}

const MaxBlockSize = 16384
const MaxBacklog = 5
//number of unfulfilled requests a client can have
//no idea whats happening here
func attemptPieceDownload(c *Client, p2dl *piecetoDl) ([]byte, error){
	state := pieceProgress{
		index: p2dl.index,
		client: c,
		buf: make([]byte, p2dl.length),
	}

	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline((time.Time{}))

	for state.downloaded < p2dl.length{
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < p2dl.length{
				blockSize := MaxBlockSize
				//Last bloack might be shorter than the typical block
				if p2dl.length - state.requested < blockSize{
					blockSize = p2dl.length - state.requested
				}

				err := c.SendRequest(p2dl.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}

			err := state.readMessage()
			if err != nil {
				return nil, err
			}
		}
	}

	return state.buf, nil
}

func (t *TorrentData) calculatePieceSize(index int) (int) {
	if len(t.PieceHashes) == index {
		remainder := t.Length - (len(t.PieceHashes) - 1)*t.PieceLength
		return remainder
	}
	return t.PieceLength
}

func (t *TorrentData) calculateBounds(index int) (int, int){
	piecesize := t.calculatePieceSize(index)
	begin := t.PieceLength*(index-1)
	end := begin + piecesize
	return begin, end
}

func checkHash(p2dl piecetoDl, buf []byte) error {
	res := bytes.Compare(p2dl.hash[:], buf)
	if res == 0 {
		return errors.New("PieceHash not equal to downloaded piece")
	}
	return nil
}