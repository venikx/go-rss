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

        #packages.default = pkgs.buildGoApplication {
        #  subPackages = [ "." ];
        #};
        packages.default = packages.go-rss;

        devShells.database = with pkgs;
          pkgs.mkShell {
            nativeBuildInputs = [ postgresql_16 ];

            PGPORT = "21358";
            PGHOST = "localhost";
            PGDATABASE = "kevin_test";
            PGUSER = "test_kevin";
            PGPASSWORD = "test_password";

            shellHook = ''
              pg_pid=""
                    set -euo pipefail
                    # TODO: explain what's happening here
                    LOCAL_PGHOST=$PGHOST
                    LOCAL_PGPORT=$PGPORT
                    LOCAL_PGDATABASE=$PGDATABASE
                    LOCAL_PGUSER=$PGUSER
                    LOCAL_PGPASSWORD=$PGPASSWORD
                    unset PGUSER PGPASSWORD
                    initdb -D $PWD/.pgdata
                    echo "unix_socket_directories = '$(mktemp -d)'" >> $PWD/.pgdata/postgresql.conf
                    # TODO: port
                    pg_ctl -D "$PWD/.pgdata" -w start || (echo pg_ctl failed; exit 1)
                    until psql postgres -c "SELECT 1" > /dev/null 2>&1 ; do
                        echo waiting for pg
                        sleep 0.5
                    done
                    psql postgres -w -c "CREATE DATABASE $LOCAL_PGDATABASE"
                    psql postgres -w -c "CREATE ROLE $LOCAL_PGUSER WITH LOGIN PASSWORD '$LOCAL_PGPASSWORD'"
                    psql postgres -w -c "GRANT ALL PRIVILEGES ON DATABASE $LOCAL_PGDATABASE TO $LOCAL_PGUSER"
            '';

          };

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
