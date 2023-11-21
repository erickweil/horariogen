package genetic

/*	https://www.youtube.com/watch?v=kHyNqSnzP8Y
	População -> Mutação -> Crossover -> Fenótipos -> Fitness -> Probabilidade -> Seleção

	https://www.youtube.com/watch?v=Fdk7ZKJHFcI&t=537s

	População Inicial
	enquanto não parar:
		Seleção
		Reprodução

	1. Inicialização
	2. Seleção de Pais
	3. Reprodução: Realizar Crossover & Mutação de descendentes
	4. Descendentes: definir a nova população com descendentes e pais
	5. Fitness: Avaliar, e selecionar os 'melhores'
	(Voltar para etapa 2)

	Inicialização/Representação do genoma:
	binária, números reais, números inteiros, permutação

	Seleção de Pais:
		Random, Roleta (Fitness Proportionate Selection), Tournament, Rank

	Reprodução:
	Mutação, CrossOver
		Mutação:
			disturbance(bit flip, small value, reset), swap, scramble, inversion
		CrossOver:
			single-point, n-point, uniform, arithmetic recombination, ox1, ox2, shuffle, ring, equivalent swap

	Sobrevivência para próxima Geração:
		random, fitness(veja Seleção de Pai), trunk ranked list, por idade

	Condição de Parada
		Nenhuma melhora em N iterações
		Atingiu máximo de iterações
		Função fitness atingiu objetivo

	Em resumo, a partir de uma população inicial cria-se descendentes (Com mutações e crossover), e seleciona para a
	próxima geração os que tiverem mais pontos de acordo com algum critério
*/

import (
	//"fmt"
	//"fmt"
	"math"
	"math/rand"
	"sort"
	//"sort"
)

// Representação do Genoma, como cada alelo é representado
type TipoGenoma int
const (
	BINARIO TipoGenoma = iota
	REAL
	INTEIRO
	PERMUTACAO
)

// Como os Pais são escolhidos para reprodução e como Decide a Sobrevivência para próxima geração
type TipoSelecao int
const (
	// Escolhe aleatoriamente
	RANDOM TipoSelecao = iota
	// Fitness Proportionate Selection, mais chances quanto mais fitness
	ROLETA
	// Simula um torneio e os N vencedores são escolhidos
	TORNEIO
	// Parecido com Roleta mas as chances dependem apenas da posição relativa e não do valor cru do fitness
	POSICAO
	// Trunked Ranked List, retira apenas o N melhores da lista
	MELHORES
	// A seleção é aleatória mas indivíduos são removidos após N gerações
	IDADE
)

/*
Tabela de compatibilidade:
			BINARIO		REAL		INTEIRO		PERMUTACAO
PERTURBACAO    x         x             x
RESET          x         x             x
TROCA          x         x             x            x
EMBARALHAR     x         x             x            x
REVERSAO       x         x             x            x
*/
type TipoMutacao int
const (
	// Uma pequena mudança é aplicada. (Em binário é um bit flip)
	PERTURBACAO TipoMutacao = iota
	// Resetado e gerado aleatoriamente denovo.
	RESET
	// Dois alelos são trocados entre si
	TROCA
	// Uma seção do genoma é embaralhada
	EMBARALHAR
	// Uma seção do genoma fica na ordem reversa.
	REVERSAO
)

type TipoCrossOver int
const (
	A TipoCrossOver = iota
	// Single-Point, um único ponto é escolhido para cross over
	UNICO_PONTO
	// N-Point, vários pontos são escolhidos para cross over
	N_PONTOS
	// Uniform, cada alelo é escolhido dos dois pais de forma aleatória
	UNIFORME
	// Faz tipo uma media ponderada. Só funciona com valores real e inteiro
	RECOMBINACAO_ARITMETICA
	OX1
	OX2
	EMBARALHAMENTO
	ANEL
	TROCA_EQUIVALENTE
)

type Cromossomo struct {
	genoma []int
	fitness float64
}

type Populacao struct {
	criaturas []Cromossomo
	tamanhoPopulacao int
	fitnessSoma float64
	tipo TipoGenoma
	selecaoPais TipoSelecao
	selecaoSobreviventes TipoSelecao
	mutacao TipoMutacao
	crossover TipoCrossOver
}

// Tipo que um dot product entre os cromossomos
func calcularDifferenca(a *Cromossomo, b *Cromossomo) float64 {
	diff := 0.0
	for i := 0; i < len(a.genoma); i++ {
		d := float64(a.genoma[i] - b.genoma[i])
		diff += math.Abs(d)
	}
	return diff
}

