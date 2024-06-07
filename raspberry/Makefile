CURRENT_DIR=$(shell pwd)

BINARY_NAME=$(CURRENT_DIR)/build/exec
SOURCE_FILE=$(CURRENT_DIR)/main.go


compile:
	go build -o $(BINARY_NAME) $(SOURCE_FILE)

clean:
	rm -f $(BINARY_NAME)

run:
	$(BINARY_NAME)
