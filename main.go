package main

import (
	"fmt"

	"github.com/erickweil/horariogen/horario"
	//"github.com/erickweil/horariogen/meucanvas"
	//"github.com/erickweil/horariogen/genetic"
)

// https://github.com/golang-standards/project-layout/tree/master
func main() {
	fmt.Println("OK")

	horario.ExecHorario()
	//sudoku.ExecSudoku()

	//meucanvas.ExecCanvas()
	//genetic.ExecHillClimbing()
}