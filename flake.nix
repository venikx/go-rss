{
  description = "RSS Reader in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
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
          overlays = [ gomod2nix.overlays.default ];
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
            inputsFrom = [ self.packages.${system}.go-rss ];

            packages = [
              air

              gomod2nix.packages.${system}.default
              go-tools
              gopls

              postgresql_17
              goose
            ];

          };
      }));
}
