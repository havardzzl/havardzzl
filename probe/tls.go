package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	host    = "tcloud.weizhipin.com:443"
	timeout = time.Second * 5
)

func ProbeTls() string {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		fmt.Printf("tcp dial error: %v\n", err)
		return err.Error()
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetLinger(0)
	}
	defer conn.Close()

	colonPos := strings.LastIndex(host, ":")
	if colonPos == -1 {
		colonPos = len(host)
	}
	hostname := host[:colonPos]

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostname,
	})
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = tlsConn.HandshakeContext(ctx)
	if err != nil {
		fmt.Printf("tls handshake error: %v\n", err)
		return err.Error()
	}
	minDurLeft := time.Hour * 500976
	var expDate time.Time
	for _, cert := range tlsConn.ConnectionState().PeerCertificates {
		valid := true
		valid = valid && !time.Now().After(cert.NotAfter)
		valid = valid && !time.Now().Before(cert.NotBefore)
		if !valid {
			return fmt.Sprintf("host %s cert not valid", host)
		}
		durLeft := time.Until(cert.NotAfter)
		if durLeft < minDurLeft {
			minDurLeft = durLeft
			expDate = cert.NotAfter
		}
	}
	return fmt.Sprintf("certificate is expiring after %v at %v", minDurLeft, expDate.Local().String())
}
