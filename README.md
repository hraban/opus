[![Test](https://github.com/hraban/opus/workflows/Test/badge.svg)](https://github.com/hraban/opus/actions?query=workflow%3ATest)

## Go wrapper for Opus

This package provides Go bindings for the xiph.org C libraries libopus and
libopusfile.

The C libraries and docs are hosted at https://opus-codec.org/. This package
just handles the wrapping in Go, and is unaffiliated with xiph.org.

Features:

- ✅ encode and decode raw PCM data to raw Opus data
- ✅ useful when you control the recording device, _and_ the playback
- ✅ decode .opus and .ogg files into raw audio data ("PCM")
- ✅ reuse the system libraries for opus decoding (libopus)
- ✅ works easily on Linux, Mac and Docker; needs libs on Windows
- ❌ does not _create_ .opus or .ogg files (but feel free to send a PR)
- ❌ does not work with .wav files (you need a separate .wav library for that)
- ❌ no self-contained binary (you need the xiph.org libopus lib, e.g. through a package manager)
- ❌ no cross compiling (because it uses CGo)

Good use cases:

- 👍 you are writing a music player app in Go, and you want to play back .opus files
- 👍 you record raw wav in a web app or mobile app, you encode it as Opus on the client, you send the opus to a remote webserver written in Go, and you want to decode it back to raw audio data on that server

## Details

This wrapper provides a Go translation layer for three elements from the
xiph.org opus libs:

* encoders
* decoders
* files & streams

### Import

```go
import "gopkg.in/hraban/opus.v2"
```

### Encoding

To encode raw audio to the Opus format, create an encoder first:

```go
const sampleRate = 48000
const channels = 1 // mono; 2 for stereo

enc, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
if err != nil {
    ...
}
```

Then pass it some raw PCM data to encode.

Make sure that the raw PCM data you want to encode has a legal Opus frame size.
This means it must be exactly 2.5, 5, 10, 20, 40 or 60 ms long. The number of
bytes this corresponds to depends on the sample rate (see the [libopus
documentation](https://www.opus-codec.org/docs/opus_api-1.1.3/group__opus__encoder.html)).

```go
var pcm []int16 = ... // obtain your raw PCM data somewhere
const bufferSize = 1000 // choose any buffer size you like. 1k is plenty.

// Check the frame size. You don't need to do this if you trust your input.
frameSize := len(pcm) // must be interleaved if stereo
frameSizeMs := float32(frameSize) / channels * 1000 / sampleRate
switch frameSizeMs {
case 2.5, 5, 10, 20, 40, 60:
    // Good.
default:
    return fmt.Errorf("Illegal frame size: %d bytes (%f ms)", frameSize, frameSizeMs)
}

data := make([]byte, bufferSize)
n, err := enc.Encode(pcm, data)
if err != nil {
    ...
}
data = data[:n] // only the first N bytes are opus data. Just like io.Reader.
```

Note that you must choose a target buffer size, and this buffer size will affect
the encoding process:

> Size of the allocated memory for the output payload. This may be used to
> impose an upper limit on the instant bitrate, but should not be used as the
> only bitrate control. Use `OPUS_SET_BITRATE` to control the bitrate.

-- https://opus-codec.org/docs/opus_api-1.1.3/group__opus__encoder.html

### Decoding

To decode opus data to raw PCM format, first create a decoder:

```go
dec, err := opus.NewDecoder(sampleRate, channels)
if err != nil {
    ...
}
```

Now pass it the opus bytes, and a buffer to store the PCM sound in:

```go
var frameSizeMs float32 = ...  // if you don't know, go with 60 ms.
frameSize := channels * frameSizeMs * sampleRate / 1000
pcm := make([]int16, int(frameSize))
n, err := dec.Decode(data, pcm)
if err != nil {
    ...
}

// To get all samples (interleaved if multiple channels):
pcm = pcm[:n*channels] // only necessary if you didn't know the right frame size

// or access sample per sample, directly:
for i := 0; i < n; i++ {
    ch1 := pcm[i*channels+0]
    // For stereo output: copy ch1 into ch2 in mono mode, or deinterleave stereo
    ch2 := pcm[(i*channels)+(channels-1)]
}
```

To handle packet loss from an unreliable network, see the
[DecodePLC](https://godoc.org/gopkg.in/hraban/opus.v2#Decoder.DecodePLC) and
[DecodeFEC](https://godoc.org/gopkg.in/hraban/opus.v2#Decoder.DecodeFEC)
options.

### Streams (and Files)

To decode a .opus file (or .ogg with Opus data), or to decode a "Opus stream"
(which is a Ogg stream with Opus data), use the `Stream` interface. It wraps an
io.Reader providing the raw stream bytes and returns the decoded Opus data.

A crude example for reading from a .opus file:

```go
f, err := os.Open(fname)
if err != nil {
    ...
}
s, err := opus.NewStream(f)
if err != nil {
    ...
}
defer s.Close()
pcmbuf := make([]int16, 16384)
for {
    n, err = s.Read(pcmbuf)
    if err == io.EOF {
        break
    } else if err != nil {
        ...
    }
    pcm := pcmbuf[:n*channels]

    // send pcm to audio device here, or write to a .wav file

}
```

See https://godoc.org/gopkg.in/hraban/opus.v2#Stream for further info.

### "My .ogg/.opus file doesn't play!" or "How do I play Opus in VLC / mplayer / ...?"

Note: this package only does _encoding_ of your audio, to _raw opus data_. You can't just dump those all in one big file and play it back. You need extra info. First of all, you need to know how big each individual block is. Remember: opus data is a stream of encoded separate blocks, not one big stream of bytes. Second, you need meta-data: how many channels? What's the sampling rate? Frame size? Etc.

Look closely at the decoding sample code (not stream), above: we're passing all that meta-data in, hard-coded. If you just put all your encoded bytes in one big file and gave that to a media player, it wouldn't know what to do with it. It wouldn't even know that it's Opus data. It would just look like `/dev/random`.

What you need is a [container format](https://en.wikipedia.org/wiki/Container_format_(computing)).

Compare it to video:

* Encodings: MPEG[1234], VP9, H26[45], AV1
* Container formats: .mkv, .avi, .mov, .ogv

For Opus audio, the most common container format is OGG, aka .ogg or .opus. You'll know OGG from OGG/Vorbis: that's [Vorbis](https://xiph.org/vorbis/) encoded audio in an OGG container. So for Opus, you'd call it OGG/Opus. But technically you could stick opus data in any container format that supports it, including e.g. Matroska (.mka for audio, you probably know it from .mkv for video).

Note: libopus, the C library that this wraps, technically comes with libopusfile, which can help with the creation of OGG/Opus streams from raw audio data. I just never needed it myself, so I haven't added the necessary code for it. If you find yourself adding it: send me a PR and we'll get it merged.

This libopus wrapper _does_ come with code for _decoding_ an OGG/Opus stream. Just not for writing one.

### API Docs

Go wrapper API reference:
https://godoc.org/gopkg.in/hraban/opus.v2

Full libopus C API reference:
https://www.opus-codec.org/docs/opus_api-1.1.3/

For more examples, see the `_test.go` files.

## Build & Installation

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

### Building Without `libopusfile`

This package can be built without `libopusfile` by using the build tag `nolibopusfile`.
This enables the compilation of statically-linked binaries with no external
dependencies on operating systems without a static `libopusfile`, such as
[Alpine Linux](https://pkgs.alpinelinux.org/contents?branch=edge&name=opusfile-dev&arch=x86_64&repo=main).

**Note:** this will disable all file and `Stream` APIs.

To enable this feature, add `-tags nolibopusfile` to your `go build` or `go test` commands:

```sh
# Build
go build -tags nolibopusfile ...

# Test
go test -tags nolibopusfile ./...
```

### Using in Docker

If your Dockerized app has this library as a dependency (directly or
indirectly), it will need to install the aforementioned packages, too.

This means you can't use the standard `golang:*-onbuild` images, because those
will try to build the app from source before allowing you to install extra
dependencies. Instead, try this as a Dockerfile:

```Dockerfile
# Choose any golang image, just make sure it doesn't have -onbuild
FROM golang:1

RUN apt-get update && apt-get -y install libopus-dev libopusfile-dev

# Everything below is copied manually from the official -onbuild image,
# with the ONBUILD keywords removed.

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["go-wrapper", "run"]
COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install
```

For more information, see <https://hub.docker.com/_/golang/>.

### Linking libopus and libopusfile

The opus and opusfile libraries will be linked into your application
dynamically. This means everyone who uses the resulting binary will need those
libraries available on their system. E.g. if you use this wrapper to write a
music app in Go, everyone using that music app will need libopus and libopusfile
on their system. On Debian systems the packages are called `libopus0` and
`libopusfile0`.

The "cleanest" way to do this is to publish your software through a package
manager and specify libopus and libopusfile as dependencies of your program. If
that is not an option, you can compile the dynamic libraries yourself and ship
them with your software as seperate (.dll or .so) files.

On Linux, for example, you would need the libopus.so.0 and libopusfile.so.0
files in the same directory as the binary. Set your ELF binary's rpath to
`$ORIGIN` (this is not a shell variable but elf magic):

```sh
patchelf --set-origin '$ORIGIN' your-app-binary
```

Now you can run the binary and it will automatically pick up shared library
files from its own directory.

Wrap it all in a .zip, and ship.

I know there is a similar trick for Mac (involving prefixing the shared library
names with `./`, which is, arguably, better). And Windows... probably just picks
up .dll files from the same dir by default? I don't know. But there are ways.

## License

The licensing terms for the Go bindings are found in the LICENSE file. The
authors and copyright holders are listed in the AUTHORS file.

The copyright notice uses range notation to indicate all years in between are
subject to copyright, as well. This statement is necessary, apparently. For all
those nefarious actors ready to abuse a copyright notice with incorrect
notation, but thwarted by a mention in the README. Pfew!
