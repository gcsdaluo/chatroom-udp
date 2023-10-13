package main

import (
	"fmt"
	"net"
	"strings"
)

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
		buffer := make([]byte, MaxBytes)
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
