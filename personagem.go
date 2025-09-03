// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"time"
	"fmt"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1 // Move para cima
	case 'a': dx = -1 // Move para a esquerda
	case 's': dy = 1  // Move para baixo
	case 'd': dx = 1  // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Função para atirar
func atirar(jogo *Jogo, x, y int){
	dx, dy := 0, 0
	switch jogo.Direcao {
	case 'w': dy = -1 // atira para cima
	case 'a': dx = -1 // atira para a esquerda
	case 's': dy = 1  // atira para baixo
	case 'd': dx = 1  // atira para a direita
	}

	for i := 1; i <= 30; i++ {
		tiroX, tiroY := x+(dx*i), y+(dy*i)
		
		// Verifica se o tiro saiu dos limites do mapa
		if tiroY < 0 || tiroY >= len(jogo.Mapa) || tiroX < 0 || tiroX >= len(jogo.Mapa[tiroY]) {
			break
		}
		
		// Verifica se atingiu uma parede
		if jogo.Mapa[tiroY][tiroX].tangivel {
			break
		}

		// Desenha o tiro
		jogo.Mapa[tiroY][tiroX] = Vegetacao
		interfaceDesenharJogo(jogo)
		time.Sleep(100 * time.Millisecond)
		
		// Apaga o tiro (restaura o elemento original)
		jogo.Mapa[tiroY][tiroX] = Vazio
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
