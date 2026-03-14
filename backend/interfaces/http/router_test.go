package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	app "github.com/user/ddd/backend/application/todo"
	"github.com/user/ddd/backend/infrastructure/memory"
	httpapi "github.com/user/ddd/backend/interfaces/http"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func postTodos(router http.Handler, title string) *httptest.ResponseRecorder {
	body := map[string]string{"title": title}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	return res
}

func newTestRouter(repository *memory.TodoRepository) http.Handler {
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	updateTitleUseCase := app.NewUpdateTodoTitleUseCase(repository)
	deleteUseCase := app.NewDeleteTodoUseCase(repository)
	reopenUseCase := app.NewReopenTodoUseCase(repository)
	return httpapi.NewRouter(createUseCase, completeUseCase, listUseCase, updateTitleUseCase, deleteUseCase, reopenUseCase)
}

func TestPOSTTodosでTodoを作成できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)
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
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)
	res := postTodos(router, "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("400を期待: got=%d", res.Code)
	}
}

func TestPATCHTodosCompleteでTodoを完了にできること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

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
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	req := httptest.NewRequest(http.MethodPatch, "/todos/not-found/complete", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("404を期待: got=%d", res.Code)
	}
}

func TestGETTodosでTodo一覧を取得できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

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

func TestPATCHTodosTitleでタイトル変更できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	body := map[string]string{"title": "パンを買う"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/todos/todo-1/title", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("200を期待: got=%d", res.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if response["title"] != "パンを買う" {
		t.Fatalf("titleが更新されていない: got=%v", response["title"])
	}
}

func TestPATCHTodosTitleで空タイトルは400になること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	body := map[string]string{"title": ""}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/todos/todo-1/title", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("400を期待: got=%d", res.Code)
	}
}

func TestDELETETodosでTodoを削除できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	req := httptest.NewRequest(http.MethodDelete, "/todos/todo-1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("204を期待: got=%d", res.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/todos", nil)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("一覧取得が失敗: got=%d", listRes.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(listRes.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if len(response) != 0 {
		t.Fatalf("削除後は0件であるべき: got=%d", len(response))
	}
}

func TestDELETETodosで存在しないTodoは404になること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	req := httptest.NewRequest(http.MethodDelete, "/todos/not-found", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusNotFound {
		t.Fatalf("404を期待: got=%d", res.Code)
	}
}

func TestGETHealthで稼働確認できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("200を期待: got=%d", res.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if response["status"] != "ok" {
		t.Fatalf("statusが一致しない: got=%s", response["status"])
	}
}

func TestPATCHTodosReopenでTodoを未完了に戻せること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	createRes := postTodos(router, "牛乳を買う")
	if createRes.Code != http.StatusCreated {
		t.Fatalf("事前作成が失敗: got=%d", createRes.Code)
	}

	completeReq := httptest.NewRequest(http.MethodPatch, "/todos/todo-1/complete", nil)
	completeRes := httptest.NewRecorder()
	router.ServeHTTP(completeRes, completeReq)
	if completeRes.Code != http.StatusOK {
		t.Fatalf("完了処理が失敗: got=%d", completeRes.Code)
	}

	reopenReq := httptest.NewRequest(http.MethodPatch, "/todos/todo-1/reopen", nil)
	reopenRes := httptest.NewRecorder()
	router.ServeHTTP(reopenRes, reopenReq)
	if reopenRes.Code != http.StatusOK {
		t.Fatalf("200を期待: got=%d", reopenRes.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(reopenRes.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗: %v", err)
	}
	if response["isCompleted"] != false {
		t.Fatalf("isCompletedはfalseのはず: got=%v", response["isCompleted"])
	}
}

func TestPATCHTodosReopenで存在しないTodoは404になること(t *testing.T) {
	repository := memory.NewTodoRepository()
	router := newTestRouter(repository)

	req := httptest.NewRequest(http.MethodPatch, "/todos/not-found/reopen", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusNotFound {
		t.Fatalf("404を期待: got=%d", res.Code)
	}
}
