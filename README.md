[![travis-ci status](https://api.travis-ci.org/travis-ci/travis-web.svg?branch=master "tarvis-ci build status")](https://travis-ci.org/hraban/opus)

## opus, by xiph.org

This package provides Go bindings for the C library libopus and libopusfile (by
xiph.org).

## Details

This wrapper provides a Go translation layer for three elements from the
xiph.org opus libs:

* encoders
* decoders
* stream handlers

### Encoding

To encode raw audio to the Opus format, create an encoder first:

```go
const sample_rate = 48000
const channels = 1 // mono; 2 for stereo

enc, err := opus.NewEncoder(sample_rate, channels, opus.APPLICATION_VOIP)
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
const buffer_size = 1000 // choose any buffer size you like. 1k is plenty.

// Check the frame size. You don't need to do this if you trust your input.
frame_size := len(pcm) // must be interleaved if stereo
frame_size_ms := float32(frame_size) / channels * 1000 / sample_rate
switch frame_size_ms {
case 2.5, 5, 10, 20, 40, 60:
    // Good.
default:
    return fmt.Errorf("Illegal frame size: %d bytes (%f ms)", frame_size, frame_size_ms)
}

data := make([]byte, buffer_size)
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
dec, err := opus.NewDecoder(sample_rate, channels)
if err != nil {
    ...
}
```

Now pass it the opus bytes, and a buffer to store the PCM sound in:

```go
var frame_size_ms float32 = ...  // if you don't know, go with 60 ms.
frame_size := channels * frame_size_ms * sample_rate / 1000
pcm := make([]byte, int(frame_size))
n, err := dec.Decode(data, pcm)
if err != nil {
    ...
}

// To get all samples (interleaved if multiple channels):
pcm = pcm[:n*channels] // only necessary if you didn't know the right frame size

// or access directly:
for i := 0; i < n; i++ {
    ch1 := pcm[i*channels+0]
    // if stereo:
    ch2 := pcm[i*channels+1]
}
```

For more examples, see the `_test.go` files.

### API Docs

Go wrapper API reference:
https://godoc.org/github.com/hraban/opus

Full libopus C API reference:
https://www.opus-codec.org/docs/opus_api-1.1.3/

## Build & installation

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

## License

The licensing terms for the Go bindings are found in the LICENSE file. The
authors and copyright holders are listed in the AUTHORS file.
