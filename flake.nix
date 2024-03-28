# JOSSO development environment: tools to work with JOSSO/IAM.tf

{
  description = "JOSSO/IAM.tf GO project flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.11";
    flake-utils.url = "github:numtide/flake-utils";
    josso-pkgs.url = "github:sgonzalezoyuela/mynix-pkgs";
  };
  outputs = { self, nixpkgs, flake-utils, josso-pkgs }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs {
            inherit  system;
            config = {
              allowUnfree = true;
            };
          };
 

        in
        with pkgs;
        {
          devShells.default = mkShell {
            buildInputs = [
              gnumake
              go
              gotools
              gopls
              go-outline
              gocode
              gopkgs
              gocode-gomod
              godef
              golint
              gofumpt
              terraform
            ];
            shellHook = ''
              echo "JOSSO/IAM.tf GO environment"
              echo " - ${pkgs.go.name}"
              echo " - ${terraform.name}"
              echo " - ${gnumake.name}"
            '';
          };
        }
      );
}
