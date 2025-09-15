// inimigo.go - Funções para movimentação e ações dos inimigos
package main

import (
	"fmt"
	"math/rand/v2"
)

// Atualiza a posição do inimigo aleatoriamente
func inimigoMover(jogo *Jogo, idx int) {
	dx, dy := 0, 0
	d := rand.IntN(4)
	switch d {
	case 0:
		dy = -1
	case 1:
		dx = -1
	case 2:
		dy = 1
	case 3:
		dx = 1
	}
	nx, ny := jogo.Entidades[idx].X+dx, jogo.Entidades[idx].Y+dy
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.Entidades[idx].X, jogo.Entidades[idx].Y, dx, dy, &jogo.Entidades[idx])
	}
}

func inimigoPerseguir(jogo *Jogo, idx, px, py int) {
	dx, dy := 0, 0
	if px > jogo.Entidades[idx].X {
		dx = 1
	} else if px < jogo.Entidades[idx].X {
		dx = -1
	}
	if py > jogo.Entidades[idx].Y {
		dy = 1
	} else if py < jogo.Entidades[idx].Y {
		dy = -1
	}
	nx, ny := jogo.Entidades[idx].X+dx, jogo.Entidades[idx].Y+dy
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.Entidades[idx].X, jogo.Entidades[idx].Y, dx, dy, &jogo.Entidades[idx])
	}
}

// inimigo detecta personagem, ou scanear em volta do personagem e se tiver um inimigo próximo, mandar ele perseguir?
// o personagem detecta inimigos e manda mensagem para eles quando está próximo
// inimigo detectar garante que o inimigo pode correr independentemente atrás do player
// inimigo manda pro player que achou está perto e eles trocam posições até o player estar longe demais
func inimigoDetectaPersonagem(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}

// Processa o evento do teclado e executa a ação correspondente
func inimigoExecutarAcao(jogo *Jogo, idx int, posPersonagem <-chan [2]int) bool {
	select {
	case pos := <-posPersonagem:
		dist := inimigoDetectaPersonagem(jogo.Entidades[idx].X, jogo.Entidades[idx].Y, pos[0], pos[1])

		if dist <= 10 {
			inimigoPerseguir(jogo, idx, pos[0], pos[1])
			jogo.StatusMsg = fmt.Sprintf("Distância do personagem: %d", dist)
			return true
		} else {
			inimigoMover(jogo, idx)
		}
	default:
		inimigoMover(jogo, idx)
		jogo.StatusMsg = "Inimigo se movendo aleatoriamente"
	}
	return true
}
