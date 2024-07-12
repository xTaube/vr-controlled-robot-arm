package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/quic-go/webtransport-go"
	"github.com/xTaube/vr-controlled-robot-arm/robot"
	"github.com/xTaube/vr-controlled-robot-arm/video"
	"go.bug.st/serial"
)

const BUFF_SIZE = 128

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFF_SIZE,
	WriteBufferSize: BUFF_SIZE,
}

func WebTransportControlRequestHandler(server *webtransport.Server) func(http.ResponseWriter, *http.Request) {
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

		log.Println("Initializing camera...")
		videoStream := video.InitVideoStream(
			os.Getenv("CAMERA_DEVICE_PATH"),
			video.Resoulution{Width: 1920, Height: 1080},
			video.FPS30,
			video.MJPEG,
			"rtsp://localhost:8554/video/feed",
		)
		defer videoStream.Stop()
		log.Println("Camera initialized")

		log.Printf("Initializing uart...")
		uart, err := robot.InitUart(
			os.Getenv("UART_PORT"),
			115200,
			serial.EvenParity,
			8,
			serial.OneStopBit,
		)
		if err != nil {
			log.Printf("Failed to initialize uart")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer uart.Close()
		log.Println("Uart initialized")

		commandHandler := InitCommandHandler(videoStream, uart)
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

func WebSocketControlRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading session...")
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket:%s\n", err)
		return
	}
	log.Println("Session upgraded to WebSocket")
	defer connection.Close()

	log.Println("Initializing camera...")
	videoStream := video.InitVideoStream(
		os.Getenv("CAMERA_DEVICE_PATH"),
		video.Resoulution{Width: 1920, Height: 1080},
		video.FPS30,
		video.MJPEG,
		"rtsp://localhost:8554/video/feed",
	)
	defer videoStream.Stop()
	log.Println("Camera initialized")

	log.Printf("Initializing uart...")
	uart, err := robot.InitUart(
		os.Getenv("UART_PORT"),
		115200,
		serial.EvenParity,
		8,
		serial.OneStopBit,
	)
	if err != nil {
		log.Printf("Failed to initialize uart")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uart.Close()
	log.Println("Uart initialized")

	commandHandler := InitCommandHandler(videoStream, uart)
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			break
		}

		log.Printf("Recived from connection: %s, type: %d", message, mt)
		output, err := commandHandler.Handle(CommandIdentifier((message[0] - 48)))
		if err != nil {
			log.Printf("Error occured: %s", err)
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %s", err)))
		} else {
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Output: %s", output)))
		}
	}
	log.Println("Session finished")
}
