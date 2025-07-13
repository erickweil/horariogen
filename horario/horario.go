package horario

import (
	"encoding/json"
	"fmt"
	"io"
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
	idDisciplinasUnidas []int	`json:"-"`
	Turma string		`json:"turma"`
	Aulas int			`json:"aulas"`
	Agrupar int			`json:"agrupar"`
	Nome string			`json:"nome"`

	// cache
	contAulas int	`json:"-"`
}

type DisciplinaUnida struct {
	id int				`json:"-"`
	Grupo string		`json:"grupo"`
	Disciplinas []string	`json:"disciplinas"`
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
	DisciplinaUnidas []DisciplinaUnida	`json:"disciplinas_unidas"`
	Professores []Professor		`json:"professores"`
}

// HorarioProcessor encapsula todos os dados e a lógica para uma única execução.
type HorarioProcessor struct {
	professores []Professor
	turmas      []Turma
	disciplinas []Disciplina
	nTempos     int
	quadro      []int
}

/*
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
}*/

func loadBytes(caminho string) ([]byte, error) {
	content, err := os.Open(caminho)
    if err != nil {
        return nil, fmt.Errorf("erro ao abrir arquivo %w",err)
    }
	defer content.Close()

	bytes, err := io.ReadAll(content)
    if err != nil {
        return nil, fmt.Errorf("erro ao ler arquivo %w",err)
    }
	return bytes, nil
}

//func loadJson(arquivoJson ArquivoJson) error {
func NewHorarioProcessor(arquivoJson ArquivoJson) (*HorarioProcessor, error) {
	/*turmas = arquivoJson.Turmas
	disciplinas = arquivoJson.Disciplinas
	professores = arquivoJson.Professores*/
	disciplinas_unidas := arquivoJson.DisciplinaUnidas
	p := &HorarioProcessor{
		turmas:      arquivoJson.Turmas,
		disciplinas: arquivoJson.Disciplinas,
		professores: arquivoJson.Professores,
	}

	for i := 0; i < len(p.turmas); i++ {
		maxTempos := p.turmas[i].Horarios.getMaxTempo()
		if maxTempos > p.nTempos {
			p.nTempos = maxTempos
		}
	}

	for i := 0; i < len(p.turmas); i++ {
		p.turmas[i].id = i
		p.turmas[i].Horarios.expandirTempos(p.nTempos)
	}

	// Mapeia nomes para IDs e processa disciplinas unidas.
	disciplinasMap := make(map[string]int)
	for i := range p.disciplinas {
		disciplinasMap[p.disciplinas[i].Nome] = i
	}

	turmasMap := make(map[string]int)
	for i := range p.turmas {
		turmasMap[p.turmas[i].Nome] = i
	}

	for i := 0; i < len(p.disciplinas); i++ {
		var d *Disciplina = &p.disciplinas[i]
		d.id = i
		//d.idTurma = getTurma(d.Turma)
		if id, ok := turmasMap[d.Turma]; ok {
			d.idTurma = id
		} else {
			return nil, fmt.Errorf("turma '%s' para a disciplina '%s' não encontrada", d.Turma, d.Nome)
		}

		var possuiUnidas bool = false
		for k := 0; k < len(disciplinas_unidas); k++ {
			var du *DisciplinaUnida = &disciplinas_unidas[k]
			for j := 0; j < len(du.Disciplinas); j++ {
				if du.Disciplinas[j] == d.Nome {
					possuiUnidas = true
					break
				}
			}
			if possuiUnidas {
				d.idDisciplinasUnidas = make([]int, 0)
				for j := 0; j < len(du.Disciplinas); j++ {
					if du.Disciplinas[j] == d.Nome { continue }
					//d.idDisciplinasUnidas = append(d.idDisciplinasUnidas,getDisciplina(du.Disciplinas[j]))
					if id, ok := disciplinasMap[du.Disciplinas[j]]; ok {
						d.idDisciplinasUnidas = append(d.idDisciplinasUnidas, id)
					} else {
						return nil, fmt.Errorf("disciplina unida '%s' não encontrada para a disciplina '%s'", du.Disciplinas[j], d.Nome)
					}
				}
			}
		}

		// printar disciplinas unidas
		/*if possuiUnidas {
			fmt.Print(d.Nome," Unida com: ")
			for j := 0; j < len(d.idDisciplinasUnidas); j++ {
				fmt.Print(disciplinas[d.idDisciplinasUnidas[j]].Nome," ")
			}
			fmt.Println()
		}*/
	}

	for i := 0; i < len(p.professores); i++ {
		var prof *Professor = &p.professores[i]
		prof.id = i
		for k := 0; k < len(prof.Disciplinas); k++ {
			//var d *Disciplina = &p.disciplinas[getDisciplina(p.Disciplinas[k])]
			//d.idProfessores = append(d.idProfessores, p.id)
			if id, ok := disciplinasMap[prof.Disciplinas[k]]; ok {
				p.disciplinas[id].idProfessores = append(p.disciplinas[id].idProfessores, prof.id)
			} else {
				return nil, fmt.Errorf("disciplina '%s' não encontrada para o professor '%s'", prof.Disciplinas[k], prof.Nome)
			}
		}
		prof.Horarios.expandirTempos(p.nTempos)
		prof.matriz.expandirTempos(p.nTempos)
	}

	// Preencher quadro com 0 onde pode ter aulas e -1 onde não haverá aulas
	p.quadro = make([]int, len(p.turmas) * 7 * p.nTempos)
	for i := 0; i < len(p.quadro); i++ { p.quadro[i] = -1 }

	for i := 0; i < len(p.turmas); i++ {
		for dia := 0; dia < 7; dia++ {
			for tempo := 0; tempo < p.nTempos; tempo++ {
				if p.turmas[i].Horarios.possui(dia,tempo) {
					p.quadro[p.toQuadroIndex(i,dia,tempo)] = 0
				}
			}
		}
	}

	fmt.Printf("Carregou %d turmas, %d disciplinas, %d professores, Horario: %d Tempos, %d espaços no quadro\n",
	len(p.turmas),len(p.disciplinas),len(p.professores),p.nTempos,len(p.quadro))

	for i := 0; i < len(p.turmas); i++ {
		// Somar a quantidade de aulas de cada turma, para validar caso aconteça aulas demais para o horário
		somaAulas := 0
		for j := 0; j < len(p.disciplinas); j++ {
			if p.disciplinas[j].idTurma == i {
				somaAulas += p.disciplinas[j].Aulas
			}
		}
		fmt.Printf("%s: %d aulas\n",p.turmas[i].Nome,somaAulas)
	}

	return p, nil
}

