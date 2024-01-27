{ config, pkgs, lib, ... }:

with lib;

let
  cfg = config.services.golightctrl;
in
{
  options = {
    services.golightctrl.enable = mkEnableOption "golightctrl";
  };

  config = mkIf (cfg.enable) {
    systemd.services.golightctrl = {
      # startAt = "*-*-* 0/6:00:00";
      path = with pkgs; [ golightctrl nodejs ];
      script = ''
        export HOME=$STATE_DIRECTORY
        golightctrl
      '';
      serviceConfig = {
        StateDirectory = "golightctrl";
        User = "golightctrl";
        DynamicUser = true;
      };
    };
  };
}
