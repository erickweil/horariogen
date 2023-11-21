package horario

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"github.com/erickweil/horariogen/pencilmark"
	"github.com/erickweil/horariogen/utils"
)

type Possib = pencilmark.Possib

type Horario struct {
	Dom []int	`json:"dom"`
	Seg []int	`json:"seg"`
	Ter []int	`json:"ter"`
	Qua []int	`json:"qua"`
	Qui []int	`json:"qui"`
	Sex []int	`json:"sex"`
	Sab []int	`json:"sab"`
}

func (h *Horario) getHorarioDia(dia int) []int {
	if dia == 0 { return h.Dom }
	if dia == 1 { return h.Seg }
	if dia == 2 { return h.Ter }
	if dia == 3 { return h.Qua }
	if dia == 4 { return h.Qui }
	if dia == 5 { return h.Sex }
	if dia == 6 { return h.Sab }
	return nil
}

func (h *Horario) setHorarioDia(dia int, tempos []int) {
	if dia == 0 { h.Dom = tempos }
	if dia == 1 { h.Seg = tempos }
	if dia == 2 { h.Ter = tempos }
	if dia == 3 { h.Qua = tempos }
	if dia == 4 { h.Qui = tempos }
	if dia == 5 { h.Sex = tempos }
	if dia == 6 { h.Sab = tempos }
}

func (h *Horario) expandirTempos(nTempos int) {
	for dia := 0; dia < 7; dia++ {
		temposExpandido := make([]int, nTempos)
		hDia := h.getHorarioDia(dia)
		if hDia != nil {
			for i :=0; i < len(hDia); i++ { 
				temposExpandido[hDia[i]-1] = 1
			}
		}
		h.setHorarioDia(dia,temposExpandido)
	}
}

func (h *Horario) copiar(outro *Horario) {
	for dia := 0; dia < 7; dia++ {
		hDia := h.getHorarioDia(dia)
		oDia := outro.getHorarioDia(dia)
		for i :=0; i < len(hDia); i++ { 
			hDia[i] = oDia[i]
		}
	}
}

func (h *Horario) possui(dia int, tempo int) bool {
	return  h.getHorarioDia(dia)[tempo] != 0
}

func (h *Horario) getMaxTempo() int {
	ret := 0
	for dia := 0; dia < 7; dia++ {
		hDia := h.getHorarioDia(dia)
		for i :=0; i < len(hDia); i++ { if ret < hDia[i] { ret = hDia[i] } }
	}
	return ret
}

type Professor struct {
	id int					`json:"-"`
	Nome string				`json:"nome"`
	Disciplinas []string	`json:"disciplinas"`
	Horarios Horario		`json:"horarios"`

	// cache
	matriz Horario	`json:"-"` 
}

type Turma struct {
	id int				`json:"-"`
	Nome string			`json:"nome"`
	Horarios Horario	`json:"horarios"`
}

type Disciplina struct {
	id int				`json:"-"`
	idProfessores []int	`json:"-"`
	idTurma int			`json:"-"`
	Turma string		`json:"turma"`
	Aulas int			`json:"aulas"`
	Agrupar int			`json:"agrupar"`
	Nome string			`json:"nome"`

	// cache
	contAulas int	`json:"-"`
}

func (d *Disciplina) possuiProfessor(prof int) bool {
	for i := 0; i < len(d.idProfessores); i++ {
		if d.idProfessores[i] == prof { return true }
	}
	return false
}

type ArquivoJson struct {
	Turmas []Turma				`json:"turmas"`
	Disciplinas []Disciplina	`json:"disciplinas"`
	Professores []Professor		`json:"professores"`
}

var professores []Professor
var turmas []Turma
var disciplinas []Disciplina
var nTempos int
var quadro []int

func getProfessor(nome string) int {
	for i := 0; i < len(professores); i++ {
		if professores[i].Nome == nome {
			return i
		}
	}
	return -1
}

func getTurma(nome string) int {
	for i := 0; i < len(turmas); i++ {
		if turmas[i].Nome == nome {
			return i
		}
	}
	return -1
}

func getDisciplina(nome string) int {
	for i := 0; i < len(disciplinas); i++ {
		if disciplinas[i].Nome == nome {
			return i
		}
	}
	return -1
}

