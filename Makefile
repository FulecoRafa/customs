build: build_plugins
	go build ./src

build_plugins:
	make -C plugins

clean:
	make -C plugins clean
