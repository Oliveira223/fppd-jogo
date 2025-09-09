// inimigo.go - Funções para movimentação e ações dos inimigos
package main

import (
	"fmt"
	"math/rand/v2"
)

// Atualiza a posição do inimigo aleatoriamente
func inimigoMover(jogo *Jogo) {

	dx, dy := 0, 0
	d := rand.IntN(4)
	switch d {
	case 0:
		dy = -1 // Move para cima
	case 1:
		dx = -1 // Move para a esquerda
	case 2:
		dy = 1 // Move para baixo
	case 3:
		dx = 1 // Move para a direita
	}
	nx, ny := jogo.xInim+dx, jogo.yInim+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.xInim, jogo.yInim, dx, dy)
		jogo.StatusMsg = fmt.Sprintf("Atualizando inimigo em (%d, %d) \n direção do inimigo: %d", jogo.xInim, jogo.yInim, d)
		jogo.xInim, jogo.yInim = nx, ny
	}

}

// inimigo detecta personagem, ou scanear em volta do personagem e se tiver um inimigo próximo, mandar ele perseguir?
// o personagem detecta inimigos e manda mensagem para eles quando está próximo
// inimigo detectar garante que o personagem pode correr independentemente atrás do player
// inimigo manda pro player que achou está perto e eles trocam posições até o player estar longe demais
func detectaPersonagem(jogo *Jogo, posXinimigo, posYinimigo int) {
}

// Processa o evento do teclado e executa a ação correspondente
func inimigoExecutarAcao(jogo *Jogo) bool {

	inimigoMover(jogo)

	return true // Continua o jogo
}
