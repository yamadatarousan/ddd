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
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	router := httpapi.NewRouter(createUseCase, completeUseCase, listUseCase)
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
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	router := httpapi.NewRouter(createUseCase, completeUseCase, listUseCase)
	res := postTodos(router, "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("400を期待: got=%d", res.Code)
	}
}

func TestPATCHTodosCompleteでTodoを完了にできること(t *testing.T) {
	repository := httpapi.NewInMemoryTodoRepository()
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	router := httpapi.NewRouter(createUseCase, completeUseCase, listUseCase)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	req := httptest.NewRequest(http.MethodPatch, "/todos/todo-1/complete", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("200を期待: got=%d", res.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if response["isCompleted"] != true {
		t.Fatalf("isCompletedはtrueのはず: got=%v", response["isCompleted"])
	}
}

func TestPATCHTodosCompleteで存在しないTodoは404になること(t *testing.T) {
	repository := httpapi.NewInMemoryTodoRepository()
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	router := httpapi.NewRouter(createUseCase, completeUseCase, listUseCase)

	req := httptest.NewRequest(http.MethodPatch, "/todos/not-found/complete", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("404を期待: got=%d", res.Code)
	}
}

func TestGETTodosでTodo一覧を取得できること(t *testing.T) {
	repository := httpapi.NewInMemoryTodoRepository()
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	router := httpapi.NewRouter(createUseCase, completeUseCase, listUseCase)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("200を期待: got=%d", res.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("1件を期待: got=%d", len(response))
	}
	if response[0]["title"] != "牛乳を買う" {
		t.Fatalf("titleが一致しない: got=%v", response[0]["title"])
	}
}