func (p *Populacao) calcularFitness(calcFitness FuncaoFitness) {
	fitnessSoma := 0.0
	for i := 0; i < len(p.criaturas); i++ {
		criatura := &p.criaturas[i]
		fitness := calcFitness(criatura)
		
		/*// aplicar operador de diversidade aumentando o fitness de indivíduos diferentes
		totalDiff := 0.0
		for k := 0; k < len(p.criaturas); k++ {
			// Cálculo da diferença entro o genoma
			totalDiff += calcularDifferenca(criatura, &p.criaturas[k])
		}
		// Quanto mais diferente, mais fitness
		fitness += 0.00005 * totalDiff*/

		criatura.fitness = fitness

		fitnessSoma += fitness
	}

	p.fitnessSoma = fitnessSoma
}

// Random Fitness Proportionate Selection
func (p *Populacao) selecionarPai(criaturas []Cromossomo,outro *Cromossomo) (*Cromossomo,int) {
	var ultimoNaoIgual *Cromossomo
	var ultimoI = 0
	chance := 1.0 / float64(len(criaturas))
	for i := 0; i < len(criaturas); i++ {
		//criatura := &criaturas[rand.Intn(len(criaturas))]
		criatura := &criaturas[i]
		if criatura == outro { continue }

		ultimoNaoIgual = criatura
		ultimoI = i

		//chance := criatura.fitness / p.fitnessSoma
		if rand.Float64() > chance { 
			continue
		} else { break }
	}
	return ultimoNaoIgual, ultimoI
}

var mutationStep = 8
var mutationChance = 0.25
var sobrevivenciaChance = 0.25

