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
  name = "go-opus";
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
