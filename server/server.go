package main

import (
	"fmt"
	"net"
	"time"
)

const (
	MaxBytes = 65535
	ADDRESS  = "127.0.0.1"
	PORT     = 1600
)

// User struct to store user information
type User struct {
	Password string
	Address  *net.UDPAddr
}

func register(conn *net.UDPConn, users map[string]User, textList []string, addr *net.UDPAddr) {
	name := textList[1]
	password := textList[2]

	if _, exists := users[name]; exists {
		conn.WriteToUDP([]byte("Error_UserExist"), addr)
	} else {
		users[name] = User{Password: password, Address: addr}
		conn.WriteToUDP([]byte("OK"), addr)
		fmt.Println(name, "is registered")
	}
}

func login(conn *net.UDPConn, users map[string]User, textList []string, addr *net.UDPAddr) {
	name := textList[1]
	password := textList[2]

	if user, exists := users[name]; exists {
		if user.Password == password {
			conn.WriteToUDP([]byte("OK"), addr)
			fmt.Println(name, "is logged in")
		} else {
			conn.WriteToUDP([]byte("Error_PasswordError"), addr)
		}
	} else {
		conn.WriteToUDP([]byte("Error_UserNotExist"), addr)
	}
}

func publicChat(conn *net.UDPConn, users map[string]User, textList []string) {
	name := textList[1]
	message := textList[3]
	data := "[" + name + "]: " + message

	for user := range users {
		if user != name {
			conn.WriteToUDP([]byte(data), users[user].Address)
		}
	}
	fmt.Println("[", time.Now(), "]", "[", name, "]:", message)
}

func privateChat(conn *net.UDPConn, users map[string]User, textList []string) {
	name := textList[1]
	message := textList[3]
	destination := textList[4]
	data := "[" + name + "]: " + message

	if user, exists := users[destination]; exists {
		conn.WriteToUDP([]byte(data), user.Address)
		fmt.Println("[", time.Now(), "]", "[", name, "] to [", destination, "]:", message)
	} else {
		fmt.Println("Error: User", destination, "not found")
	}
}

func exit(conn *net.UDPConn, users map[string]User, textList []string) {
	name := textList[1]

	if user, exists := users[name]; exists {
		conn.WriteToUDP([]byte("exit"), user.Address)
		fmt.Println(name, "has left the room")
		delete(users, name)
	}
}
