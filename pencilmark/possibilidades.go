package pencilmark
/**
Algoritmo para solucionar Sudoku e outros problemas que imita uma forma de resolver à mão

Basicamente a forma mais simples é se analisar todo o quadro marcando (Pencil mark) que possibilidades 
cada quadrado possui, de acordo com as regras do jogo. Então provavelmente haverá uma opção
com apenas uma possibilidade, escolhe esta possibilidade e continua o processo até resolver todo o quadro.

Porém, é possível que não exista nenhuma opção com apenas uma possibilidade, neste caso se utiliza da
técnica de backtracking, basicamente a diferença entre este método e o método backtracking é que 
a busca é feita escolhendo quadrados com o menor número de possibilidades e não apenas linearmente. Isto
proporciona uma melhora de desempenho, tanto que em problemas fáceis resolve em apenas algumas centenas de iterações
(E em problemas difíceis pode demorar tanto quanto o método backtracking comum)
*/

import (
	"fmt"
	"math/rand"
	//"sort"
)

type Possib struct {
	Index int // índice da possibilidade
	P []bool // array de possibilidades, true para indicar que é possível
}

func (possibs *Possib) Contar() int {
	n := 0
	for i := 0; i < len(possibs.P); i++ {
		if possibs.P[i] { n++ }
	}
	return n 
}

func (possibs *Possib) Receber(outro *Possib) {
	possibs.Index = outro.Index
	for i := 0; i < len(possibs.P); i++ {
		possibs.P[i] = outro.P[i]
	}
}

func (possibs *Possib) Resetar(valor bool) {
	for i := 0; i < len(possibs.P); i++{
		possibs.P[i] = valor
	}
}

func IniciarPossib(quadro []int, nPossibs int) []Possib {
	quadro_possib := make([]Possib, len(quadro))
	for i := 0; i < len(quadro_possib); i++ {
		quadro_possib[i] = Possib{i,make([]bool, nPossibs)}

		// Começa com todas as possibilidades são possíveis
		for j := 0; j < nPossibs; j++{
			quadro_possib[i].P[j] = true
		}
	}
	return quadro_possib
}

func PrintarPossib(quadro_possib []Possib) {
	for i := 0; i < len(quadro_possib); i++ {
		p := &quadro_possib[i]
		fmt.Print("[")
		for k := 0; k < len(p.P); k++ {
			if p.P[k] {
				fmt.Print(k+1," ")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("], ")

		if (i+1) % 9 == 0 {
			fmt.Println()
		}
	}
}

type RegrasQuadro func(quadro []int, possibs *Possib) 

// Mantém apenas as possibilidades válidas
func atualizarQuadroPossib(quadro []int, quadro_possib []Possib, regrasfn RegrasQuadro) {
	regrasfn(quadro,nil)
	for i := 0; i < len(quadro_possib); i++  {
		possibs := &quadro_possib[i]
	
		regrasfn(quadro,possibs)
	}
}

// Ao mesmo tempo que analisa as possibilidades, encontra a com menor entropia
func obterMelhorPossib(quadro []int, nPossibs int, regrasfn RegrasQuadro) (*Possib, error) {
	var p *Possib = &Possib{-1,make([]bool, nPossibs)}
	var min_p *Possib = &Possib{-1,make([]bool, nPossibs)}
	var min_cont int = -1

	// Para atualizar o cache das checagens
	regrasfn(quadro,nil)

	for index := 0; index < len(quadro); index++ {
		if quadro[index] != 0 { continue } // se já foi escolhido ignora
		p.Resetar(true)
		p.Index = index
		
		regrasfn(quadro,p)

		cont := p.Contar()
		if min_cont == -1 || cont < min_cont {
			min_cont = cont
			min_p.Receber(p) // copia os valores
		}
	}

	if min_cont <= 0 {
		return nil,fmt.Errorf("Não encontrou nenhuma possibilidade")
	}

	return min_p,nil
}

// Se não tem nenhum quadrado sem escolher, está solucionado
func checarSolucionado(quadro []int) bool {
	for i := 0; i < len(quadro); i++ {
		if quadro[i] == 0 { return false }
	}
	return true
}


func getRandomRange(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func SolucionarQuadro(quadro []int,nPossibs int, regrasfn RegrasQuadro) (int, bool) {
	var iter int = 0
	return iter,_solucionarQuadro(0,&iter,quadro,nPossibs,regrasfn)
}

func _solucionarQuadro(depth int, iter *int, quadro []int,nPossibs int, regrasfn RegrasQuadro) bool {
	*iter++
	if *iter % 10000 == 0 {
		fmt.Printf("iter: %d\n", *iter)
	}

	p, err := obterMelhorPossib(quadro,nPossibs,regrasfn)

	if err != nil {
		//fmt.Println("Quadro impossível...")
		return false		
	}

	// Uma vez escolhido o quadrado a partir do qual continuar,
	// testa cada possibilidade deste quadrado
	//for k := 0; k < len(p.P); k++ {
	//for k := len(p.P)-1; k >= 0; k-- {
	randRange := getRandomRange(make([]int, len(p.P)))
	for _, k := range randRange {
		// Se é possível colocar o valor k neste quadrado
		if p.P[k] { 
			quadro[p.Index] = k+1
			// Se com essa escolha já solucionou, retorna true
			if checarSolucionado(quadro) {
				return true
			}
			// Tenta solucionar com mais escolhas depois dessa, e se der certo retorna true
			if _solucionarQuadro(depth+1,iter,quadro,nPossibs,regrasfn) {
				return true
			}

			// Essa escolha não resolveu o sudoku, remove ela para tentar outras
			quadro[p.Index] = 0
		}
	}
	
	// Nenhuma escolha foi válida para este quadrado, ou seja, é ímpossível solucionar nesta configuração
	//fmt.Println("Backtracking...")
	return false
}

func SolucionarQuadroSemParar(quadro []int, iter *int,nPossibs int, regrasfn RegrasQuadro, results [][]int) [][]int {
	//iter++

	*iter++
	if *iter % 10000 == 0 {
		fmt.Printf("iter: %d Soluções: %d\n", *iter, len(results))
	}

	if len(results) > 100{
		fmt.Println("Já deu né! Parando...")
		return results
	}
	// A ideia é obter o quadrado com menor entropia
	// isto é, que possui a menor quantidade de escolhas possíveis
	p, err := obterMelhorPossib(quadro,nPossibs,regrasfn)

	if err != nil {
		//fmt.Println("Quadro impossível...")
		return results		
	}

	//for k := 0; k < len(p.p); k++ {
	//for k := len(p.P)-1; k >= 0; k-- {

	randRange := getRandomRange(make([]int, len(p.P)))
	for _, k := range randRange {
		// Se é possível colocar o valor k neste quadrado
		if p.P[k] { 
			quadro[p.Index] = k+1
			// Se com essa escolha já solucionou, adiciona o quadro resolvido nos resultados
			if checarSolucionado(quadro) {
				quadroCopia := make([]int, len(quadro))
				for q := 0; q < len(quadro); q++ {
					quadroCopia[q] = quadro[q]
				}
				results = append(results, quadroCopia)
				//fmt.Println("Solucionado! iter:",iter," Soluções:",len(results))
			} else {
				// Se não está solucionado, tenta solucionar com mais escolhas depois dessa
				results = SolucionarQuadroSemParar(quadro,iter,nPossibs,regrasfn,results)
			}
			// Remove a escolha para tentar outras possibilidades
			quadro[p.Index] = 0
		}
	}
	// retorna os resultados encontrados
	return results
}