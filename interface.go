package main
//===============================================================================
//IMPORTS E DEFINIÇÕES DE TIPOS
//===============================================================================

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsular as cores do termbox
type Cor = termbox.Attribute

// Representa uma ação detectada do teclado
type EventoTeclado struct {
	Tipo  string // Tipos: "sair", "mover", "bomba"
	Tecla rune   // Tecla pressionada (usado para movimento)
}

// ==============================================================================
//  CONSTANTES E CONFIGURAÇÕES DE CORES
// =============================================================================
// COres do jogo
const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorVerde           = termbox.ColorGreen
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
	CorBranco          = termbox.ColorWhite
)

// =============================================================================
//  MÓDULO DE INICIALIZAÇÃO DA INTERFACE
// =============================================================================

// Inicializa o sistema de interface gráfica
// Deve ser chamada antes de qualquer operação de desenho
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(fmt.Sprintf("Erro ao inicializar interface: %v", err))
	}
}

// Encerra o sistema de interface gráfica
func interfaceFinalizar() {
	termbox.Close()
}

// =============================================================================
// MÓDULO DE CAPTURA DE EVENTOS
// =============================================================================

// Captura e processa eventos do teclado
func interfaceLerEventoTeclado() EventoTeclado {
	// Aguarda um evento do sistema
	ev := termbox.PollEvent()
	
	// Processa apenas eventos de teclado
	if ev.Type != termbox.EventKey {
		return EventoTeclado{} // Retorna evento vazio para outros tipos
	}
	
	// Detecta tecla ESC para sair do jogo
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	
	// Detecta tecla E para colocar bomba
	if ev.Ch == 'e' || ev.Ch == 'E' {
		return EventoTeclado{Tipo: "bomba", Tecla: ev.Ch}
	}
	
	// Para outras teclas, retorna como evento de movimento
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// ==============================================================================
// MÓDULO DE RENDERIZAÇÃO PRINCIPAL
// =============================================================================


// Função principal de renderização que coordena todos os elementos visuais
func interfaceDesenharJogo(jogo *Jogo) {
	// Limpa a tela antes de desenhar
	interfaceLimparTela()
	
	// Se o jogo terminou, mostra apenas a tela de fim de jogo
	if jogo.JogoTerminado {
		interfaceDesenharTelaFimJogo(jogo)
	} else {
		// Renderiza o mapa base
		interfaceRenderizarMapa(jogo)
		
		// Renderiza bombas (antes das entidades para que fiquem atrás)
		interfaceRenderizarBombas(jogo)
		
		// Renderiza as entidades (jogador e inimigos)
		interfaceRenderizarEntidades(jogo)
		
		// Renderiza explosões (por cima de tudo)
		interfaceRenderizarExplosoes(jogo)
		
		// Renderiza elementos de interface
		interfaceDesenharIndicadorDirecao(jogo)
		interfaceDesenharBarraDeStatus(jogo)
	}
	
	// Atualiza a tela com todas as mudanças
	interfaceAtualizarTela()
}

// Desenha todos os eleentos do mapa
func interfaceRenderizarMapa(jogo *Jogo) {
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}
}

// Desenha todas as entidades do jogo
func interfaceRenderizarEntidades(jogo *Jogo) {
	for i := 0; i < len(jogo.Entidades); i++ {
		entidade := &jogo.Entidades[i]
		interfaceDesenharElemento(entidade.X, entidade.Y, entidade.Sprite)
	}
}

// Desenha todas as bombas ativas
func interfaceRenderizarBombas(jogo *Jogo) {
	for _, bomba := range jogo.Bombas {
		if bomba.Ativa {
			interfaceDesenharElemento(bomba.X, bomba.Y, BombaElem)
		}
	}
}

// Desenha todas as explosões ativas
func interfaceRenderizarExplosoes(jogo *Jogo) {
	for _, explosao := range jogo.Explosoes {
		if explosao.Ativa {
			interfaceDesenharElemento(explosao.X, explosao.Y, ExplosaoElem)
		}
	}
}

// ==============================================================================
// MÓDULO DE DESENHO DE ELEMENTOS
// =============================================================================

// Limpa completamente a tela do terminal -> não sei como
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização visual da tela
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento específico na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Mostra a direção atual do personagem
// Desenha um indicador visual na direção que o jogador está olhando
func interfaceDesenharIndicadorDirecao(jogo *Jogo) {
	// Calcula a posição do indicador baseado na direção atual
	dx, dy := 0, 0
	switch jogo.Direcao {
	case 'w': dy = -1 // Cima
	case 'a': dx = -1 // Esquerda
	case 's': dy = 1  // Baixo
	case 'd': dx = 1  // Direita
	}
	
	indicadorX := jogo.Entidades[0].X + dx
	indicadorY := jogo.Entidades[0].Y + dy
	
	// Verifica se há um inimigo na posição do indicador
	temInimigo := false
	for i := 1; i < len(jogo.Entidades); i++ {
		if jogo.Entidades[i].X == indicadorX && jogo.Entidades[i].Y == indicadorY {
			temInimigo = true
			break
		}
	}
	
	// Desenha o indicador apenas se a posição for válida e livre
	if interfacePosicaoValida(jogo, indicadorX, indicadorY) && !temInimigo {
		interfaceDesenharElemento(indicadorX, indicadorY, Direcao)
	}
}

