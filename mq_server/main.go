package main

import (
	"context"
	"fmt"
	"mq_server/internal/db"
	"mq_server/internal/db/client"
	"mq_server/internal/db/config"
	"mq_server/internal/db/transaction"
	"net"
	"strconv"
)

func main() {

	/*
		Подсоединяемся к БД и создаем клиент для дальнейшей работы
	*/
	config := config.StorageConfig{
		Username:    "postgres",
		Password:    "password",
		Host:        "localhost",
		Port:        "9999",
		MaxAttempts: 3,
		Database:    "postgres",
	}
	postgresqlClient, err := db.NewSQLClient(context.TODO(), config)
	if err != nil {
		fmt.Println("NewClient error", err)
	}

	/*
		Создаем экземпляры репозитарий для двух таблиц, клиенты и транзакции
	*/
	clientRep := client.NewRepository(postgresqlClient)
	transactionRep := transaction.NewRepository(postgresqlClient)

	/*
		Создаем слушатель, который будет слушать новые подключения от клиентов и обрабатывать их запросы
	*/
	listener, err := net.Listen("tcp", ":4545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn, clientRep, transactionRep) // запускаем горутину для обработки запроса от конкретного клиента
	}
}

/*
Функция обрабатывает подключенного клиента conn
Пытается создать клиента в БД, дальше обрабатывает запросы от этого клиента и увеличивает или уменьшает баланс, так же пишет каждую транзакцию в таблицу транзакций
conn – экземпляр подключенного клиента
clientRep – экземпляр репозитарий для таблицы клиентов
transactionRep – экземпляр репозитарий для таблицы транзакций
*/

func handleConnection(conn net.Conn, clientRep client.StorageRepository, transactionRep transaction.StorageRepository) {
	defer conn.Close()

	var message string
	// считываем имя клиента
	input := make([]byte, (1024 * 4))
	n, err := conn.Read(input)
	if n == 0 || err != nil {
		fmt.Println("Name Read error:", err)
		return
	}
	name := string(input[0:n])
	fmt.Println("Connect: ", conn.LocalAddr().String(), "Name: ", name)
	cl := client.Client{
		Name:    name,
		IP:      conn.LocalAddr().String(),
		Balance: 0,
	}
	err = clientRep.Create(context.TODO(), &cl)
	if err != nil {
		fmt.Println("Client not created", err)
	}

	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		source := string(input[0:n])

		/*
			На всякий случаи проверяем, есть ли клиент в БД, если клиент есть, начинаем принимать от него транзакции
		*/
		findCl, err := clientRep.FindOne(context.TODO(), name)
		if err != nil {
			message = "Client not find"
			fmt.Println(message)
		} else {
			sum, err := strconv.Atoi(source)
			if err != nil {
				message = "Err string to int convert"
				fmt.Println(message)
			}
			findCl.Balance += sum
			if findCl.Balance >= 0 {
				tr := transaction.Transaction{
					Info:   "Create transaction for client:" + name + " sum: " + source,
					Sum:    sum,
					Client: findCl,
				}
				err = transactionRep.Create(context.TODO(), &tr)
				if err != nil {
					message = "transaction not created"
					fmt.Println(message)
				} else {
					err := clientRep.Update(context.TODO(), findCl)
					if err != nil {
						message = "Client balance not update"
						fmt.Println(message)
					} else {
						message = "Your new balance: " + strconv.Itoa(findCl.Balance)
					}
				}
			} else {
				message = "You don't have enough funds"
			}
		}

		// отправляем данные клиенту
		conn.Write([]byte(message))
	}
}
