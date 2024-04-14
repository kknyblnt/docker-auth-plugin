


all: bin

bin:
	go build -o docker-auth-plugin

run:
	./docker-auth-plugin

dev-run:
	go run main.go

clean:
	rm -f docker-auth-plugin