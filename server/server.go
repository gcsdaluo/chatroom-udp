package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	MAX_BYTES = 65535
	ADDRESS   = "127.0.0.1"
	PORT      = 1600
)

// User struct to store user information
type User struct {
	Password string
	Address  *net.UDPAddr
}

func main() {
	// Users' information stored in a map
	users := make(map[string]User)

	// Create UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ADDRESS, PORT))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// Create UDP connection
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening on", conn.LocalAddr())

	for {
		buffer := make([]byte, MAX_BYTES)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			return
		}

		message := string(buffer[:n])
		textList := strings.Split(message, "  ")

		switch textList[0] {
		case "1":
			// Register
			register(conn, users, textList, addr)
		case "2":
			// Login
			login(conn, users, textList, addr)
		case "3":
			// Public chat
			publicChat(conn, users, textList)
		case "4":
			// Private chat
			privateChat(conn, users, textList)
		case "5":
			// Exit
			exit(conn, users, textList)
		}
	}
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
