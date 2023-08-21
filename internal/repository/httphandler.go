package repository

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octopipe/cloudx/internal/pagination"
)

type httpHandler struct {
	repositoryUseCase UseCase
}

func NewHTTPHandler(e *gin.Engine, repositoryUseCase UseCase) *gin.Engine {
	h := httpHandler{repositoryUseCase: repositoryUseCase}

	e.GET("/repositories", h.List)
	e.POST("/repositories", h.Create)
	e.GET("/repositories/:repository-name", h.Get)
	e.PUT("/repositories/:repository-name", h.Update)
	e.PATCH("/repositories/:repository-name/sync", h.Sync)
	e.DELETE("/repositories/:repository-name", h.Delete)

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

	list, err := h.repositoryUseCase.List(c.Request.Context(), namespace, pagination.ChunkingPaginationRequest{
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
	name := c.Param("repository-name")

	item, err := h.repositoryUseCase.Get(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h httpHandler) Sync(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}
	name := c.Param("repository-name")

	err := h.repositoryUseCase.Sync(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h httpHandler) Create(c *gin.Context) {
	// namespace := "default"

	// if c.Query("namespace") != "" {
	// 	namespace = c.Query("namespace")
	// }
	// name := c.Param("repository-name")

	repository := Repository{}
	if err := c.BindJSON(&repository); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.repositoryUseCase.Create(c.Request.Context(), repository)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h httpHandler) Update(c *gin.Context) {
	// namespace := "default"

	// if c.Query("namespace") != "" {
	// 	namespace = c.Query("namespace")
	// }
	// name := c.Param("repository-name")

	repository := Repository{}
	if err := c.BindJSON(&repository); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.repositoryUseCase.Update(c.Request.Context(), repository)
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
	name := c.Param("repository-name")

	err := h.repositoryUseCase.Delete(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
