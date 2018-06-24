CONF_DIR := /etc/eye
INSTALL_DIR := /opt/eye
CONF_FILE := eye.cfg
BIN_FILE := eye-agent
DEBIAN_PACKAGE := eye_1.0-0
TEMP_BASE_DIR := /tmp
DEBIAN_TEMP_DIR := $(TEMP_BASE_DIR)/$(DEBIAN_PACKAGE)

build:
	go build

install: build
	mkdir -p $(INSTALL_DIR)
	cp ./$(BIN_FILE) $(INSTALL_DIR)
	if [ ! -d $(CONF_DIR) ]; then mkdir $(CONF_DIR); fi
	if [ ! -e $(CONF_DIR)/$(CONF_FILE) ]; then cp ./$(CONF_FILE) $(CONF_DIR); fi
	cp ./eye.service /etc/systemd/system/eye.service
	systemctl enable eye.service
	systemctl start eye.service

uninstall:
	rm -Rf $(INSTALL_DIR)
	systemctl stop eye.service
	systemctl disable eye.service
	rm /etc/systemd/system/eye.service

package-for-debian: build
	if [ -d $(DEBIAN_TEMP_DIR) ]; then rm -R $(DEBIAN_TEMP_DIR); fi
	mkdir $(DEBIAN_TEMP_DIR)
	mkdir -p $(DEBIAN_TEMP_DIR)$(CONF_DIR)
	cp ./eye.cfg $(DEBIAN_TEMP_DIR)$(CONF_DIR)
	mkdir -p $(DEBIAN_TEMP_DIR)$(INSTALL_DIR)
	cp ./eye-agent $(DEBIAN_TEMP_DIR)$(INSTALL_DIR)
	mkdir -p $(DEBIAN_TEMP_DIR)/etc/systemd/system
	cp ./eye.service $(DEBIAN_TEMP_DIR)/etc/systemd/system
	mkdir $(DEBIAN_TEMP_DIR)/DEBIAN
	cp ./DEBIAN/* $(DEBIAN_TEMP_DIR)/DEBIAN
	dpkg-deb --build $(DEBIAN_TEMP_DIR)
	mv $(TEMP_BASE_DIR)/$(DEBIAN_PACKAGE).deb .
	rm -R $(DEBIAN_TEMP_DIR)

clean:
	rm -f ./$(BIN_FILE)
	rm -f ./$(DEBIAN_PACKAGE).deb
