package main

import (
	//"encoding/json"
	"fmt"
	"github.com/erickweil/horariogen/api"
	//"github.com/erickweil/horariogen/horario"
	//"github.com/erickweil/horariogen/meucanvas"
	//"github.com/erickweil/horariogen/genetic"
)

// https://github.com/golang-standards/project-layout/tree/master
func main() {
	fmt.Println("OK")

	api.InitAPIServer();

	/*result, err := horario.ExecHorario(nil)
	if err != nil {
		fmt.Println("Erro ao executar horário:", err)
		return
	}

	// unmarshal map[string]interface{} to json pretty
	prettyJson, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Erro ao formatar JSON:", err)
		return
	}
	fmt.Println("Resultado do horário:")
	fmt.Println(string(prettyJson))*/


	//sudoku.ExecSudoku()

	//meucanvas.ExecCanvas()
	//genetic.ExecHillClimbing()
}