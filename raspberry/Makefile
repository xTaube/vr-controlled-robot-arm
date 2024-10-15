CURRENT_DIR=$(shell pwd)

BINARY_NAME=$(CURRENT_DIR)/build/exec
SOURCE_FILE=$(CURRENT_DIR)/main.go

RPI_PROJECT_LOCATION=v-arm
BINARY_RPI_LOCATION_TEMP=$(RPI_PROJECT_LOCATION)/server-build/temp/exec
RPI_SCRIPTS_LOCATION=$(RPI_PROJECT_LOCATION)/vr-controlled-robot-arm/raspberry/scripts
ROBOT_SERVER_SERVICE_NAME=robot-server


compile:
	go build -o $(BINARY_NAME) $(SOURCE_FILE)

clean:
	rm -f $(BINARY_NAME)

run:
	$(BINARY_NAME)

server-remote-deploy:
	env GOARCH=arm64 GOOS=linux go build -o $(BINARY_NAME) $(SOURCE_FILE)
	scp $(BINARY_NAME) raspberry:$(BINARY_RPI_LOCATION_TEMP)
	ssh raspberry sudo $(RPI_SCRIPTS_LOCATION)/update_and_restart_robot_server_service.sh

stop-server:
	ssh raspberry sudo systemctl stop $(ROBOT_SERVER_SERVICE_NAME).service

start-server:
	ssh raspberry sudo systemctl start $(ROBOT_SERVER_SERVICE_NAME).service
	ssh raspberry sudo journalctl -u $(ROBOT_SERVER_SERVICE_NAME).service -f

format:
	gofmt -s -w .
 
make raspberry-pull-main:
	ssh raspberry $(RPI_SCRIPTS_LOCATION)/pull_main.sh

make start-rtsp-server:
	ssh raspberry ./mediamtx