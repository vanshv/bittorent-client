package main

import (
	"fmt"
	"io"
    "crypto/sha1"
    "github.com/jackpal/bencode-go"
    "log"
    "encoding/gob"
    "bytes"
)

type bencodeInfo struct {
    Pieces      string `bencode:"pieces"`
    PieceLength int    `bencode:"piece length"`
    Length      int    `bencode:"length"` // zero if there are multiple files in torrent
    Name        string `bencode:"name"`
}

type bencodeTorrent struct {
    Announce string      `bencode:"announce"`
    Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
}

func OpenBittorentFile(r io.Reader) (*bencodeTorrent, error){
    torrentinfo := bencodeTorrent{}
    err := bencode.Unmarshal(r, &torrentinfo)
    if (err != nil){
        return nil , err
    }
    return &torrentinfo, nil
}

//add error statement
func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
    torrentfile := TorrentFile{}
    torrentfile.Announce = bto.Announce
    torrentfile.Length = bto.Info.Length
    torrentfile.PieceLength = bto.Info.PieceLength
    torrentfile.Name = bto.Info.Name

    var buf bytes.Buffer
	err := bencode.Marshal(&buf, bto.Info)
    if(err != nil){
        panic(err)
    }
	torrentfile.InfoHash = sha1.Sum(buf.Bytes())

    num_pieces := torrentfile.Length/torrentfile.PieceLength
    torrentfile.PieceHashes = make([][20]byte, num_pieces*20)
    
    for i := 0; i < num_pieces; i++{
        for j := 0; j < 20; j++{
            torrentfile.PieceHashes[i][j] = bto.Info.Pieces[i*20 + j]
        }
    }
    
    return torrentfile, nil
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

func TorrentFileTester(torrentfile TorrentFile)(){
    fmt.Println(
    torrentfile.Announce, "announce \n",//
    torrentfile.PieceLength, "piecelength \n",
    torrentfile.Length, "length \n",
    torrentfile.Name, "name")
    fmt.Printf("%x infohash\n", torrentfile.InfoHash)
    for i := 0; i < torrentfile.Length/torrentfile.PieceLength; i++{
        for j := 0; j < 20; j++{
            if(i<5){
                fmt.Printf("%x ", torrentfile.PieceHashes[i][j])
            }
        }
        if(i<5){
            fmt.Println()
        }
    }
    fmt.Println("^ piecehashes")
}
