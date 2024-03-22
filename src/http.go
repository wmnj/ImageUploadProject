// package main

// // /Users/wmj/Documents/shixi/项目/bookMark/GoProgrammer/net/src/icmp.go
// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"net"
// 	"os"
// )

// func main() {
// 	if len(os.Args) != 2 {
// 		fmt.Println("Usage: ",os.Args[0],"host")
// 	}
// 	service := os.Args[1]
// 	fmt.Println(service)
// 	conn, err := net.Dial("tcp",service)
// 	checkError(err)
// 	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
// 	checkError(err)
// 	result, err := readFully(conn)
// 	checkError(err)
// 	fmt.Println(string(result))
// 	os.Exit(0)
// }

//  func checkSum(msg []byte) uint16 {
// 	sum := 0
// 	for n := 1; n < len(msg) - 1;n += 2 {
// 		sum += int(msg[n])*256 + int(msg[n+1])
// 	}
// 	sum = (sum >> 16) + (sum & 0xffff)
// 	sum += (sum >> 16)
// 	var answer uint16 = uint16(^sum)
// 	return answer
//  }

//  func checkError(err error) {
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Fatal error : %s",err.Error())
// 		os.Exit(1)
// 	}
//  }

//  func readFully(conn net.Conn)([]byte, error) {
// 	defer conn.Close()
// 	result := bytes.NewBuffer(nil)
// 	var buf [512]byte
// 	for {
// 		n, err := conn.Read(buf[0:])
// 		result.Write(buf[0:n])
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return nil,err
// 		}
// 	}
// 	return result.Bytes(),nil
//  }