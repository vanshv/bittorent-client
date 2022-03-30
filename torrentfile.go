package main

import (
	"fmt"
	"io"
    //"crypto/sha1"
    "github.com/jackpal/bencode-go"
    "bytes"
    "log"
    "encoding/gob"
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
    fmt.Println(torrentinfo.Info.PieceLength)
    if (err != nil){
        return nil , err
    }
    return &torrentinfo, nil
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile) {
    torrentfile := TorrentFile{}
    torrentfile.Announce = bto.Announce
    torrentfile.Length = bto.Info.Length
    torrentfile.PieceLength = bto.Info.PieceLength
    torrentfile.Name = bto.Info.Name

    // the dict is encoded in bencode form ( not bto.Info struct)

    num_pieces := torrentfile.Length/torrentfile.PieceLength
    torrentfile.PieceHashes = make([][20]byte, num_pieces*20)
    
    for i := 0; i < num_pieces; i++{
        for j := 0; j < 20; j++{
            torrentfile.PieceHashes[i][j] = bto.Info.Pieces[i*20 + j]
            if(i < 5){
                fmt.Printf("%x ", torrentfile.PieceHashes[i][j])
            }
        }
        if(i<5){
            fmt.Println()
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

func TorrentFileTester(torrentfile TorrentFile)(){
    fmt.Println(
    torrentfile.Announce, "announce \n",//
    torrentfile.PieceLength, "piecelength \n",
    torrentfile.Length, "length \n",
    torrentfile.Name, "name")
    for i := 0; i<20; i++{
        fmt.Printf("%x ", torrentfile.InfoHash[i])
    }
    // for i := 0; i < torrentinfo.PieceLength; i++{
    //     fmt.Println(torrentinfo.PieceHashes[i], " ")
    // }
    fmt.Println("piecehashes")
}
