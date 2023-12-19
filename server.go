package main

import (
	"bufio"
	"log"
	"net"
)

func passeur(conn net.Conn, conn2 net.Conn) {
	writer_pass := bufio.NewWriter(conn2)
	reader_pass := bufio.NewReader(conn)
	for {

		message, err := reader_pass.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			break
		}
		writer_pass.WriteString(message)
		writer_pass.Flush()
	}

}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	log.Print("En attente d'une connextion")

	defer listener.Close()
	conn1, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	writer1 := bufio.NewWriter(conn1)
	//reader1 := bufio.NewReader(conn1)

	conn2, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	writer2 := bufio.NewWriter(conn2)
	//reader2 := bufio.NewReader(conn2)

	message1 := "0\n"
	message2 := "1\n"
	writer1.WriteString(message1)
	writer2.WriteString(message2)
	writer1.Flush()
	writer2.Flush()
	log.Print("Les clients sont connectés")
	go passeur(conn1, conn2)
	go passeur(conn2, conn1)

	select {}
}
