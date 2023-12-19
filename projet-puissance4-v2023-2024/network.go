package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	writer *bufio.Writer
	reader *bufio.Reader
)

func (g *game) sendPosition() {
	// envoie la position du token
	message := strconv.Itoa(g.tokenPosition) + "\n"
	writer.WriteString(message)
	writer.Flush()
}

func (g *game) sendColor() {
	// envoie la couleur
	message := strconv.Itoa(g.p1Color) + "\n"
	writer.WriteString(message)
	writer.Flush()
}

/*  pas fonctionnel
func (g *game) sendSelector() {
	// envoie la couleur
	message := "?" + strconv.Itoa(g.p1Color) + "\n"
	writer.WriteString(message)
	writer.Flush()

}*/

func (g *game) resetColor() {
	// reset la color
	message := "-50\n"
	writer.WriteString(message)
	writer.Flush()
}

func (g *game) sendrestart() {

	message := "true\n"
	writer.WriteString(message)
	writer.Flush()
}

func (g *game) handleRead(conn net.Conn) {
	reader = bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		msg = strings.Replace(msg, "\n", "", -1)

		g.readChan <- msg

	}
}
