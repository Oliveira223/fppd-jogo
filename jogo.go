// jogo.go - Funções para manipular os elementos do jogo
package main

import (
	"bufio"
	"os"
)

// NOVO: Estrutura para uma solicitação de atualização do mapa.
// Esta definição estava faltando.
type AtualizacaoMapa struct {
	X, Y  int
	Elem  Elemento
}

// Elemento representa qualquer objeto do mapa
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento
	PosX, PosY      int
	Direcao         rune
	UltimoVisitado  Elemento
	StatusMsg       string
	// --- CAMPOS QUE ESTAVAM FALTANDO ---
	AcessoMapa      chan AtualizacaoMapa // Canal para acesso seguro ao mapa
	BombaAtiva      bool                 // Flag para garantir uma bomba por vez
}

// Elementos visuais do jogo
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

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	return Jogo{
		UltimoVisitado: Vazio,
		Direcao:        'd',
		// --- INICIALIZAÇÃO DOS NOVOS CAMPOS ---
		AcessoMapa:     make(chan AtualizacaoMapa),
		BombaAtiva:     false,
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
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

// Verifica se o personagem pode se mover para a posição (x, y)
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

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy
	elemento := jogo.Mapa[y][x]
	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}