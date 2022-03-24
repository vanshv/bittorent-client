package main

import (
	"fmt"
	"io"
    "crypto/sha1"
    "github.com/jackpal/bencode-go"
    "bytes"
    "log"
    "encoding/gob"
)

type BencodeInfo struct{
    piece_length int
    length int
    name string
    pieces string
}

type BencodeTorrent struct{
    announce string
    info BencodeInfo
}

func main() {
    //write io to check if the code works and then continue to url tracker
}

func OpenBittorentFile(r io.Reader) (*BencodeTorrent, error){
    torrentinfo := BencodeTorrent{}
    err := bencode.Unmarshal(r, &torrentinfo)
    if (err != nil){
        return nil , err
    }
    return &torrentinfo, nil
}

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
}

func (bto BencodeTorrent) toTorrentFile() (TorrentFile) {
    torrentfile := TorrentFile{}
    torrentfile.Announce = bto.announce
    torrentfile.Length = bto.info.length
    torrentfile.PieceLength = bto.info.piece_length
    torrentfile.Name = bto.info.name

    h := sha1.New()
    h.Write([]byte(EncodeToBytes(bto.info)))//convert BencodeInfo struct to []byte

    s := h.Sum(nil)
    var ret [20]byte
    copy(ret[:], s)

    torrentfile.InfoHash = ret
    

    for i := 0; i < torrentfile.PieceLength; i++{
        for j := 0; j < 20; j++{
            torrentfile.PieceHashes[i][j] = bto.info.pieces[i*20 + j]
        }
    }

    return torrentfile
}

func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}