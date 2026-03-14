package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	httpapi "github.com/user/ddd/backend/interfaces/http"
)

func postTodos(router http.Handler, title string) *httptest.ResponseRecorder {
	body := map[string]string{"title": title}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	return res
}

func TestPOSTTodosでTodoを作成できること(t *testing.T) {
	repository := httpapi.NewInMemoryTodoRepository()
	usecase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	router := httpapi.NewRouter(usecase)
	res := postTodos(router, "牛乳を買う")

	if res.Code != http.StatusCreated {
		t.Fatalf("201を期待: got=%d", res.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}

	if response["id"] != "todo-1" {
		t.Fatalf("idが一致しない: got=%v", response["id"])
	}
	if response["title"] != "牛乳を買う" {
		t.Fatalf("titleが一致しない: got=%v", response["title"])
	}
	if response["isCompleted"] != false {
		t.Fatalf("isCompletedはfalseのはず: got=%v", response["isCompleted"])
	}
}

func TestPOSTTodosで空タイトルは400になること(t *testing.T) {
	repository := httpapi.NewInMemoryTodoRepository()
	usecase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	router := httpapi.NewRouter(usecase)
	res := postTodos(router, "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("400を期待: got=%d", res.Code)
	}
}