// Verifica se uma posição é válida para desenhar
func interfacePosicaoValida(jogo *Jogo, x, y int) bool {
	return y >= 0 && y < len(jogo.Mapa) && 
		   x >= 0 && x < len(jogo.Mapa[y]) && 
		   !jogo.Mapa[y][x].tangivel
}

// =============================================================================
// MÓDULO DE INTERFACE DE STATUS
// =============================================================================

// Eexibe informações importantes do jogo
// Inclui logs de inimigos, vida do jogador e instruções de controle
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	linhaBase := len(jogo.Mapa) + 2
	
	// Desenha logs dos inimigos
	interfaceDesenharLogsInimigos(jogo, linhaBase)
	
	// Desenha barra de vida
	linhaVida := linhaBase + len(jogo.LogsInimigos) + 1
	interfaceDesenharBarraVida(jogo, linhaVida)
	
	// Desenha instruções de controle
	linhaInstrucoes := linhaVida + 2
	interfaceDesenharInstrucoes(linhaInstrucoes)
	
	// Desenha mensagem de fim de jogo se o jogo terminou
	if jogo.JogoTerminado {
		linhaFimJogo := linhaInstrucoes + 1
		interfaceDesenharMensagemFimJogo(jogo, linhaFimJogo)
	}
}

// Exibe os logs de atividade dos inimigos
func interfaceDesenharLogsInimigos(jogo *Jogo, linhaInicial int) {
	for idx, log := range jogo.LogsInimigos {
		linha := linhaInicial + idx
		
		// Desenha rótulo do inimigo
		rotulo := fmt.Sprintf("Inimigo %d: ", idx+1)
		interfaceDesenharTexto(0, linha, rotulo, CorTexto)
		
		// Desenha log do inimigo
		interfaceDesenharTexto(len(rotulo), linha, log, CorVerde)
	}
}

// Exibe a vida atual do jogador
func interfaceDesenharBarraVida(jogo *Jogo, linha int) {
	// Desenha texto "Vida: "
	vidaTexto := "Vida: "
	interfaceDesenharTexto(0, linha, vidaTexto, CorTexto)
	
	// Desenha corações representando a vida
	for i := 0; i < jogo.Vida; i++ {
		termbox.SetCell(len(vidaTexto)+i, linha, '♥', CorVermelho, CorPadrao)
	}
}

// Exibe as instruções de controle do jogo
func interfaceDesenharInstrucoes(linha int) {
	instrucoes := "Use WASD para mover. E para bomba. ESC para sair."
	interfaceDesenharTexto(0, linha, instrucoes, CorTexto)
}

// Exibe a mensagem de fim de jogo (vitória ou derrota)
func interfaceDesenharMensagemFimJogo(jogo *Jogo, linha int) {
	// Determina a cor baseada no tipo de mensagem
	var cor Cor
	if jogo.StatusMsg == "VITORIA! Todos os inimigos foram eliminados!" {
		cor = CorVerde
	} else {
		cor = CorVermelho
	}
	
	// Desenha a mensagem de status
	interfaceDesenharTexto(0, linha, jogo.StatusMsg, cor)
}

// Desenha uma tela de fim de jogo que substitui o mapa
func interfaceDesenharTelaFimJogo(jogo *Jogo) {
	// Determina a cor baseada no tipo de mensagem
	var cor Cor
	if jogo.StatusMsg == "VITORIA! Todos os inimigos foram eliminados!" {
		cor = CorVerde
	} else {
		cor = CorVermelho
	}
	
	// Centraliza a mensagem na tela
	largura, altura := termbox.Size()
	linhaMsg := altura / 2
	colunaMsg := (largura - len(jogo.StatusMsg)) / 2
	
	// Desenha a mensagem principal centralizada
	interfaceDesenharTexto(colunaMsg, linhaMsg, jogo.StatusMsg, cor)
	
	// Desenha instruções para sair
	instrucaoSair := "Pressione ESC para sair"
	colunaSair := (largura - len(instrucaoSair)) / 2
	interfaceDesenharTexto(colunaSair, linhaMsg+2, instrucaoSair, CorTexto)
}

// Função auxiliar para desenhar texto na tela
func interfaceDesenharTexto(x, y int, texto string, cor Cor) {
	for i, c := range texto {
		termbox.SetCell(x+i, y, c, cor, CorPadrao)
	}
}

// ============================================================================
// MÓDULO DE TELA DE FIM DE JOGO
// ============================================================================

// Mostra a tela de fim de jogo com mensagem de vitória ou derrota
func interfaceMostrarFimDeJogo(vitoria bool) {
	// Limpa a tela
	interfaceLimparTela()
	
	// Define a mensagem baseada no resultado
	var mensagem string
	var cor Cor
	if vitoria {
		mensagem = "VITORIA"
		cor = CorVerde
	} else {
		mensagem = "DERROTA"
		cor = CorVermelho
	}
	
	// Calcula posição central da tela
	largura, altura := termbox.Size()
	x := (largura - len(mensagem)) / 2
	y := altura / 2
	
	// Desenha a mensagem centralizada
	interfaceDesenharTexto(x, y, mensagem, cor)
	
	// Desenha instruções para sair
	instrucao := "Pressione qualquer tecla para sair..."
	xInstrucao := (largura - len(instrucao)) / 2
	interfaceDesenharTexto(xInstrucao, y+2, instrucao, CorTexto)
	
	// Atualiza a tela
	interfaceAtualizarTela()
	
	// Aguarda uma tecla ser pressionada
	termbox.PollEvent()
}

