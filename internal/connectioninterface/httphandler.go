package connectioninterface

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octopipe/cloudx/internal/pagination"
)

type httpHandler struct {
	connectionInterfaceUseCase UseCase
}

func NewHTTPHandler(e *gin.Engine, connectionInterfaceUseCase UseCase) *gin.Engine {
	h := httpHandler{connectionInterfaceUseCase: connectionInterfaceUseCase}

	e.GET("/connections-interfaces", h.List)
	e.POST("/connections-interfaces", h.Create)
	e.GET("/connections-interfaces/:name", h.Get)
	e.PUT("/connections-interfaces/:name", h.Update)
	e.DELETE("/connections-interfaces/:name", h.Delete)

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

	list, err := h.connectionInterfaceUseCase.List(c.Request.Context(), namespace, pagination.ChunkingPaginationRequest{
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

	item, err := h.connectionInterfaceUseCase.Get(c.Request.Context(), name, namespace)
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

	connectionInterface := ConnectionInterface{}
	if err := c.BindJSON(&connectionInterface); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println(connectionInterface)

	item, err := h.connectionInterfaceUseCase.Create(c.Request.Context(), connectionInterface)
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

	connectionInterface := ConnectionInterface{}
	if err := c.BindJSON(&connectionInterface); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	item, err := h.connectionInterfaceUseCase.Update(c.Request.Context(), connectionInterface)
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

	err := h.connectionInterfaceUseCase.Delete(c.Request.Context(), name, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
