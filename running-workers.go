package main

import (
	"bytes"
	"log"
	"time"
	"errors"
	"runtime"
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

func(t *TorrentData) Download() ([]byte, error){
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
	for (dledPieces < t.Length/t.PieceLength){
		res := <-results
		begin, end := t.calculateBounds(res.index)//this can use calculatePieceSize ig
		copy(buf[begin : end], res.buf)
		dledPieces++

		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		// log.Printf("Downloaded piece #%d\n", res.index)
		log.Printf("Downloaded %d out of %d or %d pieces from %d peers\n", 
		dledPieces, t.Length/t.PieceLength, len(t.PieceHashes), numWorkers)
	}

	//close(workQueue)

	return buf, nil
}

//makes connection with peer, pushes dled piece to results buffer
func (t *TorrentData) startDlWorker(peer Peer, workQueue chan *piecetoDl, results chan *pieceDled){
	c, err := NewClient(peer, t.MyPeerID, t.InfoHash)
	if err != nil{
		log.Printf("Could not perform handshake. Disconnecting %q", peer.IP)
		return
	}
	defer c.Conn.Close()
	// if c == nil{
	// 	return
	// }
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
func attemptPieceDownload(c *Client, p2dl *piecetoDl) ([]byte, error){
	log.Printf("Attempting piece #%d download from peer %q", p2dl.index, c.peer.IP)
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
				//log.Printf("Added piece %d to %q peer's backloglist", p2dl.index, c.peer.IP)
			}
		}

		err := state.readMessage()
		if err != nil {
			return nil, err			
		}
	}

	return state.buf, nil
}

func (t *TorrentData) calculatePieceSize(index int) (int) {
	if (len(t.PieceHashes) - 1) == index {
		remainder := t.Length - (len(t.PieceHashes) - 2)*t.PieceLength
		return remainder
	}
	return t.PieceLength
}

func (t *TorrentData) calculateBounds(index int) (int, int){
	piecesize := t.calculatePieceSize(index)
	begin := t.PieceLength*(index)
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