func filhoPorMutacao(criatura *Cromossomo) *Cromossomo{
	filho := &Cromossomo{make([]int, len(criatura.genoma)),0.0}
	for i := 0; i < len(criatura.genoma); i++ {
		if rand.Intn(2) > 0 {
			filho.genoma[i] = criatura.genoma[i] + rand.Intn(mutationStep)
		} else {
			filho.genoma[i] = criatura.genoma[i] - rand.Intn(mutationStep)
		}
	}
	return filho
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func aplicarMutacao(criatura *Cromossomo) {
	if rand.Float64() > mutationChance { return }
	for i := 0; i < len(criatura.genoma); i++ {
		mut := mutationStep + min(rand.Intn(20),rand.Intn(20))

		if rand.Intn(2) > 0 {
			criatura.genoma[i] = criatura.genoma[i] + rand.Intn(mut)
		} else {
			criatura.genoma[i] = criatura.genoma[i] - rand.Intn(mut)
		}
	}
}

func filhoPorCrossOver(pai *Cromossomo, mae *Cromossomo) (*Cromossomo,*Cromossomo) {
	filho1 := &Cromossomo{make([]int, len(pai.genoma)),0.0}
	filho2 := &Cromossomo{make([]int, len(pai.genoma)),0.0}

	// Single Point
	//if len(pai.genoma) == 2 {
	//	filho1.genoma[0] = pai.genoma[0]
	//	filho1.genoma[1] = mae.genoma[1]
	//	
	//	filho2.genoma[0] = mae.genoma[0]
	//	filho2.genoma[1] = pai.genoma[1]
	//} else {
		ponto := rand.Intn(len(pai.genoma)-1)+1 
		for i := 0; i < len(pai.genoma); i++ {
			if i < ponto {
				filho1.genoma[i] = pai.genoma[i]
				filho2.genoma[i] = mae.genoma[i]
			} else {
				filho1.genoma[i] = mae.genoma[i]
				filho2.genoma[i] = pai.genoma[i]
			}
		}
	//}
	
	/*

	*/

	/*// Cross Over Uniforme
	for i := 0; i < len(pai.genoma); i++ {
		if rand.Intn(2) > 0 {
			filho1.genoma[i] = pai.genoma[i]
			filho2.genoma[1] = mae.genoma[i]
		} else {
			filho2.genoma[i] = pai.genoma[i]
			filho1.genoma[1] = mae.genoma[i]
		}
	}*/
	return filho1, filho2
}

type FuncaoFitness func(criatura *Cromossomo) float64

/* 0. Inicializa população & Calcula Fitness
   1. Seleção de Pais      <-----|
   2. CrossOver                  |
   3. Mutação                    |
   4. Cálculo do fitness         |
   5. Seleção de Sobreviventes ---  */
func SimularGeracao(populacao *Populacao, calcFitness FuncaoFitness) *Populacao {
	// 1. Seleção de Pais & 2. CrossOver
	// A ideia é fazer crossOver até completar o dobro da população

	sort.Slice(populacao.criaturas,func(i, j int) bool {
		return populacao.criaturas[i].fitness < populacao.criaturas[j].fitness
	})

	populacaoPais := populacao.criaturas[:]
	for len(populacao.criaturas) < populacao.tamanhoPopulacao*2 {
		pai, _ := populacao.selecionarPai(populacaoPais,nil)
		mae, _ := populacao.selecionarPai(populacaoPais,pai)

		filho1, filho2 := filhoPorCrossOver(pai,mae)

		aplicarMutacao(filho1)
		aplicarMutacao(filho2)

		populacao.criaturas = append(populacao.criaturas, *filho1)
		populacao.criaturas = append(populacao.criaturas, *filho2)
	}

	populacao.calcularFitness(calcFitness)

	//fitnessMedia := populacao.fitnessSoma / float64(len(populacao.criaturas))
	sort.Slice(populacao.criaturas,func(i, j int) bool {
		return populacao.criaturas[i].fitness < populacao.criaturas[j].fitness
	})

	// A ideia é fazer sobreviver só metade do total da população
	k := 0
	for k < populacao.tamanhoPopulacao {
		slicePop := populacao.criaturas[k:]
		selecionado, selecionadoI := populacao.selecionarPai(slicePop,nil)

		if selecionadoI > k {
			// o selecionado será trocado pelo k-th elemento
			temp := populacao.criaturas[k]
			populacao.criaturas[k] = *selecionado
			populacao.criaturas[selecionadoI] = temp
		}
		k++
	}
	populacao.criaturas = populacao.criaturas[0:k]

	return populacao

	/*
	k := 0
	// Selecionar metade para sobreviver
	for k < len(populacao.criaturas) / 2 {
		for i := 0; i < len(populacao.criaturas); i++ {
			criatura := &populacao.criaturas[i]

			chance := criatura.fitness / populacao.fitnessSoma

			if rand.Float64() > chance { continue }
			if populacao.criaturas[i].genoma == nil { continue }

			descendentes.criaturas[k] = populacao.criaturas[i]
			populacao.criaturas[i].genoma = nil
			k++

			break
		}
	}
	*/
	// O melhor indivíduo tem a melhor chance de sobreviver
	// O segundo melhor indivíduo tem a segunda melhor chance de sobreviver
	/*
	for k < len(populacao.criaturas) / 2 {
		for i := 0; i < len(populacao.criaturas); i++ {
			if rand.Float64() > sobrevivenciaChance { continue }
			if populacao.criaturas[i].genoma == nil { continue }

			descendentes.criaturas[k] = populacao.criaturas[i]
			populacao.criaturas[i].genoma = nil
			k++
		}
	}*/
	
	/*i := 0
	for i < len(populacao.criaturas) {
		if rand.Float64() > 1.0-(float64(i)/float64(len(populacao.criaturas))) { 
			i++
			continue 
		}
		descendentes.criaturas[k] = populacao.criaturas[i]
		i++
		k++
	}*/

	/*
	//fmt.Println(k)
	for k < len(descendentes.criaturas) {
		paiAleatorio := &descendentes.criaturas[rand.Intn(k)]
		descendentes.criaturas[k] = *filhoPorMutacao(paiAleatorio)
		k++
	}


	// aplicarCrossOver(descendentes)

	return descendentes
	*/
}

type FuncaoIniciadora func(criatura *Cromossomo)
func CriarPopulacao(n int, 
	genomaN int, 
	tipo TipoGenoma,
	selecaoPais TipoSelecao, 
	selecaoSobreviventes TipoSelecao, 
	mutacao TipoMutacao,
	crossover TipoCrossOver,
	iniciar FuncaoIniciadora,calcFitness FuncaoFitness) *Populacao {

	var populacao = &Populacao{
		make([]Cromossomo, n),
		n,
		0.0,
		tipo,
		selecaoPais,
		selecaoSobreviventes,
		mutacao,
		crossover}

	for i := 0; i < n; i++ {
		populacao.criaturas[i] = Cromossomo{make([]int, genomaN),0.0}
		iniciar(&populacao.criaturas[i])
	}

	populacao.calcularFitness(calcFitness)

	return populacao
}

