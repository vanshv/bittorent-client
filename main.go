package main

import (
	"fmt"
    "os"
	"io"
    "crypto/sha1"
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

//debug PieceHashes and InfoHash in torrentfile struct and partition the functions in different folder
func main() {
    file, err := os.Open("debian.iso.torrent")
    if err != nil {
        panic(err)
    }

    bencodetorrent, err := OpenBittorentFile(file)
    if err != nil {
        panic(err)
    }

    torrentinfo := bencodetorrent.toTorrentFile()
    
    fmt.Println(
    torrentinfo.Announce, "announce \n",//
    torrentinfo.PieceLength, "piecelength \n",
    torrentinfo.Length, "length \n",
    torrentinfo.Name, "name")
    for i := 0; i<20; i++{
        fmt.Printf("%x ", torrentinfo.InfoHash[i])
    }
    // for i := 0; i < torrentinfo.PieceLength; i++{
    //     fmt.Println(torrentinfo.PieceHashes[i], " ")
    // }
    fmt.Println("piecehashes")
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

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile) {
    torrentfile := TorrentFile{}
    torrentfile.Announce = bto.Announce
    torrentfile.Length = bto.Info.Length
    torrentfile.PieceLength = bto.Info.PieceLength
    torrentfile.Name = bto.Info.Name

    h := sha1.New()
    h.Write([]byte(EncodeToBytes(bto.Info)))//convert BencodeInfo struct to []byte

    s := h.Sum(nil)
    var ret [20]byte
    copy(ret[:], s)

    torrentfile.InfoHash = ret

    usethis := torrentfile.Length/torrentfile.PieceLength
    piecehashes := make([][]byte, usethis*20)
 
    // for i := 0; i < torrentfile.PieceLength; i++{
    //     for j := 0; j < 20; j++{
    //         torrentfile.PieceHashes[i][j] = bto.Info.Pieces[i*20 + j]
    //     }
    // }
    fmt.Println(len(bto.Info.Pieces))
    for i := 0; i < usethis; i++{
        for j := 0; j < 20; j++{
            piecehashes[i][j] = bto.Info.Pieces[i*20 + j]
        }
    }
    for i := 0; i < usethis; i++{
        for j := 0; j<20; j++{
            fmt.Println(piecehashes[i][j], " ")
        }
        fmt.Println()
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