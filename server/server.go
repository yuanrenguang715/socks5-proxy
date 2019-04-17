package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"socks5-proxy-dende/utils/config"
	"socks5-proxy-dende/utils/logger"
	"strconv"
	"time"
)

func main() {
	logger.Init()
	server := config.ServerConfig()
	start(server)
}

func start(server *config.Sconfig) {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		logger.Error(err)
		return
	}
	certBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		logger.Error("Unable to read cert.pem:", err)
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		logger.Error("failed to parse root certificate:", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	s, err := tls.Listen("tcp", ":"+strconv.Itoa(server.DefaultPort), config)
	if err != nil {
		logger.Error(err)
		return
	}

	for {
		c, err := s.Accept()
		if err != nil {
			logger.Error(err)
			continue
		}
		go proxy(c)
	}
}

func proxy(client net.Conn) {
	defer client.Close()
	var b [1024]byte

	n, err := client.Read(b[:])
	if err != nil {
		logger.Error(err)
		return
	}
	var addr string
	//只接受sock5代理
	if b[0] == 0x05 {
		//回应确认代理
		client.Write([]byte{0x05, 0x00})

		n, err = client.Read(b[:])
		if err != nil {
			logger.Error(err)
			return
		}
		switch b[3] {
		case 0x01:
			//解析代理ip
			type sockIP struct {
				A, B, C, D byte
				PORT       uint16
			}
			sip := sockIP{}
			if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
				logger.Error("请求解析错误")
				return
			}
			addr = fmt.Sprintf("%d.%d.%d.%d:%d", sip.A, sip.B, sip.C, sip.D, sip.PORT)
		case 0x03:
			//解析代理域名
			host := string(b[5 : n-2])
			var port uint16
			err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
			if err != nil {
				logger.Error(err)
				return
			}
			addr = fmt.Sprintf("%s:%d", host, port)
		}

		server, err := net.DialTimeout("tcp", addr, time.Second*3)
		if err != nil {
			logger.Error(err)
			return
		}
		defer server.Close()
		//回复确定代理成功
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		//转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}
}
