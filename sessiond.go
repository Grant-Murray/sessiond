package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Grant-Murray/logdb"
	"github.com/Grant-Murray/session"
	"net"
	"net/http"
)

func main() {
	session.Configure()
	// session.log is created in plog-pg-schema.sql
	err := logdb.Initialize(session.Conf.DatabaseHandle, "", "INSERT INTO session.log (entered, msg, level) VALUES (now(), $1, $2)")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logdb: %s", err))
	}
	logdb.Info.Print("plogd server started")
	if !session.Conf.DebugVerbosely {
		logdb.Debug.Print(logdb.StopWriting)
	}

	var tlsc tls.Config
	tlsc.NextProtos = []string{"http/1.1"}
	tlsc.Certificates = make([]tls.Certificate, 1)
	if tlsc.Certificates[0], err = tls.X509KeyPair(session.Conf.HttpsCert, session.Conf.HttpsKey); err != nil {
		logdb.Attn.Fatalf("Error attempting to use HttpsCert and HttpsKey: %s", err.Error())
	} else {
		logdb.Debug.Printf("tls certificate was created")
	}

	http.HandleFunc("/session/", session.Handler)

	addr := fmt.Sprintf("%s:%d", session.Conf.HttpsHost, session.Conf.HttpsPort)
	tlsServer := &http.Server{Addr: addr, TLSConfig: &tlsc}
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		logdb.Attn.Fatalf("Error while listening: %s", err.Error())
	}

	tlsListener := tls.NewListener(conn, &tlsc)
	tlsServer.Serve(tlsListener)
}