func loadJson(caminho string) error {
	content, err := os.Open(caminho)
    if err != nil {
        return fmt.Errorf("Erro ao abrir arquivo %w",err)
    }
	defer content.Close()

	bytes, err := ioutil.ReadAll(content)
    if err != nil {
        return fmt.Errorf("Erro ao ler arquivo %w",err)
    }

	var arquivoJson ArquivoJson

	json.Unmarshal(bytes,&arquivoJson)

	turmas = arquivoJson.Turmas
	disciplinas = arquivoJson.Disciplinas
	professores = arquivoJson.Professores

	nTempos = 0
	for i := 0; i < len(turmas); i++ {
		turmas[i].id = i
		maxTempos := turmas[i].Horarios.getMaxTempo()
		if maxTempos > nTempos {
			nTempos = maxTempos
		}
	}

	for i := 0; i < len(turmas); i++ {
		turmas[i].Horarios.expandirTempos(nTempos)
	}

	for i := 0; i < len(disciplinas); i++ {
		var d *Disciplina = &disciplinas[i]
		d.id = i
		d.idTurma = getTurma(d.Turma)
	}

	for i := 0; i < len(professores); i++ {
		var p *Professor = &professores[i]
		p.id = i
		for k := 0; k < len(p.Disciplinas); k++ {
			var d *Disciplina = &disciplinas[getDisciplina(p.Disciplinas[k])]
			d.idProfessores = append(d.idProfessores, p.id)
		}
		p.Horarios.expandirTempos(nTempos)
		p.matriz.expandirTempos(nTempos)
	}

	// Preencher quadro com 0 onde pode ter aulas e -1 onde não haverá aulas
	quadro = make([]int, len(turmas) * 7 * nTempos)
	for i := 0; i < len(quadro); i++ {	quadro[i] = -1 }

	for i := 0; i < len(turmas); i++ {
		for dia := 0; dia < 7; dia++ {
			for tempo := 0; tempo < nTempos; tempo++ {
				if turmas[i].Horarios.possui(dia,tempo) {
					quadro[toQuadroIndex(i,dia,tempo)] = 0
				}
			}
		}
	}

	fmt.Printf("Carregou %d turmas, %d disciplinas, %d professores, Horario: %d Tempos, %d espaços no quadro\n",
	len(turmas),len(disciplinas),len(professores),nTempos,len(quadro))

	for i := 0; i < len(turmas); i++ {
		// Somar a quantidade de aulas de cada turma, para validar caso aconteça aulas demais para o horário
		somaAulas := 0
		for j := 0; j < len(disciplinas); j++ {
			if disciplinas[j].idTurma == i {
				somaAulas += disciplinas[j].Aulas
			}
		}
		fmt.Printf("%s: %d aulas\n",turmas[i].Nome,somaAulas)
	}


	return nil
}

func printarHorario(quadro []int) {
	for turma := 0; turma < len(turmas); turma++ {
		for tempo := 0; tempo < nTempos; tempo++ {
			for dia := 0; dia < 7; dia++ {
				idDisciplina := quadro[toQuadroIndex(turma,dia,tempo)]
				if idDisciplina > 0 {
					//fmt.Print(disciplinas[idDisciplina-1].Nome[0:8],"\t")
					fmt.Print(disciplinas[idDisciplina-1].Nome,"\t")
				} else {
					if idDisciplina == 0 {
						fmt.Print("????????","\t")
					} else {
						fmt.Print("--------","\t")
					}
				}
				
				//fmt.Print("\t\t")
			}
			fmt.Println()
		}
		fmt.Println("")
	}
}

func toQuadroIndex(turma int,dia int,tempo int) int {
	nDias := 7
	return turma * (nDias * nTempos) + dia * nTempos + tempo 
}

func fromQuadroIndex(index int) (int,int,int) {
	nDias := 7
	turma := (index / nTempos) / nDias
	dia := (index / nTempos) % nDias
	tempo := index % nTempos
	return turma, dia, tempo
}

func getQuadroValor(turma int,dia int,tempo int) int {
	if turma < 0 || turma >= len(turmas) { return -1 }
	if dia < 0 || dia >= 7 { return -1 }
	if tempo < 0 || tempo >= nTempos { return -1 }
	return quadro[toQuadroIndex(turma,dia,tempo)]
}

