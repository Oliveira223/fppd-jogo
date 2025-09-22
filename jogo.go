// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"math/rand"
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

// Bomba representa uma bomba no jogo
type Bomba struct {
	X, Y      int       // Posição da bomba
	TempoVida time.Time // Quando a bomba foi colocada
	Ativa     bool      // Se a bomba está ativa
}

// Explosao representa uma explosão temporária
type Explosao struct {
	X, Y      int       // Posição da explosão
	TempoVida time.Time // Quando a explosão começou
	Ativa     bool      // Se a explosão está ativa
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa         [][]Elemento  // grade 2D representando o mapa
	Direcao      rune          // direção atual do personagem (w, a, s, d)
	StatusMsg    string        // mensagem para a barra de status
	Entidades    []Entidade    // posicoes dos inimigos e jogador ([0] é o jogador)
	LogsInimigos []string      // logs de comportamento dos inimigos (aleatório/perseguindo)
	Vida         int           // vida atual do jogador (máximo 3 corações)
	UltimoDano   time.Time     // timestamp do último dano recebido
	CuraUsada    bool          // indica se a cura já foi utilizada (uso único)
	Bombas       []Bomba       // bombas ativas no jogo
	Explosoes    []Explosao    // explosões ativas no jogo
	Comandos     []chan string // canais para enviar comandos aos inimigos
	JogoTerminado bool         // indica se o jogo terminou (vitória ou derrota)
}

