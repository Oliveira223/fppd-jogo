// jogo.go - Funções para manipular os elementos do jogo
package main

import (
    "bufio"
    "os"
)

type AtualizacaoMapa struct {
    X, Y int
    Elem Elemento
}

type Elemento struct {
    simbolo  rune
    cor      Cor
    corFundo Cor
    tangivel bool
}

type Jogo struct {
    Mapa           [][]Elemento
    PosX, PosY     int
    Direcao        rune
    UltimoVisitado Elemento
    StatusMsg      string
    AcessoMapa     chan AtualizacaoMapa
    BombaAtiva     bool
}

var (
    Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
    Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
    Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
    Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
    Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
    Bomba      = Elemento{'o', CorVermelho, CorPadrao, true}
    Explosao   = Elemento{'*', CorVerde, CorPadrao, false}
    Direcao    = Elemento{'·', CorCinzaEscuro, CorPadrao, false}
)

func jogoNovo() Jogo {
    return Jogo{
        UltimoVisitado: Vazio,
        Direcao:        'd',
        AcessoMapa:     make(chan AtualizacaoMapa),
        BombaAtiva:     false,
    }
}

// ATUALIZADO: Garante que a célula inicial do personagem no mapa seja vazia.
func jogoCarregarMapa(nome string, jogo *Jogo) error {
    arq, err := os.Open(nome)
    if err != nil {
        return err
    }
    defer arq.Close()

    scanner := bufio.NewScanner(arq)
    y := 0
    for scanner.Scan() {
        linha := scanner.Text()
        var linhaElems []Elemento
        for x, ch := range linha {
            e := Vazio
            switch ch {
            case Parede.simbolo:
                e = Parede
            case Inimigo.simbolo:
                e = Inimigo
            case Vegetacao.simbolo:
                e = Vegetacao
            case Personagem.simbolo:
                jogo.PosX, jogo.PosY = x, y
                // **ADIÇÃO IMPORTANTE**: Garante que o chão sob o personagem seja 'Vazio'.
                // O personagem agora é uma entidade separada do mapa de dados.
                e = Vazio
            }
            linhaElems = append(linhaElems, e)
        }
        jogo.Mapa = append(jogo.Mapa, linhaElems)
        y++
    }
    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}

// Verifica se o personagem pode se mover para a posição (x, y) (sem alterações)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
    if y < 0 || y >= len(jogo.Mapa) {
        return false
    }
    if x < 0 || x >= len(jogo.Mapa[y]) {
        return false
    }
    if jogo.Mapa[y][x].tangivel {
        return false
    }
    return true
}

// REMOVIDO: A função jogoMoverElemento não é mais necessária com a nova lógica.
// Você pode apagar a função inteira ou apenas deixá-la sem uso.
/*
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
    nx, ny := x+dx, y+dy
    elemento := jogo.Mapa[y][x]
    jogo.Mapa[y][x] = jogo.UltimoVisitado
    jogo.UltimoVisitado = jogo.Mapa[ny][nx]
    jogo.Mapa[ny][nx] = elemento
}
*/
