package main

import (
    tf "/torrentfilemod/torrentfilemod" //does not work, something about environment variables
    "os"
)

//debug PieceHashes and InfoHash in torrentfile struct
func main() {
    file, err := os.Open("debian.iso.torrent")
    if err != nil {
        panic(err)
    }

    bencodetorrent, err := tf.OpenBittorentFile(file)
    if err != nil {
        panic(err)
    }

    torrentinfo := bencodetorrent.toTorrentFile()

    tf.TorrentFileTester(torrentinfo)
}