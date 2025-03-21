package linkcmd

import (
	"net"
	"sync"
	"syscall"

	"github.com/free5gc/go-gtp5gnl"
	"github.com/khirono/go-nl"
	"github.com/khirono/go-rtnllink"
)

func CmdDel(ifname string) error {
	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		return err
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	err = rtnllink.Remove(c, ifname)
	if err != nil {
		return err
	}

	return nil
}

func CmdAdd(ifname string, role int, ipAddr string, ethDev string, stopChan chan bool) error {
	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		return err
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	laddr, err := net.ResolveUDPAddr("udp4", ipAddr+":2152")
	if err != nil {
		return err
	}
	conn2, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}
	defer conn2.Close()
	f, err := conn2.File()
	if err != nil {
		return err
	}
	defer f.Close()

	linkinfo := &nl.Attr{
		Type: syscall.IFLA_LINKINFO,
		Value: nl.AttrList{
			{
				Type:  rtnllink.IFLA_INFO_KIND,
				Value: nl.AttrString("gtp5g"),
			},
			{
				Type: rtnllink.IFLA_INFO_DATA,
				Value: nl.AttrList{
					{
						Type:  gtp5gnl.IFLA_FD1,
						Value: nl.AttrU32(f.Fd()),
					},
					{
						Type:  gtp5gnl.IFLA_HASHSIZE,
						Value: nl.AttrU32(131072),
					},
					{
						Type:  gtp5gnl.IFLA_ROLE,
						Value: nl.AttrU32(role),
					},
					{
						Type:  gtp5gnl.IFLA_ETHERNET_N6_DEV,
						Value: nl.AttrString(ethDev),
					},
				},
			},
		},
	}
	err = rtnllink.Create(c, ifname, linkinfo)
	if err != nil {
		return err
	}

	err = rtnllink.Up(c, ifname)
	if err != nil {
		return err
	}

	<-stopChan

	return nil
}
