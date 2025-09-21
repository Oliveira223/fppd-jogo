// bomba.go - Lógica concorrente da bomba
package main

import (
	"time"
)

// gerenciarBomba é a função que roda como uma goroutine para a bomba. 
func gerenciarBomba(jogo *Jogo, x, y int) {
	// A bomba só precisa esperar pelo canal do temporizador.
	<-time.After(3 * time.Second) // A bomba explode após 3 segundos

	jogo.StatusMsg = "A bomba explodiu!"

	// Lógica da explosão
	raio := 2 // Raio da explosão
	for i := -raio; i <= raio; i++ {
		for j := -raio; j <= raio; j++ {
			nx, ny := x+i, y+j
			// Verifica os limites do mapa
			if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
				// Envia uma solicitação para desenhar a explosão no mapa
				// O uso do canal AcessoMapa garante a exclusão mútua. 
				jogo.AcessoMapa <- AtualizacaoMapa{X: nx, Y: ny, Elem: Explosao}
			}
		}
	}

	// Espera um pouco para o efeito visual da explosão
	time.Sleep(250 * time.Millisecond)

	// Limpa a área da explosão
	for i := -raio; i <= raio; i++ {
		for j := -raio; j <= raio; j++ {
			nx, ny := x+i, y+j
			if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
				
				jogo.AcessoMapa <- AtualizacaoMapa{X: nx, Y: ny, Elem: Vazio}
                
			}
		}
	}

	// Libera a flag para que o jogador possa plantar outra bomba
	jogo.BombaAtiva = false
}