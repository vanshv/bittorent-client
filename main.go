package main

import (
	"os"
    "log"
)

func main() {
    inPath, err := os.Open("debian.torrent")
    if err != nil {
        log.Fatal(err)
    }

    torrentinfo, err := OpenTorrentFile(inPath)
    if err != nil{
        log.Fatal(err)
    }

    torrentdata, err := connectToTracker(*torrentinfo)
    if err != nil{
        log.Fatal(err)
    }

    buf, err := torrentdata.Download()
    if err != nil{
        log.Fatal(err)
    }

    outPath := "debian.iso"
    outFile, err := os.Create(outPath)
    if err != nil{
        log.Fatal(err)
    }
    defer outFile.Close()

    _, err = outFile.Write(buf)
    if err != nil{
        log.Fatal(err)
    }
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

func OpenTorrentFile(inPath *os.File) (*TorrentFile, error){
    bencodetorrent, err := OpenBittorentFile(inPath)
    if err != nil {
        return nil, err
    }
    torrentinfo, err := bencodetorrent.toTorrentFile()
    if err != nil{
        return nil, err
    }

    return &torrentinfo, nil
}

func connectToTracker(torrentfile TorrentFile) (*TorrentData, error){
    var peerID = [20]byte{'w', 'e', 'l', 'c', 'o', 'm', 'e', 't', 'o', 'g', 'e', 't', 
    'r', 'e', 'q', 'u', 'e', 's', 't', 's'}
    myPort := 6881

    getrequest, err := torrentfile.buildTrackerURL(peerID, uint16(myPort))
    if(err != nil){
        return nil, err
    }

    tr, err := makeGetReqeust(getrequest)
    if err != nil{
        return nil, err
    }

    torrentdata := TorrentData{
        InfoHash: torrentfile.InfoHash,
        PieceHashes: torrentfile.PieceHashes,
        PieceLength: torrentfile.PieceLength,
        Length: torrentfile.Length,
        TrackerResp: *tr,
        MyPort: uint16(myPort),
        MyPeerID: peerID,
    }

    //ConnectToPeers(*tr)

    return &torrentdata, nil
}