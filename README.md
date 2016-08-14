## opus, by xiph.org

This package provides Go bindings for the C library libopus, by xiph.org. The
library, that is, not the bindings. The bindings are by me. Hraban.

### Build & installation

This package requires libopus and libopusfile development packages to be
installed on your system. These are available on Debian based systems from
aptitude as `libopus-dev` and `libopusfile-dev`, and on Mac OS X from homebrew.

They are linked into the app using pkg-config.

Debian, Ubuntu, ...:
```sh
sudo apt-get install pkg-config libopus-dev libopusfile-dev
```

Mac:
```sh
brew install pkg-config opus opusfile
```

Note that this will link the opus libraries dynamically. This means everyone who
uses the resulting binary will need to have opus and opusfile installed (or
otherwise provided).

### License

The licensing terms for the Go bindings are found in the LICENSE file.

Hraban

hraban@0brg.net
