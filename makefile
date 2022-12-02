
SIGNAL_IP = 54.91.215.145
SIGNAL_PORT = 9999
PEER1_PORT = 3333
PEER2_PORT = 4444
PEER3_PORT = 5555

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/punch
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin/arm64/punch
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/punch
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/punch
	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/punch
	GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/punch
	GOOS=android GOARCH=arm64 go build -o bin/android/arm64/punch
	zip -r builds.zip bin

deploy-signal: 
	scp -i "rendezvous.pem" bin/linux/amd64/punch ubuntu@$(SIGNAL_IP):~/

ssh-signal:
	ssh -i "rendezvous.pem" ubuntu@$(SIGNAL_IP)

p1: 
	./bin/darwin/amd64/punch client $(SIGNAL_IP):$(SIGNAL_PORT) :$(PEER1_PORT)

p2: 
	./bin/darwin/amd64/punch client $(SIGNAL_IP):$(SIGNAL_PORT) :$(PEER2_PORT)

p3: 
	./bin/darwin/amd64/punch client $(SIGNAL_IP):$(SIGNAL_PORT) :$(PEER3_PORT)

android-screen:
	scrcpy