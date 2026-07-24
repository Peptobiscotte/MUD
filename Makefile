IMAGE = localhost/fyne-dev

.PHONY: install run-server run-client run-client-gui lint clean

install:
	go mod tidy
	go mod download
	go mod tidy
	docker build -t $(IMAGE) ./gui
	cd ..

run-server:
	go run ./server

run-client:
	go run ./client

run-client-gui:
	$(MAKE) -C gui run

lint:
	go vet ./...

clean:
	go clean
