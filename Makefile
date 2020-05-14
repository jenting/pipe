all:
	go run main.go

clean:
	@if [ -a openssl.conf ]; then rm openssl.conf; fi;
