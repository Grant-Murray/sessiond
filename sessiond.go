package main

import (
	"code.grantmurray.com/session"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"net"
	"net/http"
)

func main() {

	var err error

	flag.Parse()
	session.LoadConfig()

	var tlsc tls.Config
	tlsc.NextProtos = []string{"http/1.1"}
	tlsc.Certificates = make([]tls.Certificate, 1)
	if tlsc.Certificates[0], err = tls.X509KeyPair(session.Conf.HttpsCert, session.Conf.HttpsKey); err != nil {
		panic("Error attempting to use HttpsCert and HttpsKey to create a certificate: " + err.Error())
	}

	// verify the HttpsCert
	asn1, _ := pem.Decode(session.Conf.HttpsCert)
	cert, err := x509.ParseCertificate(asn1.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Error parsing cert data: %s\n", err))
	}

	var vo x509.VerifyOptions
	chains, err := cert.Verify(vo)
	if err != nil {
		panic(fmt.Sprintf("Error verifying: %s", err))
	}

	if len(chains) <= 0 {
		panic(fmt.Sprintf("Configured HttpsCert did not verify: %s\n", err))
	}

	// set up the handler and start the server
	http.HandleFunc("/session/", session.Handler)

	addr := fmt.Sprintf("%s:%d", session.Conf.HttpsHost, session.Conf.HttpsPort)
	tlsServer := &http.Server{Addr: addr, TLSConfig: &tlsc}
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Error while listening: " + err.Error())
	}

	tlsListener := tls.NewListener(conn, &tlsc)
	tlsServer.Serve(tlsListener)
}
