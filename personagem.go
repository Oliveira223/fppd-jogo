// personagem.go - Funções para movimentação e ações do personagem
package main

import (
    //"fmt"
    "time"
)

// ATUALIZADO: Atualiza a posição do personagem e sua direção.
func personagemMover(tecla rune, jogo *Jogo) {
    dx, dy := 0, 0
    switch tecla {
    case 'w':
        dy = -1
        jogo.Direcao = 'w' // Atualiza a direção do personagem
    case 'a':
        dx = -1
        jogo.Direcao = 'a' // Atualiza a direção do personagem
    case 's':
        dy = 1
        jogo.Direcao = 's' // Atualiza a direção do personagem
    case 'd':
        dx = 1
        jogo.Direcao = 'd' // Atualiza a direção do personagem
    }

    nx, ny := jogo.PosX+dx, jogo.PosY+dy
    // Apenas verifica se pode mover e atualiza as coordenadas do personagem.
    // **NÃO** chama mais jogoMoverElemento, para não apagar o que está no mapa.
    if jogoPodeMoverPara(jogo, nx, ny) {
        jogo.PosX, jogo.PosY = nx, ny
    }
}

// Função para atirar (sem alterações)
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

// ATUALIZADO: A função de interação agora planta a bomba. A mensagem de status foi removida para não ser sobreposta.
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

// Processa o evento do teclado e executa a ação correspondente (sem alterações)
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
