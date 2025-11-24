{
  description = "Development environment for k6-nats extension";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Build the k6 extension
        k6-nats = pkgs.buildGoModule {
          pname = "k6-nats";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-BfmsCsMOh+6OGMCVJEJ+1akVFRCq3AksjRmZq/TtTm0=";

          buildInputs = with pkgs; [
            nats-server
          ];

          meta = with pkgs.lib; {
            description = "k6 extension for NATS testing";
            homepage = "https://github.com/pondigo/xk6-nats";
            license = licenses.agpl3Only;
            platforms = platforms.all;
          };
        };

        # Development shell with all necessary tools
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            go-tools
            nats-server
            golangci-lint
            delve
            just
            prettier  # Code formatting
          ];

          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$PATH:$GOPATH/bin
            echo "🚀 k6/Pondigo/nats development environment ready"
            echo "📝 Available commands:"
            echo "   go build ./... - Build the extension"
            echo "   go test ./...   - Run tests"
            echo "   golangci-lint run - Run linter"
            echo "   nats-server     - Start NATS server for testing"
          '';
        };
      in
      {
        packages = {
          default = k6-nats;
          k6-nats = k6-nats;
        };

        devShells.default = devShell;

        apps = {
          nats-server = {
            type = "app";
            program = "${pkgs.nats-server}/bin/nats-server";
          };
          k6-nats = {
            type = "app";
            program = "${k6-nats}/bin/k6-nats";
          };
        };
      });
}