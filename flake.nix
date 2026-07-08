{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-26.05";
    # For libopus 1.5.2
    nixpkgs-2511.url = "github:nixos/nixpkgs/nixos-25.11";
    # libopus 1.4
    nixpkgs-2311.url = "github:nixos/nixpkgs/nixos-23.11";
    # libopus 1.3
    nixpkgs-2305.url = "github:nixos/nixpkgs/nixos-23.05";
    flake-parts.url = "github:hercules-ci/flake-parts";
    systems.url = "github:nix-systems/default";
    treefmt = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      flake-parts,
      systems,
      ...
    }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;
      imports = [ inputs.treefmt.flakeModule ];
      perSystem =
        {
          self',
          inputs',
          lib,
          pkgs,
          ...
        }:
        {
          packages.default = pkgs.callPackage ./package.nix { };

          checks = {
            libopus16 = self'.packages.default;
            # Bit of a blunt hammer to pull in an entire nixpkgs just to control the
            # libopus version, but it’s the easiest way to make sure that specific
            # libopus version actually builds.  Of course the problem then becomes
            # that this older nixpkgs’ go build tooling doesn’t support this package
            # anymore.  It’s a trade-off.
            libopus152 = inputs'.nixpkgs-2511.legacyPackages.callPackage ./package.nix { };
            libopus131 = inputs'.nixpkgs-2305.legacyPackages.callPackage ./package.nix { };
          }
          // lib.optionalAttrs (!pkgs.stdenv.hostPlatform.isDarwin) {
            # Broken hydra build for one of the Darwin dependencies in this
            # branch
            libopus14 = inputs'.nixpkgs-2311.legacyPackages.callPackage ./package.nix { };
          };

          treefmt.programs.nixfmt = {
            enable = true;
            strict = true;
          };
        };
    };
}
