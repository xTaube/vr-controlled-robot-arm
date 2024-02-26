
BINARY_NAME=build/exec
SOURCE_FILE=main.go


compile:
	go build -o $(BINARY_NAME) $(SOURCE_FILE)

clean:
	rm -f $(BINARY_NAME)

run:
	$(BINARY_NAME)
