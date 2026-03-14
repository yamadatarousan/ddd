package http

import (
	"errors"
	"net/http"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
	"github.com/gin-gonic/gin"
)

type createTodoRequest struct {
	Title string `json:"title"`
}

type updateTodoTitleRequest struct {
	Title string `json:"title"`
}

type createTodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

func toTodoResponse(entity domain.Entity) createTodoResponse {
	return createTodoResponse{
		ID:          entity.ID(),
		Title:       entity.Title().Value(),
		IsCompleted: entity.IsCompleted(),
	}
}

func writeUseCaseError(c *gin.Context, err error) {
	if errors.Is(err, app.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func NewRouter(
	createUseCase app.CreateTodoUseCase,
	completeUseCase app.CompleteTodoUseCase,
	listUseCase app.ListTodoUseCase,
	updateTitleUseCase app.UpdateTodoTitleUseCase,
	deleteUseCase app.DeleteTodoUseCase,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/todos", func(c *gin.Context) {
		var req createTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストです"})
			return
		}

		entity, err := createUseCase.Execute(app.CreateTodoCommand{Title: req.Title})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toTodoResponse(entity))
	})

	router.PATCH("/todos/:id/complete", func(c *gin.Context) {
		entity, err := completeUseCase.Execute(app.CompleteTodoCommand{ID: c.Param("id")})
		if err != nil {
			writeUseCaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTodoResponse(entity))
	})

	router.GET("/todos", func(c *gin.Context) {
		entities, err := listUseCase.Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]createTodoResponse, 0, len(entities))
		for _, entity := range entities {
			response = append(response, toTodoResponse(entity))
		}
		c.JSON(http.StatusOK, response)
	})

	router.PATCH("/todos/:id/title", func(c *gin.Context) {
		var req updateTodoTitleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストです"})
			return
		}

		entity, err := updateTitleUseCase.Execute(app.UpdateTodoTitleCommand{
			ID:    c.Param("id"),
			Title: req.Title,
		})
		if err != nil {
			writeUseCaseError(c, err)
			return
		}

		c.JSON(http.StatusOK, toTodoResponse(entity))
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		err := deleteUseCase.Execute(app.DeleteTodoCommand{ID: c.Param("id")})
		if err != nil {
			writeUseCaseError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})

	return router
}
