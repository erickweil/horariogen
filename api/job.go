package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/erickweil/horariogen/horario"
)

//================================================================//
// 1. MODELOS E TIPOS DE DADOS (Models and Data Types)            //
//================================================================//

// JobStatus define o estado de um trabalho de processamento.
type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

// Job representa uma tarefa de geração de horário.
type Job struct {
	ID          string          `json:"id"`
	Status      JobStatus       `json:"status"`
	Config      horario.ArquivoJson `json:"-"` // Oculta a configuração completa na resposta padrão
	Result      []map[string]interface{} `json:"result,omitempty"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	CompletedAt *time.Time      `json:"completedAt,omitempty"`
}

//================================================================//
// 2. ARMAZENAMENTO EM MEMÓRIA (In-Memory Job Store)              //
//================================================================//

// JobStore é um mock para armazenamento de jobs em memória.
// Ele é seguro para uso concorrente (thread-safe).
type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

// NewJobStore cria uma nova instância do nosso armazenamento de jobs.
func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*Job),
	}
}

// AddJob adiciona um novo job ao armazenamento.
func (s *JobStore) AddJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

// GetJob recupera um job pelo seu ID.
func (s *JobStore) GetJob(id string) (*Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, found := s.jobs[id]
	return job, found
}

// UpdateJob atualiza o estado de um job existente.
func (s *JobStore) UpdateJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

//================================================================//
// 3. LÓGICA DE PROCESSAMENTO (Business Logic / Solver)           //
//================================================================//

// processHorarioJob é a função que executa o seu código de forma assíncrona.
// Ela atualiza o status do job conforme progride.
func (s *JobStore) processHorarioJob(job *Job) {
	// Atualiza o status para "running"
	job.Status = StatusRunning
	s.UpdateJob(job)

	fmt.Printf("Iniciando processamento para o Job ID: %s\n", job.ID)

	resultado, err := horario.ExecHorario(&job.Config)
	// --------------------------------

	completedTime := time.Now()
	job.CompletedAt = &completedTime

	if err != nil {
		fmt.Printf("Job ID: %s falhou: %v\n", job.ID, err)
		job.Status = StatusFailed
		job.Error = err.Error()
	} else {
		fmt.Printf("Job ID: %s concluído com sucesso!\n", job.ID)
		job.Status = StatusCompleted
		job.Result = resultado
	}

	// Atualiza o job no armazenamento com o resultado final.
	s.UpdateJob(job)
}