package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// 获取本机的MAC地址
func GetLocalMac() string {
	inters, err := net.Interfaces()
	if err != nil {
		fmt.Printf("获得本机网卡出错：%v", err)
		return ""
	}

	var localMac string
	for _, inter := range inters {
		if inter.HardwareAddr != nil {
			return inter.HardwareAddr.String()
		}
	}
	return localMac
}

func QueryMacFirm(macAddr string) string {

	cmp := []byte(macAddr)

	file, err := os.Open("./oui.txt")
	if err != nil {
		fmt.Printf("没有找到[oui.txt]：%v\n", err)
		return ""
	}

	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("查询结束，未找到该设备的生产厂商\n")
				break
			} else {
				fmt.Printf("文件读取出错%v\n", err)
			}

		}

		if len(line) > 6 {
			if bytes.Equal(line[:6], cmp){
				split := strings.Split(string(line), "\t\t")
				return split[1]
			}
		}

	}

	return ""
}

