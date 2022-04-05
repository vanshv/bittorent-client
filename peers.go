package main

import (
	"fmt"
	//"net"
	"net/http"
	"net/url"
	"strconv"
    "github.com/jackpal/bencode-go"
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

    return base.String(), nil
}

func makeGetReqeust(getrequest string)(*TrackerResponse){
    resp, err := http.Get(getrequest)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)

    TR := TrackerResponse{69420, nil}
    err = bencode.Unmarshal(resp.Body, &TR)
    if(err != nil){
        panic(err)
    }
    
    return &TR
}

type Peer struct{
    Port    uint16      `bencode:"port"`
    IP      string      `bencode:"ip"`
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
    Peers   []Peer  `bencode:"peers"`
}