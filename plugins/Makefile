SUBDIRS := $(wildcard src/*)
PACKS := $(SUBDIRS:src/%=%)

all: $(PACKS)
$(PACKS):
	go build -o ./build/$@.plugin --buildmode=plugin ./src/$@

.PHONY: all clean $(PACKS)

clean:
	rm -f build/*
