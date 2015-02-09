build/container: stage/jadvisor Dockerfile
	docker build --no-cache -t jadvisor .
	touch build/container

build/jadvisor: *.go */*.go
	GOOS=linux GOARCH=amd64 godep go build -o build/jadvisor

stage/jadvisor: build/jadvisor
	mkdir -p stage
	cp build/jadvisor stage/jadvisor

release:
	docker tag jadvisor fabric8io/jadvisor
	docker push fabric8io/jadvisor

.PHONY: clean
clean:
	rm -rf build
