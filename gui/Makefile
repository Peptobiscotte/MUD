IMAGE = localhost/fyne-dev

.PHONY: build run dev clean

build:
	docker build -t $(IMAGE) .

run:
	xhost +SI:localuser:root
	docker run --rm \
		--network host \
		-e DISPLAY=$$DISPLAY \
		-e XAUTHORITY=$$XAUTHORITY \
		-e LIBGL_ALWAYS_SOFTWARE=1 \
		-v $$XAUTHORITY:$$XAUTHORITY:ro \
		-v /tmp/.X11-unix:/tmp/.X11-unix:rw \
		-v $(CURDIR)/..:/app \
		-v fyne-go-mod:/go/pkg/mod \
		-v fyne-go-cache:/root/.cache/go-build \
		-w /app/gui \
		$(IMAGE)
	xhost -SI:localuser:root

dev: build run

clean:
	docker rmi $(IMAGE)
