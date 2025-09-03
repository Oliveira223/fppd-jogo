// main.go - Loop principal do jogo
package main

import (
	"os"
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
	//receber inimigos para a execucao

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
			evento := interfaceLerEventoTeclado()
		select{
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		inimigoExecutarAcao(&jogo)

		interfaceDesenharJogo(&jogo)
	}
}
