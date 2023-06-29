package providerconfig

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octopipe/cloudx/internal/pagination"
)

type httpHandler struct {
	providerConfigUseCase UseCase
}

func NewHTTPHandler(e *gin.Engine, providerConfigUseCase UseCase) *gin.Engine {
	h := httpHandler{providerConfigUseCase: providerConfigUseCase}

	e.GET("/providers-configs", h.List)
	e.POST("/providers-configs", h.Create)
	e.GET("/providers-configs/:name", h.Get)
	e.PUT("/providers-configs/:name", h.Update)
	e.DELETE("/providers-configs/:name", h.Delete)

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

	list, err := h.providerConfigUseCase.List(c.Request.Context(), namespace, pagination.ChunkingPaginationRequest{
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

	item, err := h.providerConfigUseCase.Get(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h httpHandler) Create(c *gin.Context) {
	// namespace := "default"

	// if c.Query("namespace") != "" {
	// 	namespace = c.Query("namespace")
	// }
	// name := c.Param("name")

	providerConfig := ProviderConfig{}
	if err := c.BindJSON(&providerConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.providerConfigUseCase.Create(c.Request.Context(), providerConfig)
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
	// name := c.Param("name")

	providerConfig := ProviderConfig{}
	if err := c.BindJSON(&providerConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.providerConfigUseCase.Update(c.Request.Context(), providerConfig)
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

	err := h.providerConfigUseCase.Delete(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
