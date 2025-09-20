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

	canais := make([]chan [2]int, len(jogo.Entidades)-1) // Canais para inimigos detectarem o personagem
	for i := range canais {
		canais[i] = make(chan [2]int, 1)
		go func(idx int, ch <-chan [2]int) { //Goroutine para cada inimigo agir independentemente
			for {
				inimigoExecutarAcao(&jogo, idx+1, ch)
				time.Sleep(500 * time.Millisecond) //Velocidade do inimigo
			}
		}(i, canais[i])
	}

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

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	go piscarcor(&jogo)



	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
