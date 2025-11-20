package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	count    = flag.Int("c", 4, "stop after <count> replies")
	ttl      = flag.Int("t", 64, "define time to live")
	interval = flag.Duration("i", time.Second, "seconds between sending each packet")
	size     = flag.Int("s", 32, "use <size> as number of data bytes to be sent")
	help     = flag.Bool("h", false, "print help and exit")
	forceV4  = flag.Bool("4", false, "force IPv4")
	forceV6  = flag.Bool("6", false, "force IPv6")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  sudo go run main.go [options] <destination>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  sudo go run main.go -c 5 -t 128 google.com\n\n")
		fmt.Fprintf(os.Stderr, "Telegram: \n  - @lochinbekdev\n")
		fmt.Fprintf(os.Stderr, "GitHub: \n  - https://github.com/lochinbekdev/networkapps\n\n")
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Xatolik: Manzil kiritilmadi!")
		flag.Usage()
		os.Exit(1)
	}

	destination := flag.Arg(0)

	// 2. IP versiyasini tanlash
	network := "ip"
	if *forceV4 {
		network = "ip4"
	} else if *forceV6 {
		network = "ip6"
	}

	dst, err := net.ResolveIPAddr(network, destination)
	if err != nil {
		fmt.Printf("Manzil topilmadi: %v\n", err)
		return
	}

	// IPv4 yoki IPv6 ekanligini aniqlash
	isIPv4 := dst.IP.To4() != nil
	if *forceV6 && isIPv4 {
		fmt.Println("Xatolik: -6 tanlangan, lekin manzil IPv4 (ehtimol domen faqat v4 ga ega)")
		return
	}

	var protocol string
	if isIPv4 {
		protocol = "ip4:icmp"
	} else {
		protocol = "ip6:ipv6-icmp"
	}

	fmt.Printf("PING %s (%s) %d data bytes\n", destination, dst.String(), *size)

	// 3. Ulanishni ochish
	c, err := net.ListenPacket(protocol, "")
	if isIPv4 {
		// IPv4 uchun 0.0.0.0 ishlatish xavfsizroq
		c, err = net.ListenPacket(protocol, "0.0.0.0")
	}
	if err != nil {
		fmt.Printf("Ulanishda xatolik (sudo kerakdir?): %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	// Wrapperlar
	var p4 *ipv4.PacketConn
	var p6 *ipv6.PacketConn

	if isIPv4 {
		p4 = ipv4.NewPacketConn(c)
		p4.SetControlMessage(ipv4.FlagTTL, true)
	} else {
		p6 = ipv6.NewPacketConn(c)
		p6.SetControlMessage(ipv6.FlagHopLimit, true)
	}

	// 4. Asosiy Loop
	for i := 0; i < *count; i++ {
		data := make([]byte, *size)
		copy(data, "GO-PING")

		// --- TUZATILGAN QISM ---
		var msgType icmp.Type // MUHIM: Umumiy interfeysdan foydalanamiz
		if isIPv4 {
			msgType = ipv4.ICMPTypeEcho
		} else {
			msgType = ipv6.ICMPTypeEchoRequest
		}
		// -----------------------

		msg := icmp.Message{
			Type: msgType,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i + 1,
				Data: data,
			},
		}

		msgBytes, _ := msg.Marshal(nil)

		// Yuborish
		start := time.Now()
		var err error
		if isIPv4 {
			p4.SetTTL(*ttl)
			_, err = p4.WriteTo(msgBytes, nil, dst)
		} else {
			p6.SetHopLimit(*ttl)
			_, err = p6.WriteTo(msgBytes, nil, dst)
		}

		if err != nil {
			fmt.Printf("Yuborish xatosi: %v\n", err)
			continue
		}

		// Javobni kutish
		reply := make([]byte, 1500)
		c.SetReadDeadline(time.Now().Add(*interval + time.Second))

		var n int
		var receivedTTL int
		var src net.Addr

		if isIPv4 {
			var cm *ipv4.ControlMessage
			n, cm, src, err = p4.ReadFrom(reply)
			if cm != nil {
				receivedTTL = cm.TTL
			}
		} else {
			var cm *ipv6.ControlMessage
			n, cm, src, err = p6.ReadFrom(reply)
			if cm != nil {
				receivedTTL = cm.HopLimit
			}
		}

		if err != nil {
			fmt.Printf("Request timeout for icmp_seq=%d\n", i+1)
		} else {
			duration := time.Since(start)
			
			// Parsing qilib javob turini tekshirish (Reply keldimi?)
			var protoNum int
			if isIPv4 {
				protoNum = 1
			} else {
				protoNum = 58
			}
			
			parsedMsg, parseErr := icmp.ParseMessage(protoNum, reply[:n])
			
			if parseErr == nil && (parsedMsg.Type == ipv4.ICMPTypeEchoReply || parsedMsg.Type == ipv6.ICMPTypeEchoReply) {
				fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%v\n",
					n, src, i+1, receivedTTL, duration)
			} else {
				// Agar echo reply bo'lmasa (masalan Destination Unreachable)
				fmt.Printf("From %s: icmp_seq=%d (boshqa xabar: %v)\n", src, i+1, parsedMsg.Type)
			}
		}

		if i < *count-1 {
			time.Sleep(*interval)
		}
	}
}