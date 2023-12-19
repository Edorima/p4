package main

import (
	"log"
	"net"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {
	if len(os.Args) != 2 {
		log.Fatal("L'adresse du serveur est manquante.")
	}

	address := os.Args[1]

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	g := game{readChan: make(chan string, 1)}
	g.p2Color = -50
	g.attente_second_joueur = false

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	go g.handleRead(conn)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
