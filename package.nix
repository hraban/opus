{
  lib,
  libogg,
  libopus,
  opusfile,
  pkg-config,
  moreutils,
  buildGoModule,
}:

buildGoModule {
  pname = "go-opus";
  version = "2" + lib.optionalString (libopus ? version) "-${libopus.version}";
  src = ./.;
  buildInputs = [
    libogg
    libopus
    opusfile
  ];
  nativeBuildInputs = [
    pkg-config
    moreutils
  ];
  checkFlags = "-race";
  postCheck = ''
    gofmt -d .
  '';
  vendorHash = null;
}
