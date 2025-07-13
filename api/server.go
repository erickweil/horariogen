package api

import (
	// "context"
	//"encoding/json"
	"fmt"
	//"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/erickweil/horariogen/horario"
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
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

//================================================================//
// 4. HANDLERS DA API (API Handlers)                              //
//================================================================//

// Handlers contém dependências para os handlers da API, como o JobStore.
type Handlers struct {
	Store *JobStore
}

// StartJobHandler lida com a requisição para iniciar um novo job.
// @Summary Inicia um novo processamento de horário
// @Description Recebe uma configuração JSON e inicia a geração do horário de forma assíncrona.
// @Accept  json
// @Produce json
// @Param   config body object true "Configuração JSON das turmas, disciplinas e professores"
// @Success 202 {object} Job "Job iniciado com sucesso"
// @Failure 400 {object} object "Erro: corpo da requisição inválido"
// @Router /horario/jobs [post]
func (h *Handlers) StartJobHandler(c *gin.Context) {
	/*var config json.RawMessage
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido: " + err.Error()})
		return
	}*/
	var config horario.ArquivoJson
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido: " + err.Error()})
		return
	}

	// Cria o novo job
	newJob := &Job{
		ID:        uuid.New().String(),
		Status:    StatusPending,
		Config:    config,
		CreatedAt: time.Now(),
	}

	// Adiciona ao armazenamento
	h.Store.AddJob(newJob)

	// Dispara o processamento em uma nova goroutine para não bloquear a resposta
	go h.Store.processHorarioJob(newJob)

	// Retorna uma resposta imediata de "Accepted" com o ID do job
	c.JSON(http.StatusAccepted, newJob)
}

// GetJobStatusHandler lida com a requisição para verificar o status de um job.
// @Summary Verifica o status de um job
// @Description Recupera o status e o resultado (se disponível) de um job de geração de horário.
// @Produce json
// @Param   id   path      string  true  "Job ID"
// @Success 200 {object} Job "Status e resultado do Job"
// @Failure 404 {object} object "Erro: job não encontrado"
// @Router /horario/jobs/{id} [get]
func (h *Handlers) GetJobStatusHandler(c *gin.Context) {
	jobID := c.Param("id")

	job, found := h.Store.GetJob(jobID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job não encontrado"})
		return
	}

	c.JSON(http.StatusOK, job)
}


//================================================================//
// 5. FUNÇÃO PRINCIPAL (Main Function)                            //
//================================================================//

func InitAPIServer() {
	// Inicializa o armazenamento em memória
	jobStore := NewJobStore()

	// Inicializa os handlers com as dependências
	handlers := &Handlers{
		Store: jobStore,
	}

	// Inicializa o router do Gin
	router := gin.Default()
	router.SetTrustedProxies(nil) // Para desenvolvimento

	// Agrupa as rotas da API
	api := router.Group("/horario")
	{
		api.POST("/jobs", handlers.StartJobHandler)
		api.GET("/jobs/:id", handlers.GetJobStatusHandler)
	}
    
    // Rota raiz para health check
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Servidor HorarioGen API está no ar!",
        })
    })

	fmt.Println("Servidor HorarioGen API iniciado em http://localhost:8080")
	// Inicia o servidor na porta 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}