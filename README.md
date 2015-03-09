## opus, by xiph.org

This package provides Go bindings for the C library libopus, by xiph.org. The
library, that is, not the bindings. The bindings are by me, Hraban.

Installation is a bit of a mess right now but I'm working on simplifying it.
Here's the summary:

```sh
git submodule init
git submodule update
make
cp libopusbuild/lib/libopus.a .
go build
```

Now, here's the annoying part: the libopus.a file from the current directory
must always be present in every directory you use this this package from. So if
you create a new media player in Go and you include libopus, you must put that
.a file in your media player project directory.

It can be changed
(https://groups.google.com/d/msg/golang-nuts/lF5skXi7OD4/SjwQbgju91QJ>) but
that's on my TODO list for now.
