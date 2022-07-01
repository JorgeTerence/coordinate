TARGET = coordinate
BIN = /usr/bin/$(TARGET)

build:
	go mod tidy
	go build -o build/$(TARGET)

install:
	@install -Dm755 build/$(TARGET) $(BIN)

config:
	@mkdir -p ~/.config/$(TARGET)
	@install -Dm666 assets/$(TARGET).yaml ~/.config/$(TARGET)/$(TARGET).yaml

clean:
	rm -rf build

uninstall:
	sudo rm -f $(BIN)
	rm -rf ~/.config/$(TARGET)
