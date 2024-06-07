package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func addWebTransportHandlers(s *webtransport.Server) {
	http.HandleFunc("/control", WebTransportControlRequestHandler(s))
}

func RunWebTransportServer(port string, certFilePath string, keyFilePath string) error {
	server := &webtransport.Server{
		H3: http3.Server{Addr: fmt.Sprintf(":%v", port)},
	}
	addWebTransportHandlers(server)

	log.Printf("Starting server on address: %s", server.H3.Addr)
	err := server.ListenAndServeTLS(certFilePath, keyFilePath)
	return err
};

func addWebSocketHandlers() {
	http.HandleFunc("/control", WebSocketControlRequestHandler)
}

func RunWebSocketServer(port string) error {
	addWebSocketHandlers();
	log.Printf("Starting server on address: :%s", port)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		nil,
	)
	return err;
};	