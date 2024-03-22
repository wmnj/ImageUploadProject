package main

// /Users/wmj/Documents/shixi/项目/bookMark/GoProgrammer/net/src/icmp.go
import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ",os.Args[0],"host")
	}
	service := os.Args[1]
	fmt.Println(service)
	conn, err := net.Dial("ip4:icmp",service)
	checkError(err)
	fmt.Println("dial...")
	var msg [512]byte
	// icmp报文 包括ip头部>=20字节 + icmp报头>=8字节 + icmp报文
	msg[0] = 8 // 类型 0-8 占一字节，标识ICMP报文的类型，目前已定义了14种，从类型值来看ICMP报文可以分为两大类。第一类是取值为1~127的差错报文，第2类是取值128以上的信息报文。ping类型字段值为8
	msg[1] = 0 // 代码 0 占一字节，标识对应ICMP报文的代码。它与类型字段一起共同标识了ICMP报文的详细类型。
	msg[2] = 0 // 校验和
	msg[3] = 0 
	msg[4] = 0   // 标识符[0] 占两字节，用于标识本ICMP进程，但仅适用于回显请求和应答ICMP报文，对于目标不可达ICMP报文和超时ICMP报文等，该字段的值为0。
	msg[5] = 13   // 标识符[1]
	msg[6] = 0   // 序列号[0]
	msg[7] = 37  // 序列号[1]
	len := 8
	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	_, err = conn.Write(msg[0:len])
	checkError(err)
	_, err = conn.Read(msg[0:])
	checkError(err)
	fmt.Println("Got Response")
	if msg[5] == 13 {
		fmt.Println("Sequence matches")
	}
	if msg[7] == 37 {
		fmt.Println("Sequence matches")
	}
	os.Exit(0)
}

 func checkSum(msg []byte)uint16 {
	sum := 0
	for n := 1; n < len(msg) - 1;n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
 }

 func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s",err.Error())
		os.Exit(1)
	}
 }

 func readFully(conn net.Conn)([]byte, error) {
	defer conn.Close()
	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil,err
		}
	}
	return result.Bytes(),nil
 }