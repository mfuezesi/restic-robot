local: fast
	./restic-robot

fast:
	go build

build:
	docker build \
		-t mfuezesi/restic-robot:v0.10.0 \
		.
