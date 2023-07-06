package execution

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octopipe/cloudx/internal/pagination"
)

type httpHandler struct {
	executionUseCase UseCase
}

func NewHTTPHandler(e *gin.Engine, executionUseCase UseCase) *gin.Engine {
	h := httpHandler{executionUseCase: executionUseCase}

	e.GET("/shared-infras/:shared-infra-name/executions", h.List)
	e.POST("/shared-infras/:shared-infra-name/executions", h.Create)
	e.GET("/shared-infras/:shared-infra-name/executions/:name", h.Get)
	e.PUT("/shared-infras/:shared-infra-name/executions/:name", h.Update)
	e.DELETE("/shared-infras/:shared-infra-name/executions/:name", h.Delete)

	return e
}

func (h httpHandler) List(c *gin.Context) {
	var err error
	namespace := "default"
	limit := 10
	chunk := ""

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	if c.Query("chunk") != "" {
		chunk = c.Query("chunk")
	}

	sharedInfraName := c.Param("shared-infra-name")

	list, err := h.executionUseCase.List(c.Request.Context(), sharedInfraName, namespace, pagination.ChunkingPaginationRequest{
		Limit: int64(limit),
		Chunk: chunk,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h httpHandler) Get(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}
	name := c.Param("name")
	sharedInfraName := c.Param("shared-infra-name")

	item, err := h.executionUseCase.Get(c.Request.Context(), sharedInfraName, name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h httpHandler) Create(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}

	execution := Execution{}
	if err := c.BindJSON(&execution); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sharedInfraName := c.Param("shared-infra-name")
	item, err := h.executionUseCase.Create(c.Request.Context(), sharedInfraName, namespace, execution)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h httpHandler) Update(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}

	execution := Execution{}
	if err := c.BindJSON(&execution); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sharedInfraName := c.Param("shared-infra-name")
	item, err := h.executionUseCase.Update(c.Request.Context(), sharedInfraName, namespace, execution)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h httpHandler) Delete(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}
	name := c.Param("name")
	sharedInfraName := c.Param("shared-infra-name")
	err := h.executionUseCase.Delete(c.Request.Context(), sharedInfraName, name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
