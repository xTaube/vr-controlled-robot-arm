SCRIPTS_DIR=vr-controlled-robot-arm/robot/scripts
ARDUINO_CLI_DIR=arduino-cli/bin

remote-deploy:
	ssh raspberry $(SCRIPTS_DIR)/remote_deploy.sh -b $(target_branch) -p $(port)


list-available-boards:
	ssh raspberry $(ARDUINO_CLI_DIR)/arduino-cli board list
