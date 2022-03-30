package main

import (
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

    TorrentFileTester(torrentinfo)
}
