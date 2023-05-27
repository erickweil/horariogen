package sudoku

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/erickweil/horariogen/pencilmark"
	"github.com/erickweil/horariogen/utils"
)

func printarQuadro(quadro []int) {
	for i := 0; i < len(quadro); i++ {
		fmt.Print(quadro[i]," ")
		
		if (i+1) % 3 == 0 {
			fmt.Print(" ")
		}

		if (i+1) % 9 == 0 {
			fmt.Println()
		}

		if (i+1) % (9*3) == 0 {
			fmt.Println()
		}
	}
}

func regrasSudoku(quadro []int, possibs *pencilmark.Possib) {
	if possibs == nil {
		return
	}

	index := possibs.Index
	py := index / 9
	px := index % 9
	// Se já foi escolhido no quadro, só tem aquela opção disponível
	if quadro[index] != 0 {
		for i := 0; i < 9; i++ {
			possibs.P[i] = false
		}
		possibs.P[quadro[index]-1] = true
		return
	}
	// verificar colunas
	for y := 0; y < 9; y++ {
		if y == py { continue }
		quadro_v := quadro[y*9 + px]
		if quadro_v == 0 {continue}

		possibs.P[quadro_v-1] = false
	}

	// verificar linhas
	for x := 0; x < 9; x++ {
		if x == px { continue }
		quadro_v := quadro[py*9 + x]
		if quadro_v == 0 {continue}

		possibs.P[quadro_v-1] = false
	}

	// verificar quadrado
	quadx := (px / 3) * 3
	quady := (py / 3) * 3
	for x := quadx; x < quadx+3; x++ {
		for y := quady; y < quady+3; y++ {
			if x == px || y == py { continue }
			quadro_v := quadro[y*9 + x]
			if quadro_v == 0 {continue}

			possibs.P[quadro_v-1] = false
		}
	}
}

func ExecSudoku() {
	
	rand.Seed(time.Now().UnixNano())
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

		0, 0, 0,  0, 0, 0,  0, 0, 0,
		0, 0, 0,  0, 0, 3,  0, 8, 5,
		0, 0, 1,  0, 2, 0,  0, 0, 0,
		
		0, 0, 0,  5, 0, 7,  0, 0, 0, 
		0, 0, 4,  0, 0, 0,  1, 0, 0,
		0, 9, 0,  0, 0, 0,  0, 0, 0,

		5, 0, 0,  0, 0, 0,  0, 7, 3,
		0, 0, 2,  0, 1, 0,  0, 0, 0, 
		0, 0, 0,  0, 4, 0,  0, 0, 9,

		/*0, 0, 0,  9, 0, 0,  7, 2, 4,
		3, 0, 0,  0, 0, 7,  5, 0, 0,
		0, 0, 9,  0, 0, 4,  0, 0, 6,
	
		0, 0, 0,  0, 4, 5,  0, 0, 8,
		8, 0, 5,  0, 0, 6,  0, 4, 3,
		0, 0, 7,  8, 0, 0,  0, 0, 2,
	
		6, 5, 8,  4, 0, 0,  0, 3, 1,
		0, 0, 0,  0, 6, 0,  8, 9, 0,
		2, 9, 0,  0, 3, 8,  0, 0, 0,*/

		/*8, 0, 0,  0, 0, 0,  0, 0, 0,
		0, 0, 3,  6, 0, 0,  0, 0, 0,
		0, 7, 0,  0, 9, 0,  2, 0, 0,
	
		0, 5, 0,  0, 0, 7,  0, 0, 0,
		0, 0, 0,  0, 4, 5,  7, 0, 0,
		0, 0, 0,  1, 0, 0,  0, 3, 0,
	
		0, 0, 1,  0, 0, 0,  0, 6, 8,
		0, 0, 8,  5, 0, 0,  0, 1, 0,
		0, 9, 0,  0, 0, 0,  4, 0, 0,*/
	}

	defer utils.TimeTrack(time.Now(),"Sudoku")

	iter, solved := pencilmark.SolucionarQuadro(quadro,9,regrasSudoku)
	if solved {
		fmt.Println("Solucionado! iter:",iter)
		printarQuadro(quadro)
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarQuadro(quadro)
	}

	/*solucoes := solucionarQuadroSemParar(quadro,9,regrasSudoku,nil)
	if solucoes != nil {
		fmt.Println("Terminou de procurar soluções! iter:",iter," nSolucoes:",len(solucoes))
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarQuadro(quadro)
	}*/
}