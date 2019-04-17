package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"socks5-proxy-dende/utils/config"
	"socks5-proxy-dende/utils/logger"
	"strconv"
)

func main() {
	logger.Init()
	client := config.ClientConfig()
	monitorClient(client)
}

func monitorClient(client *config.Cconfig) {
	go func(client *config.Cconfig) {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(client.DefaultPort))
		if err != nil {
			logger.Error(fmt.Sprintf("監聽本地端口：%v,失敗:%+v", client.DefaultPort, err))
			return
		}

		for {
			conn, err := lis.Accept()
			if err != nil {
				logger.Error(fmt.Sprintf("接受請求信息 port:%v,失敗:%+v", client.DefaultPort, err))
				continue
			}
			go handle(conn, client.ServerIp+":"+strconv.Itoa(client.ServerPort))
		}
	}(client)
	select {}
}

func handle(sconn net.Conn, remoteip string) {
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		logger.Error(err)
		return
	}
	certBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		logger.Error("Unable to read cert.pem:", err)
		return
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		logger.Error("failed to parse root certificate:", err)
		return
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	dconn, err := tls.Dial("tcp", remoteip, conf)
	if err != nil {
		logger.Error(fmt.Sprintf("连接:%v,失败:%+v", remoteip, err))
		return
	}
	ExitChan := make(chan bool, 1)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	<-ExitChan

	dconn.Close()

}
