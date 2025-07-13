package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/erickweil/horariogen/horario"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/erickweil/horariogen/docs" // Importa os docs gerados pelo swag
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//================================================================//
// 4. HANDLERS DA API (API Handlers)                              //
//================================================================//

// Handlers contém dependências para os handlers da API, como o JobStore.
type Handlers struct {
	Store *JobStore
}

// StartJobHandler lida com a requisição para iniciar um novo job.
// @Summary      Inicia um novo processamento de horário
// @Description  Recebe uma configuração JSON e inicia a geração do horário de forma assíncrona.
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Param        config  body      horario.ArquivoJson  true  "Configuração JSON das turmas, disciplinas e professores"
// @Success      202     {object}  Job                  "Job iniciado com sucesso"
// @Failure      400     {object}  object{error=string} "Erro: corpo da requisição inválido"
// @Router       /horario/jobs [post]
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
// @Summary      Verifica o status de um job
// @Description  Recupera o status e o resultado (se disponível) de um job de geração de horário.
// @Tags         Jobs
// @Produce      json
// @Param        id   path      string               true  "Job ID"
// @Success      200  {object}  Job                  "Status e resultado do Job"
// @Failure      404  {object}  object{error=string} "Erro: job não encontrado"
// @Router       /horario/jobs/{id} [get]
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


// @title           API HorarioGen
// @version         1.0
// @description     Esta é uma API para geração de horários escolares de forma assíncrona.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Seu Nome
// @contact.url    http://www.seusite.com
// @contact.email  seu@email.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Redirect para a documentação Swagger
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
    
    fmt.Println("Servidor HorarioGen API iniciado em http://localhost:8080")

	// Inicia o servidor na porta 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}