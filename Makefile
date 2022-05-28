TARGET = coordinate

build:
	go mod tidy
	go build -o build/$(TARGET)

install:
	mkdir -p /usr/share/$(TARGET)
	mkdir -p /usr/share/$(TARGET)/web
	
	install -Dm755 build/$(TARGET) /usr/bin/$(TARGET)
	install -Dm644 web/* /usr/share/$(TARGET)/web