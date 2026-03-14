package http

import (
	"net/http"

	app "github.com/user/ddd/backend/application/todo"
	"github.com/gin-gonic/gin"
)

type createTodoRequest struct {
	Title string `json:"title"`
}

type createTodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

func NewRouter(usecase app.CreateTodoUseCase) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/todos", func(c *gin.Context) {
		var req createTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストです"})
			return
		}

		entity, err := usecase.Execute(app.CreateTodoCommand{Title: req.Title})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createTodoResponse{
			ID:          entity.ID(),
			Title:       entity.Title().Value(),
			IsCompleted: entity.IsCompleted(),
		})
	})

	return router
}
