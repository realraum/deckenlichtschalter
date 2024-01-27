{ config, pkgs, lib, ... }:

with lib;

let
  cfg = config.services.golightctrl;
  format = pkgs.formats.keyValue {};
in
{
  options = {
    services.golightctrl = {
      enable = mkEnableOption "golightctrl";
      settings = mkOption {
        type = format.type;
        description = "environment variables for golightctrl";
      };
    };
  };

  config = mkIf (cfg.enable) {

    users.groups.licht = {};
    users.users.licht = {
      isSystemUser = true;
      group = "licht";
    };

    systemd.services.golightctrl = {
      wants = ["network-online.target"];
      preStart = ''
        for gpio in 4 17 18 21 22 23; do
          while ! [ -e /sys/class/gpio/gpio$gpio ]; do
            echo $gpio >| /sys/class/gpio/export;
            sleep 1;
          done;
        done
      '';
      serviceConfig = {
        Type="simple";
        User = "licht";
        DynamicUser = true;
        EnvironmentFile = format.generate "golightctrl.env" cfg.settings;
        ExecStart="${pkgs.golightctrl}/bin/golightctrl";
        SyslogIdentifier = "%i";
        Restart = "always";
        RestartSec="3s";
      };
    };
  };
}
