GO_ROOT            := $(shell go env GOROOT)

HOST               := localhost

OUT_DIR            := out
CERT               := $(OUT_DIR)/cert.pem
KEY                := $(OUT_DIR)/key.pem

OUT_DIR            := out
EXECUTABLE_NAME    := pwa-server
EXECUTABLE_VERSION := latest
IMAGE_TAG          := $(EXECUTABLE_NAME):$(EXECUTABLE_VERSION)
OUT_FILE           := $(OUT_DIR)/$(EXECUTABLE_NAME)

all: $(OUT_FILE) certs image

clean-certs:
	$(info Cleaning cert files...)
	@rm -rf $(CERT) $(KEY)
	@echo Cleaned cert files!

clean-image:
	$(info Cleaning image...)
	-docker rmi $(IMAGE_TAG)
	docker buildx prune -f

clean:
	$(info Cleaning...)
	@rm -rf $(OUT_FILE)
	@echo Cleaned!

clean-all: clean-certs clean clean-image
ca: clean-all

$(OUT_DIR):
	mkdir $(OUT_DIR)

$(CERT) $(KEY) &: | $(OUT_DIR)
	$(info Creating cert files...)
	go run "$(GO_ROOT)/src/crypto/tls/generate_cert.go" --host $(HOST)
	mv *.pem $(OUT_DIR)

$(OUT_FILE): | $(OUT_DIR)
# go build -o $(OUT_FILE) -ldflags="-extldflags=-static"
	go build -o $(OUT_FILE) -tags 'osusergo netgo static' -buildmode=pie -ldflags '-linkmode=external -extldflags "-static-pie"'

certs: $(CERT) $(KEY)
executable: $(OUT_FILE)

image:
	docker buildx build \
	  --progress=plain \
	  --tag $(IMAGE_TAG) \
	  .

.PHONY: clean clean-certs clean-all ca all certs executable image
