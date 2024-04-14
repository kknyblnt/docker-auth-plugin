


all: bin

bin:
	go build -o docker-auth-plugin

clean:
	rm -f docker-auth-plugin