func iniciarRegras(quadro []int) {
	for prof := 0; prof < len(professores); prof++ {
		professores[prof].matriz.copiar(&professores[prof].Horarios)
	}

	for disc := 0; disc < len(disciplinas); disc++ {
		disciplina := &disciplinas[disc]
		disciplina.contAulas = 0
	}

	for idTurma := 0; idTurma < len(turmas); idTurma++ {
		//turma := &turmas[idTurma]
		for tempo := 0; tempo < nTempos; tempo++ {
			for dia := 0; dia < 7; dia++ {
				index := toQuadroIndex(idTurma,dia,tempo)
				idDisciplina := quadro[index]-1
				if idDisciplina < 0 { continue }
				
				disciplina := &disciplinas[idDisciplina]
				disciplina.contAulas++

				for _, idProf := range disciplina.idProfessores {
					professores[idProf].matriz.getHorarioDia(dia)[tempo] = 0
				}
			}
		}
	}
}

func regrasHorario(quadro []int, possibs *Possib) {
	if possibs == nil { // atualizar cache com base no quadro
		iniciarRegras(quadro)
		return
	}

	index := possibs.Index
	idTurma,dia,tempo := fromQuadroIndex(index)

	// Se já foi escolhido no quadro, só tem aquela opção disponível
	// -1 indica que não haverá nenhuma escolha,
	// 0 indica que não foi escolhido
	// >0 indica que uma matéria foi escolhida
	if quadro[index] != 0 {
		possibs.Resetar(false)
		if quadro[index] > 0 {
			possibs.P[quadro[index]-1] = true
		}
		return
	}

	// A turma deve poder ter aula neste tempo
	turma := &turmas[idTurma]
	if !turma.Horarios.possui(dia,tempo) {
		possibs.Resetar(false)
		return
	}

	// verificar se já tem o número necessário de aulas nessa matéria
	for i := 0; i < len(possibs.P); i++ {
		if !possibs.P[i] { continue }
		
		possibs.P[i] = podeDisciplina(i,idTurma,dia,tempo)
	}
}

func podeDisciplina(idDisciplina int,idTurma int ,dia int, tempo int) bool {
	var disciplina *Disciplina = &disciplinas[idDisciplina]

	// A disciplina deve poder ser nesta turma
	if disciplina.idTurma != idTurma {
		return false
	}

	// esgotou o número de aulas a serem escolhidas desta matéria
	if disciplina.contAulas >= disciplina.Aulas {
		return false
	}

	// O professor deve estar disponivel neste tempo
	for k := 0; k < len(disciplina.idProfessores); k++ {
		var professor *Professor = &professores[disciplina.idProfessores[k]]
		// O professor não dá aulas neste dia/tempo ou O professor já está em outra turma neste dia/tempo
		if !professor.matriz.possui(dia,tempo) {
			return false
		}
	}

	// As aulas devem ser agrupadas de acordo com o especificado
	// Regra só funciona no horário atual, depois tem que ver isso
	if disciplina.Agrupar == 2 {
		_index := toQuadroIndex(idTurma,dia,tempo+1)
		if (tempo == 0 || tempo == 2) && (quadro[_index] > 0 && quadro[_index]-1 != idDisciplina) {
			return false
		}
		_index = toQuadroIndex(idTurma,dia,tempo-1)
		if (tempo == 1 || tempo == 3) && (quadro[_index] > 0 && quadro[_index]-1 != idDisciplina) {
			return false
		}
	}
	if disciplina.Agrupar == 4 && disciplina.contAulas > 0 {
		if getQuadroValor(idTurma,dia,0)-1 != idDisciplina &&
		getQuadroValor(idTurma,dia,1)-1 != idDisciplina &&
		getQuadroValor(idTurma,dia,2)-1 != idDisciplina &&
		getQuadroValor(idTurma,dia,3)-1 != idDisciplina {
			return false
		}
	}

	return true
}

func ExecHorario() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Horario")

	err := loadJson("./aulas.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer utils.TimeTrack(time.Now(),"Horario")

	iter, solved := pencilmark.SolucionarQuadro(quadro,len(disciplinas),regrasHorario)
	if solved {
		fmt.Println("Solucionado! iter:",iter)
		printarHorario(quadro)
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarHorario(quadro)
	}

	/*solucoes := solucionarQuadroSemParar(quadro,len(disciplinas),regrasHorario,nil)
	if solucoes != nil {
		for i:=0; i< len(solucoes) && i < 10; i++{
			fmt.Println("\nSolução:",i)
			printarHorario(solucoes[i])
		}
		fmt.Println("Terminou de procurar soluções! iter:",iter," nSolucoes:",len(solucoes))
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarHorario(quadro)
	}*/
}