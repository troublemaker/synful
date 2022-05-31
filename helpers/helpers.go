package helpers

import (
	"errors"
	"fmt"
	"net"
)

func GetInterfaceIpv4(intfName string) (net.IP, error) {

	var ief *net.Interface
	var addrs []net.Addr
	var ipv4Addr net.IP

	ief, err := net.InterfaceByName(intfName)
	if err != nil {
		return nil, err
	}

	addrs, err = ief.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipv4Addr = addr.(*net.IPNet).IP.To4()
		if ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return nil, errors.New(fmt.Sprintf("No IP4 on interface %s\n", intfName))
	}
	return ipv4Addr, nil
}
