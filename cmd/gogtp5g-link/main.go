package main

import (
	"fmt"
	"os"
	"path"

	"github.com/free5gc/go-gtp5gnl/linkcmd"
)

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "usage: %v <add|del> <ifname> <ipAddr> [ethDev] [--ran]\n", prog)
}

func main() {
	prog := path.Base(os.Args[0])
	if len(os.Args) < 4 {
		usage(prog)
		os.Exit(1)
	}
	cmd := os.Args[1]
	ifname := os.Args[2]
	ipAddr := os.Args[3]
	var ethDev string
	var role int

	if len(os.Args) == 5 {
		if os.Args[4] == "--ran" {
			role = 1
		} else {
			ethDev = os.Args[4]
		}
	} else if len(os.Args) > 5 {
		ethDev = os.Args[4]
		if os.Args[5] == "--ran" {
			role = 1
		}
	}

	switch cmd {
	case "add":
		stopChan := make(chan bool)
		err := linkcmd.CmdAdd(ifname, role, ipAddr, ethDev, stopChan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", prog, err)
			os.Exit(1)
		}
	case "del":
		err := linkcmd.CmdDel(ifname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", prog, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "%v: unknown command %q\n", prog, cmd)
		os.Exit(1)
	}
}
