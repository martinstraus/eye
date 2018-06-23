CONF_DIR := /etc/eye
INSTALL_DIR := /opt/eye
CONF_FILE := eye.cfg
BIN_FILE := eye-agent

build:
	go build

install: build
	mkdir -p $(INSTALL_DIR)
	cp ./$(BIN_FILE) $(INSTALL_DIR)
	if [ ! -d $(CONF_DIR) ]; then mkdir $(CONF_DIR); fi
	if [ -e $(CONF_DIR)/$(CONF_FILE) ]; then cp ./$(CONF_FILE) $(CONF_DIR); fi

uninstall:
	rm -Rf $(INSTALL_DIR)

clean:
	rm ./$(BIN_FILE)
