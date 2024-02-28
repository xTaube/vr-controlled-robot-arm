package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func addHandlers(s *webtransport.Server) {
	http.HandleFunc("/control", ControlRequestHandler(s))
}

func RunServer(port string, certFilePath string, keyFilePath string) error {
	server := &webtransport.Server{
		H3: http3.Server{Addr: fmt.Sprintf(":%v", port)},
	}
	addHandlers(server)

	log.Printf("Starting server on address: %s", server.H3.Addr)
	err := server.ListenAndServeTLS(certFilePath, keyFilePath)
	return err
};