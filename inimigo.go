// inimigo.go - Sistema de comportamento e ações dos inimigos
package main

// ============================================================================
// IMPORTS E DEPENDÊNCIAS
// ============================================================================

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// ============================================================================
// MÓDULO DE DETECÇÃO E DISTÂNCIA
// ============================================================================

// Calcula distância Manhattan entre inimigo e personagem
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

// ============================================================================
// MÓDULO DE MOVIMENTAÇÃO E COMPORTAMENTO
// ============================================================================

// Move inimigo aleatoriamente em uma das 4 direções
func inimigoMover(jogo *Jogo, idx int) {
	// Verifica se o índice ainda é válido
	if idx >= len(jogo.Entidades) {
		return
	}
	
	dx, dy := 0, 0
	
	// Escolhe direção aleatória (0=cima, 1=esquerda, 2=baixo, 3=direita)
	d := rand.IntN(4)
	switch d {
	case 0: dy = -1 // Cima
	case 1: dx = -1 // Esquerda
	case 2: dy = 1  // Baixo
	case 3: dx = 1  // Direita
	}
	
	// Calcula nova posição e move se válida
	nx, ny := jogo.Entidades[idx].X+dx, jogo.Entidades[idx].Y+dy
	if jogoPodeMoverParaInimigo(jogo, nx, ny, idx) {
		jogoMoverElemento(jogo, jogo.Entidades[idx].X, jogo.Entidades[idx].Y, dx, dy, &jogo.Entidades[idx])
	}
}

// Move inimigo em direção ao personagem (perseguição)
func inimigoPerseguir(jogo *Jogo, idx, px, py int) {
	// Verifica se o índice ainda é válido
	if idx >= len(jogo.Entidades) {
		return
	}
	
	dx, dy := 0, 0
	
	// Calcula direção horizontal
	if px > jogo.Entidades[idx].X {
		dx = 1
	} else if px < jogo.Entidades[idx].X {
		dx = -1
	}
	
	// Calcula direção vertical
	if py > jogo.Entidades[idx].Y {
		dy = 1
	} else if py < jogo.Entidades[idx].Y {
		dy = -1
	}
	
	// Move em direção ao personagem se possível
	nx, ny := jogo.Entidades[idx].X+dx, jogo.Entidades[idx].Y+dy
	if jogoPodeMoverParaInimigo(jogo, nx, ny, idx) {
		jogoMoverElemento(jogo, jogo.Entidades[idx].X, jogo.Entidades[idx].Y, dx, dy, &jogo.Entidades[idx])
	}
}

// ============================================================================
// MÓDULO DE SISTEMA DE DANO
// ============================================================================

// Aplica dano ao jogador quando inimigo toca nele
func inimigoAplicarDano(jogo *Jogo, inimigoIdx int) {
	// Verifica se o índice ainda é válido
	if inimigoIdx >= len(jogo.Entidades) {
		return
	}
	
	// Verifica se inimigo está na mesma posição do personagem
	if jogo.Entidades[inimigoIdx].X == jogo.Entidades[0].X && 
	   jogo.Entidades[inimigoIdx].Y == jogo.Entidades[0].Y {
		
		// Protege modificação da vida
		mutexChan <- struct{}{}
		
		// Cooldown de 2 segundos entre danos
		agora := time.Now()
		tempoDecorrido := agora.Sub(jogo.UltimoDano)
		
		// Aplica dano se cooldown passou e jogador tem vida
		if tempoDecorrido >= 2*time.Second && jogo.Vida > 0 {
			jogo.Vida--
			jogo.UltimoDano = agora
		}
		
		<-mutexChan
	}
}

// ============================================================================
// MÓDULO DE PROCESSAMENTO DE AÇÕES
// ============================================================================

// Processa comportamento do inimigo baseado na posição do personagem
func inimigoExecutarAcao(jogo *Jogo, idx int, posPersonagem <-chan [2]int) bool {
	// Verifica se o índice ainda é válido (inimigo pode ter sido removido por explosão)
	if idx >= len(jogo.Entidades) {
		return false // Inimigo foi removido, encerra a goroutine
	}
	
	logIdx := idx - 1 // Ajusta índice para o array de logs (inimigos começam no índice 1)
	
	select {
	case pos := <-posPersonagem:
		// Verifica novamente antes de acessar (pode ter mudado durante o select)
		if idx >= len(jogo.Entidades) {
			return false
		}
		
		// Calcula distância até o personagem
		dist := inimigoDetectaPersonagem(jogo.Entidades[idx].X, jogo.Entidades[idx].Y, pos[0], pos[1])

		if dist <= 10 {
			// Persegue se personagem estiver próximo
			inimigoPerseguir(jogo, idx, pos[0], pos[1])
			inimigoAplicarDano(jogo, idx)
			
			// Atualiza log de comportamento
			if logIdx >= 0 && logIdx < len(jogo.LogsInimigos) {
				jogo.LogsInimigos[logIdx] = fmt.Sprintf("Perseguindo (dist: %d)", dist)
			}
			return true
		} else {
			// Verifica se ainda é válido antes de mover
			if idx >= len(jogo.Entidades) {
				return false
			}
			
			// Movimento aleatório se personagem estiver longe
			inimigoMover(jogo, idx)
			inimigoAplicarDano(jogo, idx)
			
			// Atualiza log de comportamento
			if logIdx >= 0 && logIdx < len(jogo.LogsInimigos) {
				jogo.LogsInimigos[logIdx] = "Random"
			}
		}
	default:
		// Verifica se ainda é válido antes de mover
		if idx >= len(jogo.Entidades) {
			return false
		}
		
		// Movimento aleatório quando não há informação do personagem
		inimigoMover(jogo, idx)
		inimigoAplicarDano(jogo, idx)
		
	}
	
	return true
}
