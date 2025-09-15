// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"sync"
)

var mu sync.Mutex

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

type Entidade struct {
	Sprite         Elemento
	X, Y           int
	UltimoVisitado Elemento
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa      [][]Elemento // grade 2D representando o mapa
	StatusMsg string       // mensagem para a barra de status
	Entidades []Entidade   // posicoes dos inimigos e jogador ([0] é o jogador)
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{}
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
				ent := Entidade{X: x, Y: y, UltimoVisitado: e, Sprite: Inimigo}
				jogo.Entidades = append(jogo.Entidades, ent) // Adiciona inimigo
				e = Vazio
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				ent := Entidade{X: x, Y: y, UltimoVisitado: e, Sprite: Personagem}
				jogo.Entidades = append([]Entidade{ent}, jogo.Entidades...) // Adiciona personagem no início
				e = Vazio
				// O personagem é o primeiro elemento em jogo.Entidades[0]
				// Outros inimigos são adicionados a partir do índice 1
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	//jogo.StatusMsg = fmt.Sprintf("%d inimigos lidos", len(jogo.Entidades)-1)

	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Verifica se já existe alguma entidade nessa posição
	for _, ent := range jogo.Entidades {
		if ent.X == x && ent.Y == y {
			return false
		}
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int, ent *Entidade) {
	mu.Lock()
	defer mu.Unlock()

	// Calcula nova posição
	nx, ny := x+dx, y+dy
	jogo.Mapa[y][x] = ent.UltimoVisitado

	ent.UltimoVisitado = jogo.Mapa[ny][nx]
	// Atualiza posição da entidade
	ent.X, ent.Y = nx, ny

	/*// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento*/

}
