package sudoku

import (
	"fmt"
	"time"
)

//https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("%s took %s", name, elapsed)
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

func obterMelhorPossib(quadro []int) (*Possib, error) {
	var p *Possib = &Possib{-1,make([]bool, 9)}
	var min_p *Possib = &Possib{-1,make([]bool, 9)}
	var min_cont int = -1

	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			index := y*9+x
			if quadro[index] != 0 { continue } // se já foi escolhido ignora
			p.resetar()
			p.index = index
			
			atualizarPossib(quadro,x,y,p)

			cont := p.contar()
			if min_cont == -1 || cont < min_cont {
				min_cont = cont
				min_p.receber(p) // copia os valores
			}
		}
	}

	if min_cont == -1 {
		return nil,fmt.Errorf("Não encontrou nenhuma possibilidade")
	}

	return min_p,nil
}

func checarSolucionado(quadro []int) bool {
	for i := 0; i < 9*9; i++ {
		if quadro[i] == 0 { return false }
	}
	return true
}

var iter int = 0
func solucionarSudoku(quadro []int) bool {
	iter++

	// A ideia é obter o quadrado com menor entropia
	// isto é, que possui a menor quantidade de escolhas possíveis
	/*var quadro_possib = iniciarPossib(quadro,9)
	atualizarQuadroPossib(quadro,quadro_possib)
	min_index := -1
	min_cont := -1
	for i := 0; i < len(quadro_possib); i++ {
		possibs := &quadro_possib[i]
		if quadro[i] == 0 && (min_cont == -1 || possibs.contar() < min_cont) {
			min_cont = possibs.contar()
			min_index = i
		}
	}

	if min_cont <= 0 {
		//fmt.Println("Quadro impossível...")
		return false
	}

	p := &quadro_possib[min_index]*/

	p, err := obterMelhorPossib(quadro)

	if err != nil {
		//fmt.Println("Quadro impossível...")
		return false		
	}

	// Uma vez escolhido o quadrado a partir do qual continuar,
	// testa cada possibilidade deste quadrado
	for k := 0; k < len(p.p); k++ {
	//for k := len(p.p)-1; k >= 0; k-- {
		// Se é possível colocar o valor k neste quadrado
		if p.p[k] { 
			quadro[p.index] = k+1
			// Se com essa escolha já solucionou, retorna true
			if checarSolucionado(quadro) {
				return true
			}
			// Tenta solucionar com mais escolhas depois dessa, e se der certo retorna true
			if solucionarSudoku(quadro) {
				return true
			}

			// Essa escolha não resolveu o sudoku, remove ela para tentar outras
			quadro[p.index] = 0
		}
	}
	
	// Nenhuma escolha foi válida para este quadrado, ou seja, é ímpossível solucionar nesta configuração
	//fmt.Println("Backtracking...")
	return false
}

func Exec() {
	fmt.Println("Sudoku")

	var quadro = []int{
		/*5, 0, 0,  9, 0, 0,  7, 2, 4,
		3, 0, 0,  0, 0, 7,  5, 0, 0,
		0, 0, 9,  0, 5, 4,  0, 0, 6,
	
		0, 0, 0,  0, 4, 5,  0, 0, 8,
		8, 0, 5,  0, 0, 6,  0, 4, 3,
		0, 0, 7,  8, 0, 0,  0, 0, 2,
	
		6, 5, 8,  4, 0, 0,  0, 3, 1,
		0, 0, 0,  0, 6, 0,  8, 9, 0,
		2, 9, 0,  0, 3, 8,  0, 6, 7,*/

		8, 0, 0,  0, 0, 0,  0, 0, 0,
		0, 0, 3,  6, 0, 0,  0, 0, 0,
		0, 7, 0,  0, 9, 0,  2, 0, 0,
	
		0, 5, 0,  0, 0, 7,  0, 0, 0,
		0, 0, 0,  0, 4, 5,  7, 0, 0,
		0, 0, 0,  1, 0, 0,  0, 3, 0,
	
		0, 0, 1,  0, 0, 0,  0, 6, 8,
		0, 0, 8,  5, 0, 0,  0, 1, 0,
		0, 9, 0,  0, 0, 0,  4, 0, 0,
	}

	defer timeTrack(time.Now(),"Sudoku")

	if solucionarSudoku(quadro) {
		fmt.Println("Solucionado! iter:",iter)
		printarQuadro(quadro)
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarQuadro(quadro)
	}
}