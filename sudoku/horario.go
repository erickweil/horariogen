package sudoku

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

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
		hDia := h.getHorarioDia(dia)
		temposExpandido := make([]int, nTempos)
		for i :=0; i < len(hDia); i++ { 
			temposExpandido[hDia[i]-1] = 1
		}
		h.setHorarioDia(dia,temposExpandido)
	}
}

func (h *Horario) possui(dia int, tempo int) bool {
	return  h.getHorarioDia(dia)[tempo] == 1
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
	Nome string			`json:"nome"`
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

	return nil
}

func printarHorario(quadro []int) {
	for turma := 0; turma < len(turmas); turma++ {
		for tempo := 0; tempo < nTempos; tempo++ {
			for dia := 0; dia < 7; dia++ {
				idDisciplina := quadro[toQuadroIndex(turma,dia,tempo)]
				if idDisciplina > 0 {
					fmt.Print(disciplinas[idDisciplina-1].Nome[0:8],"\t")
				} else {
					if idDisciplina == 0 {
						fmt.Print("????????","\t")
					} else {
						fmt.Print("--------","\t")
					}
				}
				
				fmt.Print("\t\t")
			}
			fmt.Println()
		}
		fmt.Println("")
	}
}

func contarAulas(quadro []int,materia int) int {
	count := 0
	for i := 0; i < len(quadro); i++ {
		if quadro[i] == materia {
			count++
		}
	}
	return count
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

func regrasHorario(quadro []int, possibs *Possib) {
	index := possibs.index
	idTurma,dia,tempo := fromQuadroIndex(index)

	// Se já foi escolhido no quadro, só tem aquela opção disponível
	// -1 indica que não haverá nenhuma escolha,
	// 0 indica que não foi escolhido
	// >0 indica que uma matéria foi escolhida
	if quadro[index] != 0 {
		possibs.resetar(false)
		if quadro[index] >= 1 {
			possibs.p[quadro[index]-1] = true
		}
		return
	}

	// A turma deve poder ter aula neste tempo
	turma := turmas[idTurma]
	if !turma.Horarios.possui(dia,tempo) {
		possibs.resetar(false)
		return
	}

	// verificar se já tem o número necessário de aulas nessa matéria
	for i := 0; i < len(possibs.p); i++ {
		if !possibs.p[i] { continue }

		// esgotou o número de aulas a serem escolhidas desta matéria
		if contarAulas(quadro,i+1) >= disciplinas[i].Aulas {
			possibs.p[i] = false
		}
	}

	// As aulas devem ser de 2 em 2, sem cruzar o intervalo
	for i := 0; i < len(possibs.p); i++ {
		_index := toQuadroIndex(idTurma,dia,tempo+1)
		if (tempo == 0 || tempo == 2) && (quadro[_index] > 0 && quadro[_index]-1 != i) {
			possibs.p[i] = false
		}
		_index = toQuadroIndex(idTurma,dia,tempo-1)
		if (tempo == 1 || tempo == 3) && (quadro[_index] > 0 && quadro[_index]-1 != i) {
			possibs.p[i] = false
		}
	}

	// O professor deve estar disponivel neste tempo
	for i := 0; i < len(possibs.p); i++ {
		var disciplina *Disciplina = &disciplinas[i]
		for k := 0; k < len(disciplina.idProfessores); k++ {
			var professor *Professor = &professores[disciplina.idProfessores[k]]
			// O professor não dá aulas neste dia/tempo
			if !professor.Horarios.possui(dia,tempo) {
				possibs.p[i] = false
				break
			}

			// O professor já está em outra turma neste dia/tempo
			for cturma := 0; cturma < len(turmas); cturma++ {
				if cturma == idTurma { continue }

				caula := quadro[toQuadroIndex(cturma,dia,tempo)]
				if caula > 0 && disciplinas[caula-1].possuiProfessor(professor.id) {
					possibs.p[i] = false
					break
				}
			}
		}
	}

	// A disciplina deve poder ser nesta turma
	for i := 0; i < len(possibs.p); i++ {
		var disciplina *Disciplina = &disciplinas[i]

		if disciplina.idTurma != idTurma {
			possibs.p[i] = false
		}
	}
}

func ExecHorario() {
	fmt.Println("Horario")

	err := loadJson("./aulas.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer timeTrack(time.Now(),"Horario")

	if solucionarQuadro(quadro,len(disciplinas),regrasHorario) {
		fmt.Println("Solucionado! iter:",iter)
		printarHorario(quadro)
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarHorario(quadro)
	}

	/*solucoes := solucionarQuadroSemParar(quadro,len(materias),regrasHorario,nil)
	if solucoes != nil {
		fmt.Println("Terminou de procurar soluções! iter:",iter," nSolucoes:",len(solucoes))
	} else {
		fmt.Println("Não conseguiu solucionar iter:",iter)
		printarHorario(quadro)
	}*/
}