package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
    "encoding/binary"
)

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error){
    //builds the tracker URL, which will make GET request to torrentfile.announce
    base, err := url.Parse(t.Announce)
    if err != nil {
        return "", err
    }

    torrentLinkInfo := url.Values{//elements in values are strings mapped to an array of strings
        "info_hash":    []string{string(t.InfoHash[:])},
        "peer_id":      []string{string(peerID[:])},
        "port":         []string{strconv.Itoa(int(port))},
        "uploaded":     []string{"0"},
        "downloaded":   []string{"1"},
        "left":         []string{strconv.Itoa(t.Length)},
    }

    base.RawQuery = torrentLinkInfo.Encode()

    //fmt.Println(base.String())
    return base.String(), nil
}

func makeGetReqeust(getrequest string){
    resp, err := http.Get(getrequest)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)

    scanner := bufio.NewScanner(resp.Body)
    for i := 0; scanner.Scan() && i < 10; i++ {
        fmt.Println(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
    //hard part is to deencode bencoded data
}

type Peer struct{
    IP net.IP
    Port uint16
}

//store peer data in structs once it is deencoded
func Unmarshal(peerslist []byte)([]Peer, error){
    const peerSize = 6 
    numPeers := len(peerslist) / peerSize
    if len(peerslist)%peerSize != 0 {
        err := fmt.Errorf("Received malformed peers")
        return nil, err
    }
    peers := make([]Peer, numPeers)
    for i := 0; i < numPeers; i++ {
        offset := i * peerSize
        peers[i].IP = net.IP(peerslist[offset : offset + 4])
        peers[i].Port = binary.BigEndian.Uint16(peerslist[offset + 4 : offset + 6])
    }
    return peers, nil
}