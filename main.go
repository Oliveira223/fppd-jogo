// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Inicializa logs para cada inimigo (exceto o personagem que é índice 0)
	jogo.LogsInimigos = make([]string, len(jogo.Entidades)-1)
	for i := range jogo.LogsInimigos {
		jogo.LogsInimigos[i] = "Aguardando..."
	}

	///////////////////////////////////////////////////////////////////////////////////

	// CANAL PARA GERENCIAR VIDA DO PERSONAGEM
	chanVida := make(chan int) // Canal para gerenciar vida do personagem
	go func() {
		for {
			select {
			case v := <-chanVida:
				jogo.Vida += v
				if jogo.Vida > 5 {
					jogo.Vida = 5 // Limita vida máxima
				}
				jogo.StatusMsg = "Vida aumentada!"
			}
		}
	}()

	///////////////////////////////////////////////////////////////////////////////////

	//BOMBAS ALEATORIAS

	// Canal para posição do personagem
	chanPosPersonagem := make(chan [2]int, 1)
	// Canais para posição dos inimigos
	chansPosInimigos := make([]chan [2]int, len(jogo.Entidades)-1)
	for i := range chansPosInimigos {
		chansPosInimigos[i] = make(chan [2]int, 1)
	}
	// Goroutine para enviar posição do personagem

	go func() {
		for {
			chanPosPersonagem <- [2]int{jogo.Entidades[0].X, jogo.Entidades[0].Y}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Goroutines para enviar posição dos inimigos
	for i := range chansPosInimigos {
		go func(idx int, ch chan<- [2]int) {
			for {
				if idx+1 < len(jogo.Entidades) {
					ch <- [2]int{jogo.Entidades[idx+1].X, jogo.Entidades[idx+1].Y}
				}
				time.Sleep(200 * time.Millisecond)
			}
		}(i, chansPosInimigos[i])
	}

	go iniciarBombasInteligentes(&jogo, 3, 5, chanPosPersonagem, chansPosInimigos, chanVida) // 3 bombas a cada 5 segundos

	///////////////////////////////////////////////////////////////////////////////////

	//INIMIGOS AUTONOMOS
	canais := make([]chan [2]int, len(jogo.Entidades)-1) // Canais para inimigos detectarem o personagem
	for i := range canais {
		canais[i] = make(chan [2]int, 1)
		go func(idx int, ch <-chan [2]int) { //Goroutine para cada inimigo agir independentemente
			for {
				inimigoExecutarAcao(&jogo, idx+1, ch, chanVida) // +1 porque o personagem é o índice 0
				time.Sleep(500 * time.Millisecond)              //Velocidade do inimigo
			}
		}(i, canais[i])
	}

	//POSICAO DO PERSONAGEM ENVIADA PARA INIMIGOS
	go func() { //Envia a posição do personagem para os inimigos periodicamente
		for {
			for _, ch := range canais {
				select {
				case ch <- [2]int{jogo.Entidades[0].X, jogo.Entidades[0].Y}:
					//jogo.StatusMsg = "Posição do personagem enviada para inimigos"
				default:
					//jogo.StatusMsg = "Posição do personagem não enviada"
				}
			}
			time.Sleep(500 * time.Millisecond) //Intervalo de detecção do personagem
		}
	}()

	///////////////////////////////////////////////////////////////////////////////////

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	//CURAS PRESENTES NO MAPA
	go piscarcor(&jogo)
	///////////////////////////////////////////////////////////////////////////////////

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, chanVida); !continuar {
			break
		}

		// Atualiza bombas e explosões
		jogoAtualizarBombas(&jogo)
		jogoAtualizarExplosoes(&jogo)

		interfaceDesenharJogo(&jogo)
	}
}
