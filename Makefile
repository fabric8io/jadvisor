build/jadvisor: *.go */*.go
	GOOS=linux GOARCH=amd64 godep go build -o build/jadvisor

image:
	docker build --no-cache -t jadvisor .

release:
	docker tag jadvisor fabric8/jadvisor
	docker push fabric8/jadvisor

.PHONY: clean
clean:
	rm -rf build
