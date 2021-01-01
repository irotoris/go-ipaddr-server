package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var amazonCheckIpUrl = "http://checkip.amazonaws.com"

func main() {
	var addr = ":8080"
	http.HandleFunc("/", GetIpHandler)
	http.HandleFunc("/egress", GetEgressIpHandler)
	log.Printf("[START] ipaddr-server. port %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func GetIpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	ipaddr, err := GetIp(r)
	if err != nil {
		log.Printf("GetIp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := []byte(fmt.Sprintf("Internal Server Error - %v", err))
		w.Write(resp)
		return
	}

	resp := []byte(ipaddr + "\n")
	w.Write(resp)
}

func GetEgressIpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	CheckIpUrl := r.URL.Query().Get("CheckIpUrl")
	if CheckIpUrl == "" {
		CheckIpUrl = amazonCheckIpUrl
	}
	_, err := url.ParseRequestURI(CheckIpUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := []byte(fmt.Sprintf("Invalid \"CheckIpUrl\" Parameter - %v", err))
		w.Write(resp)
		return
	}

	egressIpaddr, err := GetEgressIp(CheckIpUrl)
	if err != nil {
		log.Printf("GetEgressIp: %v!", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := []byte(fmt.Sprintf("Internal Server Error - %v", err))
		w.Write(resp)
		return
	}

	resp := []byte(egressIpaddr)
	w.Write(resp)
}

func GetIp(r *http.Request) (string, error) {
	realip := r.Header.Get("X-REAL-IP")
	if realip != "" {
		return realip, nil
	}
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded, nil
	}
	ip, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip.IP.String(), nil
}

func GetEgressIp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
