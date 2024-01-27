final: prev: rec {
  golightctrl = prev.callPackage ./linux/golightctrl {};

  default = prev.releaseTools.aggregate {
    name = "all-packages";
    constituents = with final; [ golightctrl ];
  };

}


