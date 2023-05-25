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

func printarQuadro(quadro []int) {
	for i := 0; i < len(quadro); i++ {
		fmt.Print(quadro[i]," ")
		
		if (i+1) % 9 == 0 {
			fmt.Println()
		}
	}
}

func atualizarPossib(quadro []int, px int, py int, possibs *Possib) {
	// Se já foi escolhido no quadro, só tem aquela opção disponível
	if quadro[py*9+px] != 0 {
		for i := 0; i < 9; i++ {
			possibs.p[i] = false
		}
		possibs.p[quadro[py*9+px]-1] = true
		return
	}
	// verificar colunas
	for y := 0; y < 9; y++ {
		if y == py { continue }
		quadro_v := quadro[y*9 + px]
		if quadro_v == 0 {continue}

		possibs.p[quadro_v-1] = false
	}

	// verificar linhas
	for x := 0; x < 9; x++ {
		if x == px { continue }
		quadro_v := quadro[py*9 + x]
		if quadro_v == 0 {continue}

		possibs.p[quadro_v-1] = false
	}

	// verificar quadrado
	quadx := (px / 3) * 3
	quady := (py / 3) * 3
	for x := quadx; x < quadx+3; x++ {
		for y := quady; y < quady+3; y++ {
			if x == px || y == py { continue }
			quadro_v := quadro[y*9 + x]
			if quadro_v == 0 {continue}

			possibs.p[quadro_v-1] = false
		}
	}
}

// Mantém apenas as possibilidades válidas
func atualizarQuadroPossib(quadro []int, quadro_possib []Possib) {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			possibs := &quadro_possib[y*9+x]
		
			atualizarPossib(quadro,x,y,possibs)
		}
	}
}

func colapsar(quadro []int,quadro_possib []Possib,index int) error {
	valor := -1
	possibs := &quadro_possib[index]
	for i := 0; i < len(possibs.p); i++ {
		if possibs.p[i] { 
			valor = i
		}
	}
	if valor == -1 {
		return fmt.Errorf("Não tem como colapsar, 0 possibilidades")
	}

	quadro[index] = valor+1
	return nil
}

func checarSolucionado(quadro []int) bool {
	for i := 0; i < 9*9; i++ {
		if quadro[i] == 0 { return false }
	}
	return true
}

func solucionarSudoku(quadro []int, quadro_possib []Possib) bool {
	for iter := 0; iter < 1000000; iter++{
		if checarSolucionado(quadro) {
			fmt.Println("Iter:",iter)
			return true
		}

		atualizarQuadroPossib(quadro,quadro_possib)

		min_index := -1
		min_cont := -1
		for i := 0; i < len(quadro); i++ {
			possibs := &quadro_possib[i]
			if quadro[i] == 0 && (min_cont == -1 || possibs.contar() < min_cont) {
				min_cont = possibs.contar()
				min_index = i
			}
		}

		err := colapsar(quadro,quadro_possib,min_index)
		if err != nil {
			fmt.Println(err)
			printarPossib(quadro_possib)
			return false
		}
		//printarPossib(quadro_possib)
		//fmt.Println("\n--------------------")
	}

	if checarSolucionado(quadro) {
		return true
	} else {
		return false
	}
}

func Exec() {
	fmt.Println("Sudoku")

	var quadro = []int{
		5, 0, 0,  9, 0, 0,  7, 2, 4,
		3, 0, 0,  0, 0, 7,  5, 0, 0,
		0, 0, 9,  0, 5, 4,  0, 0, 6,
	
		0, 0, 0,  0, 4, 5,  0, 0, 8,
		8, 0, 5,  0, 0, 6,  0, 4, 3,
		0, 0, 7,  8, 0, 0,  0, 0, 2,
	
		6, 5, 8,  4, 0, 0,  0, 3, 1,
		0, 0, 0,  0, 6, 0,  8, 9, 0,
		2, 9, 0,  0, 3, 8,  0, 6, 7,
	}
	// 9 possibilidades em cada quadrado
	var quadro_possib = iniciarPossib(quadro,9)


	if solucionarSudoku(quadro,quadro_possib) {
		fmt.Println("Solucionado!")
		printarQuadro(quadro)
	} else {
		fmt.Println("Não conseguiu solucionar")
		printarQuadro(quadro)
	}
}