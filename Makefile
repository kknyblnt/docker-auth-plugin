


all: bin

bin:
	go build -o docker-auth-plugin

run:
	./docker-auth-plugin

dev-run:
	go run main.go

priv-dev-run:
	make bin
	sudo ./docker-auth-plugin

clean:
	rm -f docker-auth-plugin

