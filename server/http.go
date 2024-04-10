package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quic-go/webtransport-go"
	"github.com/xTaube/vr-controlled-robot-arm/video"
)

const BUFF_SIZE = 128

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

		videoStream := video.InitVideoStream(
			os.Getenv("CAMERA_DEVICE_PATH"),
			video.Resoulution{Width: 1920, Height: 1080},
			video.FPS30,
			video.MJPEG,
			"rtsp://localhost:8554/video/feed",
		)
		defer videoStream.Stop()

		commandHandler := InitCommandHandler(videoStream)
		for {
			buf := make([]byte, BUFF_SIZE)
			log.Println("Waiting for message...")
			n, err := stream.Read(buf)
			if err != nil {
				break
			}
			log.Printf("Recived from stream %v: %s\n", stream.StreamID(), buf[:n])
			output, err := commandHandler.Handle(CommandIdentifier((buf[0] - 48)))
			if err != nil {
				log.Printf("Error occured: %s", err)
				stream.Write([]byte(fmt.Sprintf("Error: %s", err)))
			} else {
				stream.Write([]byte(fmt.Sprintf("Output: %s", output)))
			}
		}
		log.Println("Session finished")
	}
}
