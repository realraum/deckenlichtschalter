{ lib
, buildGoModule
}:

# mkNode {
#   root = ./.;
#   pnpmLock = ./pnpm-lock.yaml;
#   nodejs = nodejs;
# } rec {
#   buildInputs = [
#     curl.dev
#     curl.out
#   ];

#   nativeBuildInputs = [
#     makeWrapper
#   ];

#   postInstall = ''
    
#   '';
# }



#////////////////////////


buildGoModule rec {
  name = "golightctrl";

  src = ./.;

  vendorHash = "sha256-VByr4pX4D3+perO+oQN8d1kuO0hr3d3uKxfQRHHzqoY=";

  meta = with lib; {
    description = "Control BasicLights and RF433, handle ActionAlias and read hardware buttons";
  };
}
