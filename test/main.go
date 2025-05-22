package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"math/rand"
	"bytes"
  	"encoding/binary"
	"errors"
	"encoding/hex"
	"encoding/base32"
)


func nop(e error){
	fmt.Println(">>>", e)
}

const PROTOCOL_ID uint64 = 0x41727101980


func parseHashInfo(hash_info string) ([]byte, error){

	if len(hash_info) == 40{
		data, err := hex.DecodeString(hash_info) 
		if err != nil{
			return nil, err
		}
		fmt.Println(len(data))
		if len(data) != 20{
			return nil, errors.New("Invalid hash_info")
		}
		return data, nil

	} else if len(hash_info) == 32{
		data, err := base32.StdEncoding.DecodeString(hash_info)
		if err != nil{
			return nil, err
		}
		if len(data) != 20{
			return nil, errors.New("Invalid hash_info")
		}
		return data, nil
	}

	return nil, errors.New("Invalid hash_info")
}

func main(){

	magnet := "magnet:?xt=urn:btih:E274E0F3389D2E6F3AD36484D2177AFCA3F98B03&dn=Alien+Invasion+Rise+of+the+Phoenix+2025+1080p+BluRay+HEVC+x265+5.1+BONE&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopentracker.io%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2F&tr=udp%3A%2F%2Fttk2.nbaonlineservice.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.darkness.services%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.0x7c0.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fttk2.nbaonlineservice.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fbandito.byterunner.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fevan.im%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce"




	if strings.HasPrefix(magnet , "magnet:?") == false{
		nop(errors.New("not a valid magnet link "))
	}

	// removing the schema
	magnet =  magnet[len("magnet:?"):]




	values, err := url.ParseQuery(magnet) 
	hash_info := values["xt"]
	if err != nil{
		nop(err)
	}
	hash_info_bytes, err := parseHashInfo(hash_info[0][len("urn:btih:"):]) 
	if err != nil{
		nop(err)
	}

	


	var trUdp []string

	// collecting trackers addresses
	for _, endpoint := range(values["tr"]){
		// udp://tracker.zer0day.to:1337/announce
		if  strings.HasPrefix(endpoint, "udp://"){
			endpoint = endpoint[len("udp://"):]

			if  strings.HasSuffix(endpoint,"/announce"){
				endpoint = endpoint[:len(endpoint) - len("/announce")]
			}
			endpoint = strings.Replace(endpoint, "/", "", -1)
			trUdp = append(trUdp , endpoint)
		}


	}

	for _, tr := range(trUdp){
		raddr, err := net.ResolveUDPAddr("udp", tr)
		if err != nil {
			return
		}
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			return
		}

		var connection_id  uint64
		err  = GetConnectionId(conn, &connection_id)
		if err != nil{
			nop(err)
		}
		err = GetPeers(conn, connection_id, hash_info_bytes)
		if err != nil{
			nop(err)
		}
		conn.Close()
	}


}

func GetConnectionId(con *net.UDPConn, connectionID  *uint64) error{
	var action int32 
	var transaction_id_req  int32 
	var transaction_id_resp	int32 


	action = 0
	transaction_id_req = rand.Int31()
	buf := new(bytes.Buffer)

	// requesting the connection id 
	err := binary.Write(buf,binary.BigEndian, PROTOCOL_ID)
	if err != nil{
		return err
	}
	err = binary.Write(buf,binary.BigEndian, action)
	if err != nil{
		return err
	}
	err = binary.Write(buf,binary.BigEndian, transaction_id_req)
	if err != nil{
		return err
	}


	n,err := con.Write(buf.Bytes())
	if err != nil || n != 16 {
		if err != nil{
			return err
		}
		return errors.New("we didn't sent the hole data ")
	}

	bf := make([]byte, 500)
	n,err  = con.Read(bf)
	bf = bf[:n]

	if err != nil || n <  16 {
		if err != nil {
			return err
		}
		return errors.New("we didn't receive the hole data ")
	}


    	reader := bytes.NewReader(bf)

    	if err := binary.Read(reader, binary.BigEndian, &action); err != nil {
        	return err
    	}
    	if err := binary.Read(reader, binary.BigEndian, &transaction_id_resp); err != nil {
        	return err
    	}

	if transaction_id_resp != transaction_id_req{

		return errors.New("transaction id of resp doesn´t equal the one in the request") 
	}
    	if err := binary.Read(reader, binary.BigEndian, connectionID); err != nil {
        	return err
    	}
	return  nil
}


type AnnounceRequest struct {
	ConnectionID   uint64      // 0: 64-bit connection ID
	Action         uint32      // 8: Action (1 = announce)
	TransactionID  uint32      // 12: Random transaction ID
	InfoHash       [20]byte    // 16: Info hash (torrent identifier)
	PeerID         [20]byte    // 36: Peer ID
	Downloaded     uint64      // 56: Bytes downloaded
	Left           uint64      // 64: Bytes left
	Uploaded       uint64      // 72: Bytes uploaded
	Event          uint32      // 80: Event (0 = none)
	IPAddress      uint32      // 84: IP address (default 0)
	Key            uint32      // 88: Random key
	NumWant        int32       // 92: Number of peers wanted (-1 = default)
	Port           uint16      // 96: Port the client is listening on
}



func (a *AnnounceRequest) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, a)
	return buf.Bytes(), err
}


func GetPeers(con *net.UDPConn, connectionID  uint64, hash_info []byte) error{
	req := AnnounceRequest{
		ConnectionID:  connectionID,
		Action:        1,
		TransactionID: rand.Uint32(),
		InfoHash:      [20]byte(hash_info),
		PeerID:        [20]byte{/* your peer ID here */},
		Downloaded:    0,
		Left:          100000,
		Uploaded:      0,
		Event:         2,        // started
		IPAddress:     0,
		Key:           rand.Uint32(),
		NumWant:       -1,
		Port:          6881,
	}	
	bt, err := req.ToBytes()
	if err != nil{
		return err
	}

	_,err = con.Write(bt)
	if err != nil{
		return err
	}
	
	buf := make([]byte,65507)
	n,err  := con.Read(buf)
	buf = buf[:n]
	
	if n < 20 {
		return errors.New("announce response too short")
	}	

	
	return nil
}
