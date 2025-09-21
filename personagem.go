// personagem.go - Funções para movimentação e ações do personagem
package main

import (
    //"fmt"
    "time"
)

// Atualiza a posição do personagem e sua direção.
func personagemMover(tecla rune, jogo *Jogo) {
    dx, dy := 0, 0
    switch tecla {
    case 'w':
        dy = -1
        jogo.Direcao = 'w' 
    case 'a':
        dx = -1
        jogo.Direcao = 'a' 
    case 's':
        dy = 1
        jogo.Direcao = 's' 
    case 'd':
        dx = 1
        jogo.Direcao = 'd' 
    }

    nx, ny := jogo.PosX+dx, jogo.PosY+dy
    // verifica se pode mover e atualiza as coordenadas do personagem.
    if jogoPodeMoverPara(jogo, nx, ny) {
        jogo.PosX, jogo.PosY = nx, ny
    }
}

// Função para atirar
func atirar(jogo *Jogo, x, y int) {
    dx, dy := 0, 0
    switch jogo.Direcao {
    case 'w':
        dy = -1
    case 'a':
        dx = -1
    case 's':
        dy = 1
    case 'd':
        dx = 1
    }

    for i := 1; i <= 30; i++ {
        tiroX, tiroY := x+(dx*i), y+(dy*i)
        if tiroY < 0 || tiroY >= len(jogo.Mapa) || tiroX < 0 || tiroX >= len(jogo.Mapa[tiroY]) {
            break
        }
        if jogo.Mapa[tiroY][tiroX].tangivel {
            break
        }
        jogo.Mapa[tiroY][tiroX] = Vegetacao
        interfaceDesenharJogo(jogo)
        time.Sleep(100 * time.Millisecond)
        jogo.Mapa[tiroY][tiroX] = Vazio
    }
}

// A função de interação planta a bomba. A mensagem de status foi removida para não ser sobreposta.
func personagemInteragir(jogo *Jogo) {
    if !jogo.BombaAtiva {
        jogo.BombaAtiva = true
        bombaX, bombaY := jogo.PosX, jogo.PosY

        // Envia uma solicitação para colocar a bomba no mapa
        jogo.AcessoMapa <- AtualizacaoMapa{X: bombaX, Y: bombaY, Elem: Bomba}
        jogo.StatusMsg = "Bomba plantada! Corra!"

        // Inicia a goroutine da bomba, que cuidará da explosão
        go gerenciarBomba(jogo, bombaX, bombaY)

    } else {
        jogo.StatusMsg = "Aguarde a bomba anterior explodir!"
    }
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
    switch ev.Tipo {
    case "sair":
        return false
    case "interagir":
        personagemInteragir(jogo)
    case "mover":
        personagemMover(ev.Tecla, jogo)
    }
    return true
}
