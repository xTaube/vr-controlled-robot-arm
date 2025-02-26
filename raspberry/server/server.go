package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/xTaube/vr-controlled-robot-arm/robot"
	"github.com/xTaube/vr-controlled-robot-arm/video"
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
}

func addWebSocketHandlers(robot *robot.Robot, video0 *video.VideoStream, video1 *video.VideoStream) {
	http.HandleFunc("/control", WebSocketControlRequestHandler(robot, video0, video1))
}

func RunWebSocketServer(port string, robot *robot.Robot, video0 *video.VideoStream, video1 *video.VideoStream) error {
	addWebSocketHandlers(robot, video0, video1)
	log.Printf("Starting server on address: :%s", port)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		nil,
	)
	return err
}
