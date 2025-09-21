// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"	
)
func gerenciadorDeEstado(jogo *Jogo, redraw chan bool) {
	for {
		// Espera por uma solicitação de atualização vinda de qualquer goroutine
		update := <-jogo.AcessoMapa
		
		// Aplica a atualização no mapa de forma segura
		jogo.Mapa[update.Y][update.X] = update.Elem
		
		// Envia um sinal para o loop principal redesenhar a tela
		select {
		case redraw <- true:
		default:
		}
	}
}

func main() {
	// Inicializa a interface 
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Canal para eventos de teclado
	eventoTeclado := make(chan EventoTeclado)
	go func() {
		for {
			eventoTeclado <- interfaceLerEventoTeclado()
		}
	}()
	
	// Canal para sinalizar que a tela precisa ser redesenhada
	redraw := make(chan bool, 1) // Buffer de 1 para não bloquear

	// Inicia o gerenciador de estado em uma goroutine separada
	go gerenciadorDeEstado(&jogo, redraw)
	
	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// O loop principal agora escuta múltiplos canais
	for {
		select {
		case ev := <-eventoTeclado: // Escuta por ações do jogador
			if continuar := personagemExecutarAcao(ev, &jogo); !continuar {
				return // Sai do programa
			}
			interfaceDesenharJogo(&jogo)

		case <-redraw: // Escuta por pedidos de redesenho (ex: da bomba)
			interfaceDesenharJogo(&jogo)
		
        // Ticker para manter a barra de status sempre atualizada
        case <-time.After(100 * time.Millisecond):
             interfaceDesenharJogo(&jogo)
		}
	}
}