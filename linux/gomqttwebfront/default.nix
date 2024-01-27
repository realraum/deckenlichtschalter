{ lib
, buildGoModule
}:

buildGoModule rec {
  name = "gomqttwebfront";

  src = ./.;

  vendorHash = "sha256-VByr4pX4D3+perO+oQN8d1kuO0hr3d3uKxfQRHHzqow=";

  meta = with lib; {
    description = "licht.realraum.at webfrontend and js to mqtt bridge";
  };
}
