SCRIPTS_DIR=v-arm/vr-controlled-robot-arm/robot/scripts
ARDUINO_CLI_DIR=local/bin

remote-deploy:
	ssh raspberry $(SCRIPTS_DIR)/robot_remote_deploy.sh -b $(target_branch) -p $(port)


list-available-boards:
	ssh raspberry $(ARDUINO_CLI_DIR)/arduino-cli board list
