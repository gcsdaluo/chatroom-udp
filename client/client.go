package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	MaxBytes = 65535
	ADDRESS  = "127.0.0.1"
	PORT     = 1600
)

func main() {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ADDRESS, PORT))
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	var choice int
	for {
		fmt.Println("请进行选择:")
		fmt.Println("1. 注册")
		fmt.Println("2. 登录")

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Error reading choice:", err)
			continue
		}

		if choice != 1 && choice != 2 {
			fmt.Println("Unknown command")
			continue
		}

		signal, name, address := personMessage(conn, choice, reader)
		if signal == "OK" {
			fmt.Println("\t\t\t\t欢迎进入聊天室\t\t")
			chatMessage(conn, name, address, reader)
			break
		} else if signal == "Error_UserExist" {
			fmt.Println("User already exists!")
		} else if signal == "Error_PasswordError" {
			fmt.Println("Password error!")
		} else if signal == "Error_UserNotExist" {
			fmt.Println("User does not exist!")
		}
	}
}

func personMessage(conn net.Conn, choice int, reader *bufio.Reader) (string, string, string) {
	fmt.Print("请输入名字: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("请输入密码: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	text := fmt.Sprintf("%d  %s  %s", choice, name, password)
	conn.Write([]byte(text))

	buffer := make([]byte, MaxBytes)
	n, _ := conn.Read(buffer)
	return string(buffer[:n]), name, fmt.Sprintf("%s:%d", ADDRESS, PORT)
}

func chatMessage(conn net.Conn, name, address string, reader *bufio.Reader) {
	fmt.Println("请输入聊天信息:")
	fmt.Println("(input \033[1;44mExit\033[0m to quit the room,")
	fmt.Println("input \033[1;44ms/name/message\033[0m for private chat)")

	go receiveMessages(conn)

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "Exit" {
			text := fmt.Sprintf("5  %s  %s", name, address)
			conn.Write([]byte(text))
			fmt.Println("你已经退出聊天室")
			os.Exit(0)
		}

		if strings.HasPrefix(message, "s/") {
			parts := strings.SplitN(message[2:], "/", 2)
			if len(parts) == 2 {
				destination := strings.TrimSpace(parts[0])
				message := strings.TrimSpace(parts[1])
				text := fmt.Sprintf("4  %s  %s  %s  %s", name, address, message, destination)
				conn.Write([]byte(text))
				fmt.Println("OK!")
			} else {
				fmt.Println("Invalid private message format")
			}
		} else {
			text := fmt.Sprintf("3  %s  %s  %s", name, address, message)
			conn.Write([]byte(text))
		}
	}
}

func receiveMessages(conn net.Conn) {
	buffer := make([]byte, MaxBytes)
	for {
		n, _ := conn.Read(buffer)
		message := string(buffer[:n])

		if message == "exit" {
			fmt.Println("You have been kicked out of the chat room")
			os.Exit(0)
		}

		fmt.Print("\t\t\t\t\t\t" + message + "\n请输入聊天信息:\n")
	}
}
