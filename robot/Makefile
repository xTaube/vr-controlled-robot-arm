PROJECT_DIR=$(shell pwd)
BUILD_DIR=$(PROJECT_DIR)/build

RPI_BUILD_DIR=v-arm/robot-build
RPI_SCRIPTS_DIR=v-arm/vr-controlled-robot-arm/robot/scripts

remote-deploy:
	arduino-cli compile --build-path $(BUILD_DIR) --fqbn arduino:avr:uno $(PROJECT_DIR)
	ssh raspberry rm -rf $(RPI_BUILD_DIR)
	scp -pr $(BUILD_DIR) raspberry:$(RPI_BUILD_DIR) > /dev/null
	ssh raspberry $(RPI_SCRIPTS_DIR)/robot_remote_deploy.sh -p $(port)


list-available-boards:
	ssh raspberry $(ARDUINO_CLI_DIR)/arduino-cli board list
