{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages = flake-utils.lib.flattenTree {
          cointop = let lib = pkgs.lib; in
            pkgs.buildGo117Module {
              pname = "cointop";
              version = "1.6.9";

              modSha256 = lib.fakeSha256;
              vendorSha256 = null;

              src = ./.;

              meta = {
                description = "A fast and lightweight interactive terminal based UI application for tracking cryptocurrencies 🚀";
                homepage = "https://cointop.sh/";
                license = lib.licenses.mit;
                maintainers = [ "johnrichardrinehart" ]; # flake maintainers, not project maintainers
                platforms = lib.platforms.linux ++ lib.platforms.darwin;
              };
            };
        };

        defaultPackage = packages.cointop;
        defaultApp = packages.cointop;
      }
    );
}
