build:
	 go build -o build/cli -i main.go

run:
	./build/cli

clean:
	echo "{}" > ~/.uv/config.json
