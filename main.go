package main

import (
	"fmt"
	"os"
)

func main() {
    file, err := os.Open("debian.iso.torrent")
    //file, err := os.Open("3-gatsu.torrent")
    if err != nil {
        panic(err)
    }

    bencodetorrent, err := OpenBittorentFile(file)
    if err != nil {
        panic(err)
    }

    torrentinfo := bencodetorrent.toTorrentFile()

    var peerID = [20]byte{'w', 'e', 'l', 'c', 'o', 'm', 'e', 't', 'o', 'g', 'e', 't', 
    'r', 'e', 'q', 'u', 'e', 's', 't', 's'}
    myPort := 6881

    getrequest, err := torrentinfo.buildTrackerURL(peerID, uint16(myPort))
    if(err != nil){
        panic(err)
    }

    tr := makeGetReqeust(getrequest)
    torrentdata := TorrentData{
        InfoHash: torrentinfo.InfoHash,
        PieceHashes: torrentinfo.PieceHashes,
        PieceLength: torrentinfo.PieceLength,
        Length: torrentinfo.Length,
        TrackerResp: tr,
        MyPort: uint16(myPort),
        MyPeerID: peerID,
    }

    fmt.Println(tr.Interval)
    fmt.Println(tr.Peers)

    ConnectToPeers(tr)
    torrentdata.Download()

}

type TorrentData struct {
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    TrackerResp TrackerResponse
    MyPort      uint16
    MyPeerID    [20]byte
}