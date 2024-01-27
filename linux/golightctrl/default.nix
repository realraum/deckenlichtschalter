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

  #vendorHash = "sha256-ciBIR+a1oaYH+H1PcC8cD8ncfJczk1IiJ8iYNM+R6aA=";

  meta = with lib; {
    description = "Control BasicLights and RF433, handle ActionAlias and read hardware buttons";
  };
}
