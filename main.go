package main

import (
	//"fmt"
	"os"
)

//debug and InfoHash in torrentfile struct
//continue work on getting peers
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

    var port = [20]byte{'w', 'e', 'l', 'c', 'o', 'm', 'e', 't', 'o', 'g', 'e', 't', 
    'r', 'e', 'q', 'u', 'e', 's', 't', 's'}

    getrequest, err := torrentinfo.buildTrackerURL(port, 6881)
    if(err != nil){
        panic(err)
    }

    makeGetReqeust(getrequest)
}
