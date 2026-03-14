package acceptance_test

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

func Test主要フロー_作成更新完了再開削除まで一貫して成功すること(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repository := memory.NewTodoRepository()
	createUseCase := app.NewCreateTodoUseCase(repository, func() string { return "todo-1" })
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	updateTitleUseCase := app.NewUpdateTodoTitleUseCase(repository)
	deleteUseCase := app.NewDeleteTodoUseCase(repository)
	reopenUseCase := app.NewReopenTodoUseCase(repository)
	router := httpapi.NewRouter(
		createUseCase,
		completeUseCase,
		listUseCase,
		updateTitleUseCase,
		deleteUseCase,
		reopenUseCase,
	)

	createRes := doJSONRequest(t, router, http.MethodPost, "/todos", map[string]string{"title": "牛乳を買う"})
	if createRes.Code != http.StatusCreated {
		t.Fatalf("Todo作成が失敗: got=%d body=%s", createRes.Code, createRes.Body.String())
	}

	updateRes := doJSONRequest(t, router, http.MethodPatch, "/todos/todo-1/title", map[string]string{"title": "パンを買う"})
	if updateRes.Code != http.StatusOK {
		t.Fatalf("タイトル変更が失敗: got=%d body=%s", updateRes.Code, updateRes.Body.String())
	}

	completeRes := doRequest(router, http.MethodPatch, "/todos/todo-1/complete", nil)
	if completeRes.Code != http.StatusOK {
		t.Fatalf("完了処理が失敗: got=%d body=%s", completeRes.Code, completeRes.Body.String())
	}

	completedListRes := doRequest(router, http.MethodGet, "/todos?completed=true", nil)
	if completedListRes.Code != http.StatusOK {
		t.Fatalf("完了一覧取得が失敗: got=%d body=%s", completedListRes.Code, completedListRes.Body.String())
	}
	var completedList []map[string]any
	if err := json.Unmarshal(completedListRes.Body.Bytes(), &completedList); err != nil {
		t.Fatalf("完了一覧の解析に失敗: %v", err)
	}
	if len(completedList) != 1 || completedList[0]["title"] != "パンを買う" {
		t.Fatalf("完了一覧の内容が不正: %#v", completedList)
	}

	reopenRes := doRequest(router, http.MethodPatch, "/todos/todo-1/reopen", nil)
	if reopenRes.Code != http.StatusOK {
		t.Fatalf("未完了戻しが失敗: got=%d body=%s", reopenRes.Code, reopenRes.Body.String())
	}

	deleteRes := doRequest(router, http.MethodDelete, "/todos/todo-1", nil)
	if deleteRes.Code != http.StatusNoContent {
		t.Fatalf("削除が失敗: got=%d body=%s", deleteRes.Code, deleteRes.Body.String())
	}

	listRes := doRequest(router, http.MethodGet, "/todos", nil)
	if listRes.Code != http.StatusOK {
		t.Fatalf("最終一覧取得が失敗: got=%d body=%s", listRes.Code, listRes.Body.String())
	}
	var list []map[string]any
	if err := json.Unmarshal(listRes.Body.Bytes(), &list); err != nil {
		t.Fatalf("最終一覧の解析に失敗: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("最終的に0件であるべき: %#v", list)
	}
}

func doJSONRequest(t *testing.T, router http.Handler, method string, path string, payload map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("リクエスト生成に失敗: %v", err)
	}
	return doRequest(router, method, path, body)
}

func doRequest(router http.Handler, method string, path string, body []byte) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
