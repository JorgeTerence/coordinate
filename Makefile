TARGET = coordinate
BIN_DIR = /usr/bin/$(TARGET)
ASSETS_DIR = /usr/share/$(TARGET)
AUR_DIR = assets/package/aur

build:
	go mod tidy
	go build -o build/$(TARGET)

install:
	@mkdir -p $(ASSETS_DIR)
	@mkdir -p $(ASSETS_DIR)/web
	
	@sudo install -Dm755 build/$(TARGET) $(BIN_DIR)
	@sudo install -Dm644 web/* $(ASSETS_DIR)/web

clean: # Clean `build` and `pkg` artifacts
	@rm -rf build

	@printf "AUR [assets/package/aur] "
	@rm -rf \
		$(AUR_DIR)/$(TARGET) \
		$(AUR_DIR)/pkg \
		$(AUR_DIR)/src \
		$(AUR_DIR)/$(TARGET)*.pkg.tar.zst
	@echo "✔️"

erase: # Remove current installation
	@sudo rm -f $(BIN_DIR)
	@sudo rm -rf $(ASSETS_DIR)

pkg: # Generate and install Arch Linux package
	@echo "Erasing current installation 🧼 [make erase]"
	@make erase
	
	@echo "Generating package 📦 [makepkg]"
	@cd assets/package/aur; makepkg
	
	@echo "Installing package 💻"
	@sudo pacman -U $(AUR_DIR)/$(TARGET)*.pkg.tar.zst

	@echo "Don't forget to clean up 🧹 [make clean]"
