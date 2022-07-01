TARGET = coordinate
BIN = /usr/bin/$(TARGET)

build:
	go mod tidy
	go build -o build/$(TARGET)

install:
	@install -Dm755 build/$(TARGET) $(BIN)
	@mkdir -p ~/.config/coordinate
	@cp assets/$(TARGET).yaml ~/.config/coordinate

clean:
	@rm -rf build

uninstall:
	@sudo rm -f $(BIN)