func trimToSize(str string, size int) string {
	if len(str) > size {
		return str[0:size]
	}
	return str
}

func (p *HorarioProcessor) printarHorario() {
	for turma := 0; turma < len(p.turmas); turma++ {
		for tempo := 0; tempo < p.nTempos; tempo++ {
			for dia := 0; dia < 7; dia++ {
				idDisciplina := p.quadro[p.toQuadroIndex(turma,dia,tempo)]
				if idDisciplina > 0 {
					//fmt.Print(disciplinas[idDisciplina-1].Nome[0:8],"\t")
					fmt.Printf("%-32s\t",trimToSize(p.disciplinas[idDisciplina-1].Nome,32))
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

func (p *HorarioProcessor) gerarHorarioJson() ([]map[string]interface{}, error) {
	horarioJson := make([]map[string]interface{}, 0)
	for turma := 0; turma < len(p.turmas); turma++ {
		turmaNome := p.turmas[turma].Nome
		turmaJson := make([]map[string]interface{}, 0)
		for dia := 0; dia < 7; dia++ {
			diaStr := ""
			switch dia {
			case 0: diaStr = "dom"
			case 1: diaStr = "seg"
			case 2: diaStr = "ter"
			case 3: diaStr = "qua"
			case 4: diaStr = "qui"
			case 5: diaStr = "sex"
			case 6: diaStr = "sab"
			}
			diaJson := make([]interface{}, 0)
			for tempo := 0; tempo < p.nTempos; tempo++ {
				idDisciplina := p.quadro[p.toQuadroIndex(turma,dia,tempo)]
				if idDisciplina > 0 {
					disciplina := p.disciplinas[idDisciplina-1]
					
					diaJson = append(diaJson, disciplina.Nome)
				} else if idDisciplina == 0 {
					diaJson = append(diaJson, "????????") // Não resolvido
				} else {
					diaJson = append(diaJson, nil) // Não há aula nesse tempo mesmo
				}
			}
			//turmaJson[diaStr] = diaJson
			turmaJson = append(turmaJson, map[string]interface{}{
				"dia": diaStr,
				"tempos": diaJson,
			})
		}
		//horarioJson[turmaNome] = turmaJson
		horarioJson = append(horarioJson, map[string]interface{}{
			"turma": turmaNome,
			"horario": turmaJson,
		})
	}
	
	return horarioJson, nil
}

func (p *HorarioProcessor) toQuadroIndex(turma int,dia int,tempo int) int {
	nDias := 7
	return turma * (nDias * p.nTempos) + dia * p.nTempos + tempo 
}

func (p *HorarioProcessor) fromQuadroIndex(index int) (int,int,int) {
	nDias := 7
	turma := (index / p.nTempos) / nDias
	dia := (index / p.nTempos) % nDias
	tempo := index % p.nTempos
	return turma, dia, tempo
}

func (p *HorarioProcessor) getQuadroValor(turma int,dia int,tempo int) int {
	if turma < 0 || turma >= len(p.turmas) { return -1 }
	if dia < 0 || dia >= 7 { return -1 }
	if tempo < 0 || tempo >= p.nTempos { return -1 }
	return p.quadro[p.toQuadroIndex(turma,dia,tempo)]
}

func (p *HorarioProcessor) iniciarRegras(quadro []int) {
	for prof := 0; prof < len(p.professores); prof++ {
		p.professores[prof].matriz.copiar(&p.professores[prof].Horarios)
	}

	for disc := 0; disc < len(p.disciplinas); disc++ {
		disciplina := &p.disciplinas[disc]
		disciplina.contAulas = 0
	}

	for idTurma := 0; idTurma < len(p.turmas); idTurma++ {
		//turma := &turmas[idTurma]
		for tempo := 0; tempo < p.nTempos; tempo++ {
			for dia := 0; dia < 7; dia++ {
				index := p.toQuadroIndex(idTurma,dia,tempo)
				idDisciplina := quadro[index]-1
				if idDisciplina < 0 { continue }
				
				disciplina := &p.disciplinas[idDisciplina]
				disciplina.contAulas++

				for _, idProf := range disciplina.idProfessores {
					p.professores[idProf].matriz.getHorarioDia(dia)[tempo] = 0
				}
			}
		}
	}
}

func (p *HorarioProcessor) regrasHorario(quadro []int, possibs *Possib) {
	if possibs == nil { // atualizar cache com base no quadro
		p.iniciarRegras(quadro)
		return
	}

	index := possibs.Index
	idTurma,dia,tempo := p.fromQuadroIndex(index)

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
	turma := &p.turmas[idTurma]
	if !turma.Horarios.possui(dia,tempo) {
		possibs.Resetar(false)
		return
	}

	// verificar se já tem o número necessário de aulas nessa matéria
	for i := 0; i < len(possibs.P); i++ {
		if !possibs.P[i] { continue }
		
		possibs.P[i] = p.podeDisciplina(i,idTurma,dia,tempo)
	}
}

func (p *HorarioProcessor) podeDisciplina(idDisciplina int,idTurma int ,dia int, tempo int) bool {
	var disciplina *Disciplina = &p.disciplinas[idDisciplina]

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
		var professor *Professor = &p.professores[disciplina.idProfessores[k]]
		// O professor não dá aulas neste dia/tempo ou O professor já está em outra turma neste dia/tempo
		if !professor.matriz.possui(dia,tempo) {
			return false
		}
	}

	// As aulas devem ser agrupadas de acordo com o especificado
	// Regra só funciona no horário atual, depois tem que ver isso
	if disciplina.Agrupar == 2 {
		if tempo == 0 || tempo == 2 {
			// Verifica se o tempo depois desse não é da mesma disciplina
			_index := p.toQuadroIndex(idTurma,dia,tempo+1)
			if p.quadro[_index] > 0 && p.quadro[_index]-1 != idDisciplina {
				return false
			}
		}

		if tempo == 1 || tempo == 3 {
			// Verifica se o tempo antes desse não é da mesma disciplina
			_index := p.toQuadroIndex(idTurma,dia,tempo-1)
			if p.quadro[_index] > 0 && p.quadro[_index]-1 != idDisciplina {
				return false
			}
		}

		if tempo == 0 || tempo == 1 {
			// Impedir que fique 4 seguidas
			if p.getQuadroValor(idTurma,dia,3)-1 == idDisciplina &&
			   p.getQuadroValor(idTurma,dia,4)-1 == idDisciplina {
				return false
			}
		} else {
			// Impedir que fique 4 seguidas
			if p.getQuadroValor(idTurma,dia,0)-1 == idDisciplina &&
				p.getQuadroValor(idTurma,dia,1)-1 == idDisciplina {
				return false
			}	
		}
	}
	if disciplina.Agrupar == 4 && disciplina.contAulas > 0 {
		if p.getQuadroValor(idTurma,dia,0)-1 != idDisciplina &&
		p.getQuadroValor(idTurma,dia,1)-1 != idDisciplina &&
		p.getQuadroValor(idTurma,dia,2)-1 != idDisciplina &&
		p.getQuadroValor(idTurma,dia,3)-1 != idDisciplina {
			return false
		}
	}

	// As disciplinas unidas devem ser escolhidas juntas
	if len(disciplina.idDisciplinasUnidas) > 0 {
		// Verifica se alguma das disciplinas unidas já foram escolhidas
		// Deve poder escolher as outras disciplinas nesse horario também
		for i := 0; i < len(disciplina.idDisciplinasUnidas); i++ {
			var discUnida *Disciplina = &p.disciplinas[disciplina.idDisciplinasUnidas[i]]
			var marcacaoQuadro = p.getQuadroValor(discUnida.idTurma,dia,tempo)

			if marcacaoQuadro != 0 && marcacaoQuadro != discUnida.id+1 {
				// Se foi escolhido, deve ser a mesma disciplina
				return false
			}
		}
	}

	return true
}

// ExecHorario executa o solucionador para esta instância do processador.
func (p *HorarioProcessor) ExecHorario() ([]int, error) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Iniciando solucionador de horário...")
	defer utils.TimeTrack(time.Now(), "ExecHorario")

    // A função de regras agora é um método do nosso processador
	iter, solved := pencilmark.SolucionarQuadro(p.quadro, len(p.disciplinas), p.regrasHorario)

	if solved {
		fmt.Println("Solucionado! Iterações:", iter)
	} else {
		fmt.Println("Não conseguiu solucionar! Iterações:", iter)
	}
	
	p.printarHorario()
	return p.quadro, nil
}

func ExecHorario(arquivoJson *ArquivoJson) ([]map[string]interface{}, error) {
	fmt.Println("Horario")

	if (arquivoJson == nil) {
		bytes, err := loadBytes("./aulas.json")
		if err != nil {
			return nil, fmt.Errorf("erro ao ler json: %w", err)
		}

		arquivoJson = &ArquivoJson{}
		json.Unmarshal(bytes, arquivoJson)
	}
	
	p, err := NewHorarioProcessor(*arquivoJson)
	if err != nil {
		return nil, fmt.Errorf("erro ao construir horario processor: %w", err)
	}
	
	_, err = p.ExecHorario()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar horário: %w", err)
	}
	fmt.Println("Quadro finalizado!")

	return p.gerarHorarioJson()

	/*var iter int = 0
	solucoes := pencilmark.SolucionarQuadroSemParar(quadro,&iter,len(disciplinas),regrasHorario,nil)
	if solucoes != nil {
		for i:=0; i< len(solucoes) && i < 10; i++{
			fmt.Println("\nSolução:",i)
			printarHorario(solucoes[i])
		}
		fmt.Println("Terminou de procurar soluções! nSolucoes:",len(solucoes))
	} else {
		fmt.Println("Não conseguiu solucionar!")
		printarHorario(quadro)
	}*/
}