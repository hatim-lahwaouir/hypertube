package main


import (
	"math/rand"
	"net" 
	"encoding/binary"
	"errors"
	"time"
	"log"
)

func sendUdp(c *net.UDPConn, n int, data any) error{
	buf := make([]byte, n)
	_ , err := binary.Encode(buf, binary.BigEndian, data)
	
	if err != nil{
		return err
	}
	nSent, err := c.Write(buf)
	
	if err != nil{
		return err
	}
	if n > nSent {
		return errors.New("no every thins was sent over udp")
	}

	return nil
}

func readUdpToBytes(c *net.UDPConn, n int) ([]byte, error){
	buf := make([]byte, 65507)
	nRead, err := c.Read(buf)
	if err != nil{
		return nil, err
	}
	if nRead <  n{
		return nil, errors.New("We didn't receive the hole data")
	}
	buf = buf[:nRead]
	return buf, nil 
}

func readUdp(c *net.UDPConn, n int, data any) error{
	buf := make([]byte, 65507)
	nRead, err := c.Read(buf)
	if err != nil{
		return err
	}
	if nRead <  n{
		return errors.New("We didn't receive the hole data")
	}

	_, err = binary.Decode(buf, binary.BigEndian, data)
	if err != nil{
		return err
	}

	return nil
}

func getConnection(c *net.UDPConn) (*ConnectResp, error){
	var conResp ConnectResp

	
	transaction := rand.Uint32()
	conReq :=  ConnectReq{Protocol: 0x41727101980, Action: 0, Trans : transaction}

	
	if err := sendUdp(c, 16, conReq); err != nil{
		return  nil, err
	}


	if err := readUdp(c, 16, &conResp); err != nil{
		return  nil, err
	}
	if conResp.Trans != conReq.Trans{
		return nil, errors.New("Invalid transaction id, the one we sent not equal the one we receive ")  
	}
	return &conResp, nil

}

func RequestAnnonce(c *net.UDPConn, conId uint64) ([]byte, error){
	AnReq := AnnounceRequest{
		ConId : conId,
		Num_want: 1,
		Trans: rand.Uint32(),
		Event: 0,
		Action: 1,
	}
	var AnResp  AnnounceRespHeaders

	if err := sendUdp(c,100, AnReq); err != nil{
		return nil, err
	}

	buf, err := readUdpToBytes(c, 20)
	if err != nil {
		return nil, err
	}
	_, err = binary.Decode(buf, binary.BigEndian, &AnResp)
	if err != nil {
		return nil, err
	}

	if AnReq.Trans != AnResp.Trans{
		return nil, errors.New("Invalid transaction id, the one we sent not equal the one we receive ")  
	}
	// remov the headers from buf
	buf = buf[20:]
	return buf , nil

}



func GetPeers(magnet *MagnetLink){


		var peers  []Peer
	for _, tr := range(magnet.Trackers){


		

		raddr, err := net.ResolveUDPAddr("udp", tr.Endpoint)	
		if err != nil{
			log.Println(err)
			continue
		}
		c, err := net.DialUDP("udp",nil, raddr)

		c.SetReadDeadline(time.Now().Add(1 * time.Second))
		
		if err != nil{
			log.Println(err)
			continue
		}
		
		conResp , err := getConnection(c)

		if err != nil{
			log.Println(err)

			continue
		}

		buf , err := RequestAnnonce(c, conResp.ConId)
		
		if err != nil {
			log.Println(err)
			return 
		}

		for i := 0; i + 6 <= len(buf); i += 6{
			ip := net.IPv4(buf[i], buf[i+1], buf[i+2], buf[i+3])
			port := binary.BigEndian.Uint16(buf[i+4 : i+6])
			peers = append(peers, Peer{IP: ip, Port: port})	 		

			log.Println("Peer successfuly found ", ip.String(), port)
		}


		defer c.Close()
	}


}
