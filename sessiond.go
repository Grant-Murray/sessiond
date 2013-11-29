package main

import (
  "crypto/tls"
  "flag"
  "fmt"
  "github.com/Grant-Murray/session"
  "net"
  "net/http"
)

func main() {

  var err error

  flag.Parse()

  var tlsc tls.Config
  tlsc.NextProtos = []string{"http/1.1"}
  tlsc.Certificates = make([]tls.Certificate, 1)
  if tlsc.Certificates[0], err = tls.X509KeyPair(session.Conf.HttpsCert, session.Conf.HttpsKey); err != nil {
    panic("Error attempting to use HttpsCert and HttpsKey to create a certificate: " + err.Error())
  }

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
