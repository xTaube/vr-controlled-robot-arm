package server

import (
	"errors"
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
	log.Println("ControlRequestHander registered.")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upgrading session...")
		session, err := server.Upgrade(w, r)
		if err != nil {
			log.Printf("Failed to upgrade to WebTransport: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Session upgraded to WebTransport.")

		log.Println("Accepting stream...")
		stream, err := session.AcceptStream(r.Context())
		if err != nil {
			log.Printf("Failed to accept stream: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Stream accepted.")
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
		log.Println("Camera initialized.")

		log.Println("Initializing robot arm...")
		robot, err := robot.InitRobot(
			robot.UartConfig{
				PortName: os.Getenv("UART_PORT"),
				Parity:   serial.EvenParity,
				StopBits: serial.OneStopBit,
				BaudRate: 115200,
				DataBits: 8,
			},
		)
		if err != nil {
			log.Printf("Error initializing robot arm: %s.\n", err)
			return
		}
		defer robot.ShutDown()
		log.Println("Robot arm initialized.")

		commandHandler := InitCommandHandler(videoStream, robot)
		for {
			buf := make([]byte, BUFF_SIZE)
			log.Println("Waiting for message...")
			n, err := stream.Read(buf)
			if err != nil {
				break
			}
			command_id, args := ParseRequestArguments(string(buf[:n]))
			output, err := commandHandler.Handle(command_id, args)
			if err != nil {
				log.Printf("Error occured: %s.", err)
				webtransportHandlerError(stream, err)
			} else {
				stream.Write([]byte(fmt.Sprintf("Output: %s", output)))
			}
		}
		log.Println("Session finished")
	}
}

func webtransportHandlerError(stream webtransport.Stream, err error) {
	if errors.Is(err, &CommandNotFound{}) {
		stream.Write(ErrorResponse(UNKNOWN_COMMAND, err))
	} else {
		stream.Write(ErrorResponse(UNKNOWN_ERROR, err))
	}
}

func WebSocketControlRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading session...")
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket:%s.\n", err)
		return
	}
	log.Println("Session upgraded to WebSocket.")
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
	log.Println("Camera initialized.")

	log.Println("Initializing robot arm...")
	robot, err := robot.InitRobot(
		robot.UartConfig{
			PortName: os.Getenv("UART_PORT"),
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
			BaudRate: 115200,
			DataBits: 8,
		},
	)
	if err != nil {
		log.Printf("Error initializing robot arm: %s.\n", err)
		return
	}
	defer robot.ShutDown()
	log.Println("Robot arm initialized.")

	commandHandler := InitCommandHandler(videoStream, robot)
	for {
		_, request, err := connection.ReadMessage()
		if err != nil {
			break
		}

		command_id, args := ParseRequestArguments(string(request))
		result, err := commandHandler.Handle(command_id, args)
		if err != nil {
			log.Printf("Error occured: %s", err)
			websocketHandleError(connection, err)
		} else {
			connection.WriteMessage(websocket.TextMessage, result)
		}
	}
	log.Println("Session finished")
}


func websocketHandleError(connection *websocket.Conn, err error) {
	if errors.Is(err, &CommandNotFound{}) {
		connection.WriteMessage(websocket.TextMessage, ErrorResponse(UNKNOWN_COMMAND, err))
	} else {
		connection.WriteMessage(websocket.TextMessage, ErrorResponse(UNKNOWN_ERROR, err))
	}
}