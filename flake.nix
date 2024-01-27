{
  description = "r3 Deckenlichter";

  ## base system
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:

    let
      supportedSystems = [ "x86_64-linux" "armv6l-linux" "arm64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f system);
    in

    {
      overlays.default = import ./overlay.nix;

      defaultPackage = forAllSystems (system: (import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      }).golightctrl);

      nixosModules = {
        golightctrl = import ./linux/golightctrl/module.nix;
      };

    };
}
