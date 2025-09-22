// personagem.go - Sistema de controle e ações do personagem
package main

// ============================================================================
// MÓDULO DE MOVIMENTAÇÃO E POSICIONAMENTO
// ============================================================================
//import "time"

// Move o personagem baseado na tecla WASD pressionada
func personagemMover(tecla rune, jogo *Jogo, chanVida chan int) {
	dx, dy := 0, 0

	// Define direção baseada na tecla
	switch tecla {
	case 'w':
		dy = -1 // Cima
	case 'a':
		dx = -1 // Esquerda
	case 's':
		dy = 1 // Baixo
	case 'd':
		dx = 1 // Direita
	}

	// Atualiza direção atual do personagem
	jogo.Direcao = tecla

	// Calcula nova posição
	nx, ny := jogo.Entidades[0].X+dx, jogo.Entidades[0].Y+dy

	// Verifica se movimento é válido e executa
	if jogoPodeMoverParaPersonagem(jogo, nx, ny) {
		// Verifica se há cura na posição de destino
		if jogo.Mapa[ny][nx].simbolo == Cura.simbolo {
			coletarCura(jogo, nx, ny, chanVida)
		}
		jogoMoverElemento(jogo, jogo.Entidades[0].X, jogo.Entidades[0].Y, dx, dy, &jogo.Entidades[0])
	}
}

// ============================================================================
// MÓDULO DE SISTEMA DE CURA
// ============================================================================

// Coleta uma cura na posição especificada
func coletarCura(jogo *Jogo, x, y int, chanVida chan int) {
	// Remove cura do mapa
	jogo.Mapa[y][x] = Vazio

	// Marca cura como usada
	jogo.CuraUsada = true

	chanVida <- 1 // Envia sinal para aumentar vida
}

// ============================================================================
// MÓDULO DE PROCESSAMENTO DE EVENTOS
// ============================================================================

// A função de interação planta a bomba. A mensagem de status foi removida para não ser sobreposta.
func personagemInteragir(jogo *Jogo) {
    if !jogo.BombaAtiva {
        jogo.BombaAtiva = true
        bombaX, bombaY := jogo.Entidades[0].X, jogo.Entidades[0].Y

        // Envia uma solicitação para colocar a bomba no mapa
        jogo.AcessoMapa <- AtualizacaoMapa{X: bombaX, Y: bombaY, Elem: Bomba}
        jogo.StatusMsg = "Bomba plantada! Corra!"

        // Inicia a goroutine da bomba, que cuidará da explosão
        go gerenciarBomba(jogo, bombaX, bombaY)

    } else {
        jogo.StatusMsg = "Aguarde a bomba anterior explodir!"
    }
}


// Processa eventos de teclado e executa ações do personagem
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, chanVida chan int) bool {
	switch ev.Tipo {
	case "sair":
		// Termina o jogo
		return false

	case "interagir":
		//Persongaem Coloca Bomba
        personagemInteragir(jogo)

	case "mover":
		// Executa movimento do personagem
		personagemMover(ev.Tecla, jogo, chanVida)

	case "bomba":
		// Coloca uma bomba na posição atual
		jogoColocarBomba(jogo)
	}

	// Continua o jogo
	return true
}
