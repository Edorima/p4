package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func connection() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	log.Println("Je suis connecté")
	writer := bufio.NewWriter(conn)
	message := "Hello, serveur!"
	_, err = writer.WriteString(message)
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	writer.Flush()
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message reçu du client: \n", message)
	} else {
		log.Println("Erreur de lecture du message:", scanner.Err())
	}
}
