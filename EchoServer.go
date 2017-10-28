package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var port = flag.String("p", "8888", "Set Listen Port, default 8888")
var seconds = flag.Int("t", 600, "Set Timeout, default 600, 10minutes")
var useBroadcast = flag.Bool("b", false, "Use broadcast mode")

const timeFormatStr = "2006-01-02T15:04:05.999"

var connMap map[string]net.Conn
var connMutex sync.Mutex

func handleConn(c net.Conn, timeout time.Duration) {
	defer (func() {
		timeStr := time.Now().Format(timeFormatStr)
		fmt.Printf("%s Disconnected=%s\n", timeStr, c.RemoteAddr().String())
		connMutex.Lock()
		delete(connMap, c.RemoteAddr().String())
		connMutex.Unlock()
		c.Close()
	})()

	timeStr := time.Now().Format(timeFormatStr)
	fmt.Printf("%s Connected=%s\n", timeStr, c.RemoteAddr().String())
	connMutex.Lock()
	connMap[c.RemoteAddr().String()] = c
	connMutex.Unlock()

	data := make([]byte, 2048)
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		fmt.Println(err)
		return
	}
	for {
		if readLen, err := c.Read(data); err != nil {
			fmt.Println(err)
			return
		} else if readLen > 0 {
			if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
				fmt.Println(err)
				return
			}
			timeStr = time.Now().Format(timeFormatStr)
			fmt.Printf("%s From=%v, Len=%v, Data=[ ", timeStr, c.RemoteAddr().String(), readLen)
			for _, d := range data[0:readLen] {
				fmt.Printf("%02X ", d)
			}
			fmt.Printf("]\n")
			if *useBroadcast {
				connMutex.Lock()
				for _, conn := range connMap {
					if _, err := conn.Write(data[0:readLen]); err != nil {
						fmt.Println(err)
						continue
					}
				}
				connMutex.Unlock()
			} else {
				if _, err := c.Write(data[0:readLen]); err != nil {
					fmt.Println(err)
					continue
				}
			}

		}
	}
}

func main() {
	flag.Parse()
	timeout := time.Duration((*seconds)) * time.Second
	defer (func() {
		timeStr := time.Now().Format(timeFormatStr)
		fmt.Println(timeStr, "Server is stopped")
	})()
	for {
		l, err := net.Listen("tcp", ":"+*port)
		if err != nil {
			fmt.Println("Listen error:", err)
			return
		}
		timeStr := time.Now().Format(timeFormatStr)
		modeStr := "Single"
		if *useBroadcast {
			modeStr = "Broadcast"
		}

		fmt.Printf("%s Listen on Port=%v, Timeout=%v, EchoMode=%s\n", timeStr, *port, timeout, modeStr)
		connMap = make(map[string]net.Conn)
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println("Accept Error: ", err)
				break
			}
			go handleConn(c, timeout)
		}
	}
}
