package server

import (
	"log"
	"net/http"

	"github.com/quic-go/webtransport-go"
)

const BUFF_SIZE = 1024


func ControlRequestHandler(server *webtransport.Server) func(http.ResponseWriter, *http.Request) {
	log.Println("ControlRequestHander registered")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upgrading session...")
		session, err := server.Upgrade(w, r)
		if err != nil {
			log.Printf("Failed to upgrade to WebTransport: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Session upgraded to WebTransport")

		log.Println("Accepting stream...")
		stream, err := session.AcceptStream(r.Context())
		if err != nil {
			log.Printf("Failed to accept stream: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Stream accepted")
		defer stream.Close()
		for {
			buf := make([]byte, BUFF_SIZE)
			log.Println("Waiting for message...")
			n, err := stream.Read(buf)
			if err != nil {
				break
			}
			log.Printf("Recived from stream %v: %s\n", stream.StreamID(), buf[:n])
		}
	}
}