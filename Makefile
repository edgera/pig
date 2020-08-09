EXE=pigm-armv6
DEPLOY_TARGET=

build:
	mkdir -p  build/bin/
	GOARM=6 GOARCH=arm GOOS=linux go build -o build/bin/${EXE} ./main.go

deploy: build
	scp -r  build/bin/${EXE} usbgadget_config  pi@${DEPLOY_TARGET}:/home/pi/


run: deploy
	ssh -t pi@${DEPLOY_TARGET} "/home/pi/${EXE}"


.PHONY: build deploy run config