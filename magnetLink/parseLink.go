package main 


import (
	"net/url"
	"strings"
	"errors"
	"encoding/hex"
	"encoding/base32"
)


func validXt(xt string) bool {

	if strings.HasPrefix(xt , "urn:btih:") == false{
		return false
	}
	if len(xt) -  len("urn:btih:") != 32  && len(xt) -  len("urn:btih:") !=  40{
		return false
	}
	return true
}

func parseXt(xt string) ([]byte ,bool) {
	xt = xt[len("urn:btih:"):]

	if len(xt) == 40{
		decode, err :=	hex.DecodeString(xt)
		if err != nil{
			return nil, false
		}
		return decode, true
	}
	if len(xt) == 32{
		decode, err := base32.StdEncoding.DecodeString(xt)
		if err != nil{
			return nil, false
		}
		return decode, true
	}
	return nil, false
}

func parseTracker(link string ) *Tracker {
	

	if strings.HasPrefix(link, "udp:") == true {
		if strings.HasSuffix(link, "/announce"){
			link = link[:len(link) - len("/announce")]
		}
		link = link[len("udp://"):]
		return  &Tracker{Schema: Udp, Endpoint: link } 
	}

	if strings.HasPrefix(link, "http") == true {
		return  &Tracker{Schema: Http, Endpoint: link } 
	}
	return nil
}

func ParseMagnetLink(link string) (*MagnetLink, error) {

	// magnet:?
	var res MagnetLink

	if strings.HasPrefix(link, "magnet:?") == false{
		return nil, errors.New("Invalid magnet Link")
	}
	link = link[len("magnet:?"):]

	// parsing queries
	v, err := url.ParseQuery(link)
	if err != nil {
		return nil, errors.New("Invalid magnet link, we can't parse queries")
	}

	// adding display name
	if v["dn"] != nil{
		res.Dn = v["dn"][0]
	}

	
	// validate hashInfo 

	if v["xt"] == nil || validXt(v["xt"][0]) == false {
		return nil, errors.New("Invalid magnet link, we can't parse the hash_info")
	}
	// parse hashInfo 
	info_hash , ok := parseXt(v["xt"][0])
	if !ok {
		return nil, errors.New("Invalid magnet link, we can't parse the hash_info")
	}
	res.Hash_info = [20]byte(info_hash)
	res.Xt =  v["xt"][0][len("urn:btih:"):]
	// validate the trackers

	if v["tr"] == nil{
		return nil, errors.New("Invalid magnet link, we don't have trackers ")
	}

	for _, tr := range(v["tr"]){
		tracker := parseTracker(tr)
		if tracker != nil{
			res.Trackers = append(res.Trackers,*tracker)
		}
	}

	return &res, nil
}
