GO_ROOT         := $(shell go env GOROOT)

HOST            := localhost

OUT_DIR         := out
CERT            := $(OUT_DIR)/cert.pem
KEY             := $(OUT_DIR)/key.pem

OUT_DIR         := out
EXECUTABLE_NAME := pwa-server
OUT_FILE        := $(OUT_DIR)/$(EXECUTABLE_NAME)

all: $(OUT_FILE) $(CERT) $(KEY)

clean-certs:
	$(info Cleaning cert files...)
	@rm -rf $(CERT) $(KEY)
	@echo Cleaned cert files!

clean:
	$(info Cleaning...)
	@rm -rf $(OUT_FILE)
	@echo Cleaned!

clean-all: clean-certs clean

$(OUT_DIR):
	mkdir $(OUT_DIR)

$(CERT) $(KEY) &: | $(OUT_DIR)
	$(info Creating cert files...)
	go run "$(GO_ROOT)/src/crypto/tls/generate_cert.go" --host $(HOST)
	mv *.pem $(OUT_DIR)

$(OUT_FILE): | $(OUT_DIR)
	go build -o $(OUT_FILE)

.PHONY: clean clean-certs clean-all all $(OUT_FILE)
