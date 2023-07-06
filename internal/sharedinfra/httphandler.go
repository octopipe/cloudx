package sharedinfra

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octopipe/cloudx/internal/pagination"
)

type httpHandler struct {
	sharedInfraUseCase UseCase
}

func NewHTTPHandler(e *gin.Engine, sharedInfraUseCase UseCase) *gin.Engine {
	h := httpHandler{sharedInfraUseCase: sharedInfraUseCase}

	e.GET("/shared-infras", h.List)
	e.POST("/shared-infras", h.Create)
	e.GET("/shared-infras/:shared-infra-name", h.Get)
	e.PUT("/shared-infras/:shared-infra-name", h.Update)
	e.PATCH("/shared-infras/:shared-infra-name/reconcile", h.Reconcile)
	e.DELETE("/shared-infras/:shared-infra-name", h.Delete)

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

	list, err := h.sharedInfraUseCase.List(c.Request.Context(), namespace, pagination.ChunkingPaginationRequest{
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
	name := c.Param("shared-infra-name")

	item, err := h.sharedInfraUseCase.Get(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h httpHandler) Reconcile(c *gin.Context) {
	namespace := "default"

	if c.Query("namespace") != "" {
		namespace = c.Query("namespace")
	}
	name := c.Param("shared-infra-name")

	err := h.sharedInfraUseCase.Reconcile(c.Request.Context(), name, namespace)
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
	// name := c.Param("shared-infra-name")

	sharedInfra := SharedInfra{}
	if err := c.BindJSON(&sharedInfra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.sharedInfraUseCase.Create(c.Request.Context(), sharedInfra)
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
	// name := c.Param("shared-infra-name")

	sharedInfra := SharedInfra{}
	if err := c.BindJSON(&sharedInfra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.sharedInfraUseCase.Update(c.Request.Context(), sharedInfra)
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
	name := c.Param("shared-infra-name")

	err := h.sharedInfraUseCase.Delete(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
