package main

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var restart = false
var p1choisit = false

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.titleUpdate() {
			g.gameState++
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			g.gameState++
		}
	case playState:
		if restart {
			restart = false
		}

		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {

			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {

			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()

		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:

		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	if !g.otherReady {
		select {
		case turn := <-g.readChan:
			g.otherReady = true

			g.turn, _ = strconv.Atoi(turn)

		default:
		}

	}
	return g.otherReady && inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {
	if !p1choisit {
		col := g.p1Color % globalNumColorCol
		line := g.p1Color / globalNumColorLine

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			col = (col + 1) % globalNumColorCol

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			col = (col - 1 + globalNumColorCol) % globalNumColorCol

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			line = (line + 1) % globalNumColorLine

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			line = (line - 1 + globalNumColorLine) % globalNumColorLine
		}
		g.p1Color = line*globalNumColorLine + col
	}
	/*
		var p2Selecteur string
		select {
		case p2Selecteur = <-g.colorchan:
			g.p2Selector, _ = strconv.Atoi(p2Selecteur)
		default:

		}
	*/
	var colorp2 string
	select {
	case colorp2 = <-g.readChan:
		g.p2Color, _ = strconv.Atoi(colorp2)
	default:

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.p1Color != g.p2Color && !g.attente_second_joueur {

		go g.sendColor()
		g.attente_second_joueur = true
		p1choisit = true

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // permet d'annulé le choix de couleur
		go g.resetColor()
		g.attente_second_joueur = false
		p1choisit = false
	}

	if g.attente_second_joueur { // permet au joueur ayant choisit en premier d'être lancé au moment où le second joueur choisit sa couleur
		select {
		case colorp2 = <-g.readChan:
			g.p2Color, _ = strconv.Atoi(colorp2)
		default:

		}
	}
	return g.attente_second_joueur && g.p2Color >= 0
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos

			// Envoyer la position au serveur
			go g.sendPosition()
			log.Print("envoie des donnée")
			g.turn = 1
			log.Print("en attente de donnée")
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update() (int, int) {

	select {
	case message_server := <-g.readChan:

		position, _ := strconv.Atoi(message_server)
		updated, yPos := g.updateGrid(p2Token, position)
		for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
			position = (position + 1) % globalNumTilesX
		}
		g.turn = 0
		return position, yPos
	default:
		return -15, -15
	}
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !restart {
		go g.sendrestart()
		restart = true
		log.Print("attente du second joueur")
		g.attente_second_joueur = false
	}
	if restart {
		select {
		case message := <-g.readChan:
			return (message == "true")
		default:
			return false
		}
	}
	return false
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
