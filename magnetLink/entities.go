package main 


import (

	"net"

)

type SchemaType int

const (
	Http = iota
	Udp = iota
)


type Tracker struct {

	Schema 		SchemaType // http or udp 
	Endpoint	string	// the host and the port for udp, and for http the hole endpoint where we should do Get request

}

type MagnetLink struct {
	
	Trackers	[]Tracker
	Dn		string
	Hash_info	[20]byte
	Xt 		string
}


type  ConnectReq struct {

	Protocol	uint64
	Action		uint32
	Trans		uint32
}


type  ConnectResp struct {
	
	Action		uint32 // action
	Trans		uint32 // transaction id
	ConId		uint64 // connection Id 
}


type AnnounceRequest struct {
	ConId		uint64 // connection Id 
	Action		uint32
	Trans		uint32 // transaction id
	Info_hash	[20]byte
	Peer_id		[20]byte
	Downloaded	uint64
	Left		uint64
	Upload		uint64
	Event		uint32
	IP		uint32
	Key		uint32
	Num_want	int32
	Port		uint16
}

type AnnounceRespHeaders struct {

	Action		uint32
	Trans		uint32 // transaction id
	Interval	uint32 
	Leechers	uint32 
	Seeders		uint32 
}

type Peer struct {
	IP   net.IP
	Port 	uint16
}


