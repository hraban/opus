## opus, by xiph.org

This package provides Go bindings for the C library libopus, by xiph.org. The
library, that is, not the bindings. The bindings are by me, Hraban.

### Build & installation

Installation is a bit of a mess right now but I'm working on simplifying it.
Here's the summary:

```sh
make
go build
```

Now, here's the annoying part: the libopus.a and libopusfile.a files from the
current directory must always be present in every directory you use this this
package from. So if you create a new media player in Go and you include libopus,
you must put those .a files in your media player project directory.

The advantage is that this puts the burden of managing libopus on the user of
this go opus package (i.e. you), but not on the user of *their* program. In the
media player example: you, the programmer, will have to manage the .a files, but
the users of your media player can use those bundled .a files without having to
apt-get or otherwise install anything related to opus. This is a worthwile
trade-off as long as opus is not widely available for all major platforms, yet.

It can be changed
(<https://groups.google.com/d/msg/golang-nuts/lF5skXi7OD4/SjwQbgju91QJ>) but
that's on my TODO list for now, under "not worth the trouble until a lot of
people actually use this."

### License

The licensing terms for the Go bindings are found in the LICENSE file. The
licensing terms for libopus and libopusfile are probably found in their
respective source directories, which are checked out as git submodules as part
of the build process.

However, because libopus and libopusfile are not (by source nor binaries)
included in this package, this package, in source code form, is (should be?)
unaffected by their licenses. That changes once you run make, which will
download the Opus source code, and go build, which will link it.

Hraban

hraban@0brg.net
