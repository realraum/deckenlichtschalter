// (c) Bernhard Tittelbach, 2013, 2015

package main

import "os"

func getFileMTime(filename string) (int64, error) {
	keysfile, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer keysfile.Close()
	stat, err := keysfile.Stat()
	if err != nil {
		return 0, err
	}
	return stat.ModTime().Unix(), nil
}

func EnvironOrDefault(envvarname, defvalue string) string {
	if len(os.Getenv(envvarname)) > 0 {
		return os.Getenv(envvarname)
	} else {
		return defvalue
	}
}

func IfThenElseStr(c bool, strue, sfalse string) string {
	if c {
		return strue
	} else {
		return sfalse
	}
}
