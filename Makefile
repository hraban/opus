BUILDDIR := build

PWD=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))
export PKG_CONFIG_PATH :=$(PWD)/$(BUILDDIR)/lib/pkgconfig:$(PKG_CONFIG_PATH)

.PHONY: clean distclean default test build all libopus libopusfile
# Don't delete config.h files after succesful builds
.PRECIOUS: %/config.h %/autogen.sh

all: libopus libopusfile

libopus: libopus.a
libopusfile: libopus libopusfile.a

%/autogen.sh:
	git submodule init
	git submodule update

%/config.h: %/autogen.sh
	(cd "$*"; ./autogen.sh)
	(cd "$*"; ./configure --prefix="$$PWD/../$(BUILDDIR)" --enable-fixed-point)

%.a: %/config.h
	$(MAKE) -C "$*"
	$(MAKE) -C "$*" install
	cp "$(BUILDDIR)/lib/$@" .

clean-%:
	$(MAKE) -C "$*" clean

distclean-%:
	$(MAKE) -C "$*" distclean

clean distclean: %: %-libopus %-libopusfile
	rm -rf $(BUILDDIR) *.a
