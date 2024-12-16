{
  description = "A command line interface for Pleasant Password Server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Version extraction from source code
        rootGoContent = builtins.readFile ./cmd/root.go;
        versionMatch = builtins.match
          ".*var version = \"([^\"]*)\".*" rootGoContent;
        version =
          if versionMatch == null
          then throw "Version not found in cmd/root.go"
          else builtins.head versionMatch;

      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "pleasant-cli";
          inherit version;
          src = ./.;
          vendorHash = "sha256-mRAlpDUg+2O/ShhxNVGGfK8YUJvWcR/ojossgZ7rCEY=";
        };

        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      }
    );
}
