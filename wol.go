package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {

	inputMac := flag.String("mac", "112233445566", "目标MAC地址，例如：\n\t1). 11-22-33-44-55-66\n\t2). 11:22:33:44:55:66\n\t3). 112233445566\n\t")
	flag.Parse()

	// 去除-和:分隔符
	stripMac := strings.ReplaceAll(strings.ReplaceAll(*inputMac, ":", ""), "-", "")

	// 发送者：即本机
	sender := net.UDPAddr{}

	// 接受者：广播地址-即网段内的所有主机
	target := net.UDPAddr{
		IP: net.IPv4bcast,
	}

	conn, err := net.DialUDP("udp", &sender, &target)
	if err != nil {
		fmt.Printf("创建连接对象出现错误：%v\n", err)
		return
	}

	var bcastMac = []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	var buffer bytes.Buffer
	buffer.Write(bcastMac)

	macHex, err := hex.DecodeString(stripMac)
	for i := 0; i < 16; i++ {
		buffer.Write(macHex)
	}

	magicPacket := buffer.Bytes()
	lens, err := conn.Write(magicPacket)
	if err != nil {
		fmt.Printf("发送网络数据报文出错：%v\n", err)
		return
	}
	_ = conn.Close()

	if lens == 102 {
		fmt.Printf("发送魔包成功\n")
	} else {
		fmt.Printf("数据已经发送，但是MAC地址输入不合法，请确认MAC地址：[%v]\n", *inputMac)
	}
}
