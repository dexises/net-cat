package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var joinMessageForUser = []string{
	"Welcome to TCP-Chat!\n",
	"         _nnnn_\n",
	"        dGGGGMMb\n",
	"       @p~qp~~qMb\n",
	"       M|@||@) M|\n",
	"       @,----.JM|\n",
	"      JS^\\__/  qKL\n",
	"     dZP        qKRb\n",
	"    dZP          qKKb\n",
	"   fZP            SMMb\n",
	"   HZM            MMMM\n",
	"   FqM            MMMM\n",
	" __| \".        |\\dS\"qML\n",
	" |    `.       | `' \\Zq\n",
	"_)      \\.___.,|     .'\n",
	"\\____   )MMMMMP|   .'\n",
	"     `-'       `--'\n",
	"[ENTER YOUR NAME]: ",
}

var (
	clients           = make(map[string]net.Conn)
	leaving           = make(chan message)
	messages          = make(chan message)
	historyOfMessages = [][]string{}
)

type message struct {
	time     string
	userName string
	text     string
	address  string
}

func main() {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	var mutex sync.Mutex

	go broadcaster(&mutex)

	for {

		// fmt.Println(len(clients))
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	if len(clients) > 9 {
		fmt.Fprintln(conn, "Number of user is 10, comeback later :)")
		conn.Close()
	} else {

		for _, s := range joinMessageForUser {
			fmt.Fprint(conn, s)
		}

		user := userExist(conn)

		for _, s := range historyOfMessages {
			for _, historyOfMessagesString := range s {
				fmt.Fprint(conn, historyOfMessagesString)
			}
		}
		messages <- newMessage(user, " has joined our chat...", conn)

		input := bufio.NewScanner(conn)
		for input.Scan() {
			messages <- newMessage(user, input.Text(), conn)
		}

		// Delete client form map
		delete(clients, user)

		leaving <- newMessage(user, " has left our chat...", conn)

		conn.Close()
	}
}

func userExist(conn net.Conn) string {
	userName := ""
	userExist := true
	for userExist {
		scanner := bufio.NewScanner(conn)
		scanner.Scan()
		exist := false
		for k := range clients {
			if k == scanner.Text() {
				exist = true
				fmt.Fprint(conn, "Username already exist\n"+"[ENTER YOUR NAME]: ")
				break
			}
		}
		if !exist {
			clients[scanner.Text()] = conn
			userName = scanner.Text()
			userExist = false
		}
	}
	return userName
}

func newMessage(user string, msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	msgTime := "[" + time.Now().Format("01-02-2006 15:04:05") + "]"
	if msg != " has joined our chat..." && msg != " has left our chat..." && msg != "" {
		temp := []string{msgTime, "[" + user + "]", msg, "\n"}
		historyOfMessages = append(historyOfMessages, temp)
	}
	return message{
		time:     msgTime,
		userName: user,
		text:     msg,
		address:  addr,
	}
}

func broadcaster(mutex *sync.Mutex) {
	for {
		mutex.Lock()
		select {
		case msg := <-messages:
			for k, conn := range clients {
				if msg.address == conn.RemoteAddr().String() {
					if msg.text == " has joined our chat..." {
						time.Sleep(time.Millisecond * 300)
					}
					fmt.Fprint(conn, msg.time+"["+msg.userName+"]"+":")
					continue
				}
				if msg.text == " has joined our chat..." {
					fmt.Fprintln(conn, "\n"+msg.userName+msg.text)
					fmt.Fprint(conn, msg.time+"["+k+"]"+":")
					continue
				}
				if msg.text != "" {

					fmt.Fprintln(conn, "\n"+msg.time+"["+msg.userName+"]"+":"+msg.text)
					fmt.Fprint(conn, msg.time+"["+k+"]"+":")
				}

			}

		case msg := <-leaving:
			for k, conn := range clients {
				fmt.Fprintln(conn, "\n"+msg.userName+msg.text)
				fmt.Fprint(conn, msg.time+"["+k+"]"+":")
			}

		}
		mutex.Unlock()
	}
}
