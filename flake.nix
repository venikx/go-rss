{
  description = "RSS Reader in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };
  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        version = builtins.substring 0 8 self.lastModifiedDate;
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            (final: prev: {
              go = prev.go_1_21;
              buildGoModule = prev.buildGo121Module;
            })
            gomod2nix.overlays.default
          ];
        };
      in rec {
        packages.go-rss = pkgs.buildGoApplication {
          inherit version;
          pname = "go-rss";
          src = ./.;
          modules = ./gomod2nix.toml;
        };
        packages.default = packages.go-rss;

        devShells.default = with pkgs;
          pkgs.mkShell {
            packages = [
              air

              gomod2nix.packages.${system}.default
              gotools
              go-tools
              gopls
              go-outline
              gopkgs
              gocode-gomod
              godef
              golint
              templ

              # debugger
              delve
              gdlv

              goose
              sqlc
            ];

            nativeBuildInputs = [ go_1_21 ];
          };
      }));
}
