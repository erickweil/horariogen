package sudoku

import "fmt"

type Possib struct {
	index int // índice da possibilidade
	p []bool // array de possibilidades, true para indicar que é possível
}

func (possibs *Possib) contar() int {
	n := 0
	for i := 0; i < len(possibs.p); i++ {
		if possibs.p[i] { n++ }
	}
	return n 
}

func (possibs *Possib) receber(outro *Possib) {
	possibs.index = outro.index
	for i := 0; i < len(possibs.p); i++ {
		possibs.p[i] = outro.p[i]
	}
}

func (possibs *Possib) resetar() {
	for i := 0; i < len(possibs.p); i++{
		possibs.p[i] = true
	}
}

func iniciarPossib(quadro []int, nPossibs int) []Possib {
	quadro_possib := make([]Possib, len(quadro))
	for i := 0; i < len(quadro_possib); i++ {
		quadro_possib[i] = Possib{i,make([]bool, nPossibs)}

		// Começa com todas as possibilidades são possíveis
		for j := 0; j < nPossibs; j++{
			quadro_possib[i].p[j] = true
		}
	}
	return quadro_possib
}

func printarPossib(quadro_possib []Possib) {
	for i := 0; i < len(quadro_possib); i++ {
		p := &quadro_possib[i]
		fmt.Print("[")
		for k := 0; k < len(p.p); k++ {
			if p.p[k] {
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
