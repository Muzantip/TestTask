package main

import (
	"flag"
	"fmt"
	"mq_client/utils"
	"net"
	"strconv"
	"time"
)

func main() {

	var name string
	flag.StringVar(&name, "n", "", "set client name")
	flag.Parse()

	if len(name) == 0 {
		fmt.Println("Need set client name, -n clientName")
		return
	}

	file, err := utils.CreateFile("logTransactions.txt")
	defer file.Close()

	conn, err := net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//отпраляем имя клиента серверу
	if n, err := conn.Write([]byte(name)); n == 0 || err != nil {
		fmt.Println(err)
		return
	}

	for {
		var source string
		fmt.Print("Введите сумму: ")
		_, err1 := fmt.Scanln(&source)
		err2 := inputValidate(source)
		if err1 != nil || err2 != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}
		// отправляем сообщение серверу
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			text := "The transaction failed in the amount of " + source + ", Date:" + time.Now().String()
			err = utils.WriteFile(text, file)
			if err != nil {
				fmt.Println("failed to write a transaction to a file for the amount of ", source)
			}
			return
		}
		// получем ответ
		fmt.Print("Новый баланс: ")
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Print(string(buff[0:n]))
		fmt.Println()
	}
}

func inputValidate(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	return nil
}
