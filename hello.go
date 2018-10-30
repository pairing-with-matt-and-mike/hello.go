package main

import "bufio"
import "fmt"
import "net"

type Chatroom struct {
	cs []net.Conn
}

func (r Chatroom) send(data []byte, sender net.Conn) {
	for _, c := range r.cs {
		if sender == c { continue }
		if _, err := c.Write(data); err != nil {
			fmt.Printf("Whoops! %v\n", err)
		}		
	}
}

func handleConnection(r *Chatroom, c net.Conn) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		if data, err := bufio.NewReader(c).ReadString('\n'); err != nil {
			fmt.Println("Socket Closed :(")
			remove(r, c)
			break
		} else {
			r.send([]byte(data), c)
		}
	}
}

func remove(r *Chatroom, c net.Conn) {
	for i, element := range r.cs {
		if element == c {
			r.cs[i] = r.cs[len(r.cs) - 1]
			r.cs = r.cs[:len(r.cs)-1]
			return
		}
	}
}

func main() {
	r := Chatroom{}
	l, _ := net.Listen("tcp", ":8080")
	defer l.Close()
	for {
		c, _ := l.Accept()
		r.cs = append(r.cs, c)
		go handleConnection(&r, c)
	}
}
