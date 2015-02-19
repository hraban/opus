BUILDDIR := libopusbuild

.PHONY: libopus clean default test build

default: libopus

build:
	go build -o opus

test:
	go test

libopus/config.h: libopus/autogen.sh
	(cd libopus; ./autogen.sh)
	(cd libopus; ./configure --prefix="$$PWD/../$(BUILDDIR)" --enable-fixed-point)

libopus/autogen.sh:
	git submodule init
	git submodule update

libopus: libopus/config.h
	$(MAKE) -C libopus
	$(MAKE) -C libopus install
	cp $(BUILDDIR)/lib/libopus.a .

%clean:
	$(MAKE) -C libopus $@
	rm -rf $(BUILDDIR) libopus/configure.h 
