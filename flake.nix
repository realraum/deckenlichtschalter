{
  description = "r3 Deckenlichter";

  ## base system
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:

    let
      supportedSystems = [ "x86_64-linux" "armv6l-linux" "aarch64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f system);
    in

    {
      overlays.default = import ./overlay.nix;

      packages = forAllSystems (system: (import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      }));

      nixosModules = {
        golightctrl = import ./linux/golightctrl/module.nix;
        gomqttwebfront = import ./linux/gomqttwebfront/module.nix;
      };

    };
}
