{ pkgs ? import <nixpkgs> {} }:

let
  inherit (pkgs) lib;
in
pkgs.buildGoModule {
  name = "go-opus";
  src = lib.cleanSource ./.;
  buildInputs = with pkgs; [
    libogg
    libopus
    opusfile
  ];
  nativeBuildInputs = with pkgs; [
    pkg-config
    moreutils
  ];
  checkFlags = "-race";
  postCheck = ''
    gofmt -d .
  '';
  vendorHash = null;
}
