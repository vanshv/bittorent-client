package main

import (
	"log"
	"time"
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

func(t *TorrentFile) run(){
	workQueue := make(chan * piecetoDl, len(t.PieceHashes))
	results := make(chan *pieceDled)
	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		//isn't piecesize calulated at start
		//piece size is piece size or + remainder
		workQueue <- &piecetoDl{index, hash, length}
	}
	//we push these structs in workQueue so that it can take them and dl them
	//thus workQueue stores undownloaded pieces.
	//once dled, it will send the piece to results buffer
	//results buffer capacity is 0, so it will probably send it for hash check.

	//we probably need the list of unchoked peers, which we do not have
	//however we haev the list of all peers in TrackerResponse struct
	for _, peer := range t.xyz{//list of peers, this is a placeholder
		go t.startDlWorker(peer, workQueue, results)
	}
	//this method of assigning pieces to peers
	//optimization ideas(depend on how we define peers)
	//-download from unchoked peers
	//-download from a peer which we're not already downloading from(at least come back to this peer after going to  all ohter peers)

	//basically, the whole file is stored in this random buffer
	buf := make([]byte, t.Length)
	dledPieces := 0
	for dledPieces < len(t.PieceHashes){
		res := <-results//read from results(?)
		begin, end := t.calculateBounds(res.index)//this can use calculatePieceSize ig
		copy(buf[begin : end], res.buf)
		dledPieces++
	}

	close(workQueue)
}

//makes connection with peer, pushes dled piece to results buffer
func (t *TorrentFile) startDlWorker(peer Peer, workQueue chan *piecetoDl, results chan *pieceDled){
	c, err := NewClient(peer, t.peerID, t.InfoHash)//t.peerID is our PeerID
	//ok so we need to create a class with torrentfile as its parent, 
	//we can add fields to it that come up further or in the code
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

		c.SendHave(p2dl.index)//what does this message mean?
		results <- &pieceDled{p2dl.index, buf}
	}
}

const MaxBlockSize = 16384
//pieces are made of blocks, in reality we request blocks
//from the peer
const MaxBacklog = 5
//number of unfulfilled requests a client can have
//no idea whats happening here
func attemptDownloadPiece(c *Client, p2dl *piecetoDl) ([]byte, error){
	state := pieceProgress{
		index: p2dl.index,
		client: c,
		buf: make([]byte, p2dl.length),
	}

	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline((time.Time{}))

	for state.downloaded < p2dl.length{
		//if unchoked, send requests until we have enough
		//unfulfilled requests
		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pw.length
				//Last bloack might be shorter than the typical block
				if pw.length - state.requested < blockSize{
					blockSize = pw.length - state.requested
				}

				err := c.SendRequest(p2dl.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}

			err := state.readMessage{
				if err != nil {
					return nil, err
				}
			}
	}

	return state.buf, nil
}

func (t *TorrentFile) calculatePieceSize(index int) (int) {
	//if statement clearly wrong lmao dumbass wtf
	if len(t.InfoHash) == index {
		remainder := t.Length - len(t.PieceHashes - 1)*t.PieceLength
		return remainder
	}
	return t.PieceLength
	
}//needs to be tested