// Elementos visuais do jogo
var (
	Personagem   = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo      = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede       = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao    = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio        = Elemento{' ', CorPadrao, CorPadrao, false}
	Cura         = Elemento{'+', CorVerde, CorPadrao, false}
	Direcao      = Elemento{'•', CorCinzaEscuro, CorPadrao, false}
	BombaElem    = Elemento{'●', CorVermelho, CorPadrao, false}
	ExplosaoElem = Elemento{'*', CorVermelho, CorPadrao, false}
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

// ============================================================================
// MÓDULO DE SISTEMA DE BOMBAS
// ============================================================================

// Coloca uma bomba na posição atual do jogador
func jogoColocarBomba(jogo *Jogo) {
	x, y := jogo.Entidades[0].X, jogo.Entidades[0].Y

	// Verifica se já existe uma bomba nesta posição
	for _, bomba := range jogo.Bombas {
		if bomba.X == x && bomba.Y == y && bomba.Ativa {
			jogo.StatusMsg = "Já existe uma bomba aqui!"
			return
		}
	}

	// Cria nova bomba
	novaBomba := Bomba{
		X:         x,
		Y:         y,
		TempoVida: time.Now(),
		Ativa:     true,
	}

	jogo.Bombas = append(jogo.Bombas, novaBomba)
	jogo.StatusMsg = "Bomba colocada!"
}

// Atualiza o estado das bombas (verifica se devem explodir)
func jogoAtualizarBombas(jogo *Jogo) {
	tempoAtual := time.Now()

	for i := len(jogo.Bombas) - 1; i >= 0; i-- {
		bomba := &jogo.Bombas[i]

		if bomba.Ativa && tempoAtual.Sub(bomba.TempoVida) >= 3*time.Second {
			// Bomba explode após 3 segundos
			jogoExplodirBomba(jogo, bomba.X, bomba.Y)
			bomba.Ativa = false

			// Remove bomba da lista
			jogo.Bombas = append(jogo.Bombas[:i], jogo.Bombas[i+1:]...)
		}
	}
}

// Cria uma explosão na posição especificada com raio 5
func jogoExplodirBomba(jogo *Jogo, x, y int) {
	raio := 5
	tempoAtual := time.Now()

	// Cria explosões em todas as direções dentro do raio
	for dx := -raio; dx <= raio; dx++ {
		for dy := -raio; dy <= raio; dy++ {
			// Calcula distância Manhattan (mais apropriada para jogos em grade)
			distancia := abs(dx) + abs(dy)
			if distancia <= raio {
				ex, ey := x+dx, y+dy

				// Verifica se a posição está dentro dos limites do mapa
				if ex >= 0 && ex < len(jogo.Mapa[0]) && ey >= 0 && ey < len(jogo.Mapa) {
					// Não explode através de paredes
					if jogo.Mapa[ey][ex].simbolo != Parede.simbolo {
						explosao := Explosao{
							X:         ex,
							Y:         ey,
							TempoVida: tempoAtual,
							Ativa:     true,
						}
						jogo.Explosoes = append(jogo.Explosoes, explosao)

						// Verifica se há inimigos na posição da explosão
						jogoVerificarInimigoNaExplosao(jogo, ex, ey)
					}
				}
			}
		}
	}

	jogo.StatusMsg = "BOOM! Bomba explodiu!"
}

// Verifica se há inimigos na posição da explosão e os elimina
func jogoVerificarInimigoNaExplosao(jogo *Jogo, x, y int) {
	for i := len(jogo.Entidades) - 1; i >= 1; i-- { // Começa do 1 para não afetar o jogador
		if jogo.Entidades[i].X == x && jogo.Entidades[i].Y == y {
			// Remove inimigo
			jogo.Entidades = append(jogo.Entidades[:i], jogo.Entidades[i+1:]...)

			// Remove log correspondente
			if i-1 < len(jogo.LogsInimigos) {
				jogo.LogsInimigos = append(jogo.LogsInimigos[:i-1], jogo.LogsInimigos[i:]...)
			}

			jogo.StatusMsg = "Inimigo eliminado pela explosão!"
		}
	}
}

// Atualiza o estado das explosões (remove as que expiraram)
func jogoAtualizarExplosoes(jogo *Jogo) {
	tempoAtual := time.Now()

	for i := len(jogo.Explosoes) - 1; i >= 0; i-- {
		explosao := &jogo.Explosoes[i]

		if explosao.Ativa && tempoAtual.Sub(explosao.TempoVida) >= 500*time.Millisecond {
			// Explosão dura 500ms
			explosao.Ativa = false

			// Remove explosão da lista
			jogo.Explosoes = append(jogo.Explosoes[:i], jogo.Explosoes[i+1:]...)
		}
	}
}

// Função auxiliar para calcular valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Função para criar bombas inteligentes aleatórias
func iniciarBombasInteligentes(jogo *Jogo, numBombas, raio int, chanPersonagem chan [2]int, chansInimigos []chan [2]int, chanVida chan int) {
	for {
		for i := 0; i < numBombas; i++ {
			go func() {
				// Gera posição aleatória válida
				var x, y int
				for {
					x = rand.Intn(len(jogo.Mapa[0]))
					y = rand.Intn(len(jogo.Mapa))
					if jogo.Mapa[y][x].simbolo == Vazio.simbolo {
						break
					}
				}
				// Adiciona bomba visualmente (opcional)
				jogo.Bombas = append(jogo.Bombas, Bomba{X: x, Y: y, Ativa: true, TempoVida: time.Now()})

				for {
					select {
					case pos := <-chanPersonagem:
						if inimigoDetectaPersonagem(x, y, pos[0], pos[1]) <= raio {
							jogo.StatusMsg = "Bomba explodiu: personagem próximo!"
							jogoExplodirBomba(jogo, x, y)
							// Dano ao personagem
							mutexChan <- struct{}{}
							if jogo.Vida > 0 {
								chanVida <- -1
								jogo.UltimoDano = time.Now()
							}
							<-mutexChan
							return
						}
					default:
						for _, ch := range chansInimigos {
							select {
							case pos := <-ch:
								if inimigoDetectaPersonagem(x, y, pos[0], pos[1]) <= raio {
									jogo.StatusMsg = "Bomba explodiu: inimigo próximo!"
									jogoExplodirBomba(jogo, x, y)
									// Aqui você pode eliminar o inimigo se quiser
									return
								}
							case <- time.After(2* time.Second):
							}
						}
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()
		}
		time.Sleep(15 * time.Second) // Espera antes de criar novas bombas
	}
}

// ============================================================================
// MÓDULO DE SISTEMA DE VITÓRIA E DERROTA
// ============================================================================

// Verifica e processa condição de vitória (não há mais inimigos no mapa)
func jogoVerificarVitoria(jogo *Jogo) {
	// Só verifica se o jogo ainda não terminou
	if jogo.JogoTerminado {
		return
	}
	
	// Conta quantos inimigos restam (excluindo o jogador que está no índice 0)
	numInimigos := len(jogo.Entidades) - 1
	if numInimigos <= 0 {
		jogo.JogoTerminado = true
		jogo.StatusMsg = "VITORIA! Todos os inimigos foram eliminados!"
	}
}

// Verifica e processa condição de derrota (vida chegou a 0)
func jogoVerificarDerrota(jogo *Jogo) {
	// Só verifica se o jogo ainda não terminou
	if jogo.JogoTerminado {
		return
	}
	
	if jogo.Vida <= 0 {
		jogo.JogoTerminado = true
		jogo.StatusMsg = "DERROTA! Sua vida chegou a zero!"
	}
}
