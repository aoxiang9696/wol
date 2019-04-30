package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"strings"
	"wol/utils"
)

func main() {

	localMac := utils.GetLocalMac()

	inputMac := flag.String("mac", localMac, "目标MAC地址，例如：\n\t\t\t11-22-33-44-55-66" +
		"\n\t\t\t11:22:33:44:55:66\n\t\t\t112233445566")
	flag.Parse()

	// 去除-和:分隔符
	stripMac := strings.ReplaceAll(strings.ReplaceAll(*inputMac, ":", ""), "-", "")
	if len(stripMac) != 12 {
		fmt.Printf("MAC地址输入有误：%v\n", *inputMac)
		return
	}

	// 发送者：即本机
	sender := net.UDPAddr{}

	// 接受者：广播地址-即网段内的所有主机
	target := net.UDPAddr{
		IP: net.IPv4bcast,
	}

	conn, err := net.DialUDP("udp", &sender, &target)
	if err != nil {
		fmt.Printf("\033[31m创建连接对象出现错误：%v\033[0m\n", err)
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
		fmt.Printf("\033[31m发送网络数据报文出错：%v\033[0m\n", err)
		return
	}
	_ = conn.Close()

	if lens == 102 {
		macFirm := utils.QueryMacFirm(strings.ToUpper(stripMac[:6]))
		fmt.Printf("\033[32mMagic Packet发送成功，目标MAC地址为：[%v]\t生产厂商：[%v]\033[0m\n", *inputMac, macFirm)
	} else {
		fmt.Printf("\033[33mMagic Packet已经发送，但是输入的MAC地址可能不合法，目标主机可能不会被唤醒，请确认MAC地址：[%v]\033[0m\n", *inputMac)
	}
}
