// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"time"
)

// Canal para exclusão mútua - funciona como semáforo binário
var mutexChan = make(chan struct{}, 1)

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
	Mapa           [][]Elemento // grade 2D representando o mapa
	Direcao        rune         // direção atual do personagem (w, a, s, d)
	StatusMsg      string       // mensagem para a barra de status
	Entidades      []Entidade   // posicoes dos inimigos e jogador ([0] é o jogador)
	LogsInimigos   []string     // logs de comportamento dos inimigos (aleatório/perseguindo)
	Vida           int          // vida atual do jogador (máximo 3 corações)
	UltimoDano     time.Time    // timestamp do último dano recebido
	CuraUsada      bool         // indica se a cura já foi utilizada (uso único)
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Cura       = Elemento{'+', CorVerde, CorPadrao, false}
	Direcao    = Elemento{'•', CorCinzaEscuro, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{
		Mapa:         make([][]Elemento, 0),
		Direcao:      'w',
		StatusMsg:    "Jogo iniciado",
		Entidades:    make([]Entidade, 0),
		LogsInimigos: make([]string, 0),
		Vida:         3, // jogador começa com 3 corações

	}
}

func piscarcor(jogo *Jogo) {
	pisca := false
	for {
		var novaCor Cor
		if pisca {
			novaCor = CorBranco
		} else {
			novaCor = CorVerde 
		}
		pisca = !pisca

		// Adquire o lock para modificar o mapa
		mutexChan <- struct{}{}
		curaEncontrada := false
		for y := range jogo.Mapa {
			for x := range jogo.Mapa[y] {
				if jogo.Mapa[y][x].simbolo == Cura.simbolo {
					jogo.Mapa[y][x].cor = novaCor
					curaEncontrada = true
				}
			}
		}
		<-mutexChan // Libera o lock

		// Se não há mais curas no mapa, para de piscar
		if !curaEncontrada {
			return
		}

		interfaceDesenharJogo(jogo)
		time.Sleep(1000 * time.Millisecond)
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
			case Cura.simbolo:
				e = Cura
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
	return nil
}

// Verifica se uma entidade pode se mover para a posição (x, y)
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

	// Pode mover para a posição (permitindo sobreposição para dano)
	return true
}

// Verifica se uma entidade pode se mover para a posição, excluindo o personagem
func jogoPodeMoverParaInimigo(jogo *Jogo, x, y int, inimigoIdx int) bool {
	// Verifica limites e tangibilidade
	if !jogoPodeMoverPara(jogo, x, y) {
		return false
	}

	// Verifica se já existe outro inimigo nessa posição (mas permite posição do personagem)
	for i, ent := range jogo.Entidades {
		if i != 0 && i != inimigoIdx && ent.X == x && ent.Y == y {
			return false // Bloqueia movimento para posição de outro inimigo
		}
	}

	return true // Permite movimento para posição do personagem ou vazia
}

// Verifica se o personagem pode se mover para a posição (bloqueia movimento para inimigos)
func jogoPodeMoverParaPersonagem(jogo *Jogo, x, y int) bool {
	// Verifica limites e tangibilidade
	if !jogoPodeMoverPara(jogo, x, y) {
		return false
	}

	// Verifica se já existe algum inimigo nessa posição
	for i, ent := range jogo.Entidades {
		if i != 0 && ent.X == x && ent.Y == y {
			return false // Bloqueia movimento para posição de inimigo
		}
	}

	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int, ent *Entidade) {
	// Adquire o lock usando canal (semáforo binário)
	mutexChan <- struct{}{}
	defer func() { <-mutexChan }() // Libera o lock

	// Calcula nova posição
	nx, ny := x+dx, y+dy
	jogo.Mapa[y][x] = ent.UltimoVisitado

	ent.UltimoVisitado = jogo.Mapa[ny][nx]
	// Atualiza posição da entidade
	ent.X, ent.Y = nx, ny
}
