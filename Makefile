all:
	go mod download && go build cmd/packetrusher.go
clean:
	rm ./packetrusher
