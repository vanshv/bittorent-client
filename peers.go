package main

import (
	"time"
	"net/http"
	"net/url"
	"strconv"
	"github.com/jackpal/bencode-go"
)

func (t TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error){
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

func makeGetReqeust(getrequest string)(*TrackerResponse, error){
    conn := &http.Client{Timeout: 15 * time.Second}
    resp, err := conn.Get(getrequest)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    TR := TrackerResponse{69420, nil}
    err = bencode.Unmarshal(resp.Body, &TR)
    //i don't understand completely, but as far as I can tell we don't need to pass the TR to 
    //unmarshal into the struct, we just need to define it and tag the variable names
    if(err != nil){
        return nil, err
    }
    
    return &TR, nil
}

type Peer struct{
    Port    uint16      `bencode:"port"`
    IP      string      `bencode:"ip"`
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
    Peers   []Peer  `bencode:"peers"`
}