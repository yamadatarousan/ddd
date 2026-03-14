import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  completeTodo,
  createTodo,
  deleteTodo,
  fetchTodos,
  reopenTodo,
  Todo,
  updateTodoTitle
} from "./api";

type Filter = "all" | "completed" | "active";

function toCompletedValue(filter: Filter): boolean | undefined {
  if (filter === "all") return undefined;
  return filter === "completed";
}

export function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTitle, setNewTitle] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingTitle, setEditingTitle] = useState("");
  const [filter, setFilter] = useState<Filter>("all");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function loadTodos(nextFilter: Filter = filter): Promise<void> {
    setIsLoading(true);
    setError(null);
    try {
      const list = await fetchTodos(toCompletedValue(nextFilter));
      setTodos(list);
    } catch (err) {
      setError(err instanceof Error ? err.message : "一覧取得に失敗しました");
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    void loadTodos("all");
    // 初回だけ読み込む。filterの変更はselectハンドラーで明示的に反映する。
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function onCreate(event: FormEvent): Promise<void> {
    event.preventDefault();
    setError(null);
    try {
      await createTodo(newTitle);
      setNewTitle("");
      await loadTodos();
    } catch (err) {
      setError(err instanceof Error ? err.message : "作成に失敗しました");
    }
  }

  async function toggleComplete(todo: Todo): Promise<void> {
    setError(null);
    try {
      if (todo.isCompleted) {
        await reopenTodo(todo.id);
      } else {
        await completeTodo(todo.id);
      }
      await loadTodos();
    } catch (err) {
      setError(err instanceof Error ? err.message : "状態変更に失敗しました");
    }
  }

  async function onDelete(id: string): Promise<void> {
    setError(null);
    try {
      await deleteTodo(id);
      await loadTodos();
    } catch (err) {
      setError(err instanceof Error ? err.message : "削除に失敗しました");
    }
  }

  function startEdit(todo: Todo): void {
    setEditingId(todo.id);
    setEditingTitle(todo.title);
  }

  async function saveEdit(id: string): Promise<void> {
    setError(null);
    try {
      await updateTodoTitle(id, editingTitle);
      setEditingId(null);
      setEditingTitle("");
      await loadTodos();
    } catch (err) {
      setError(err instanceof Error ? err.message : "タイトル更新に失敗しました");
    }
  }

  const title = useMemo(() => {
    if (filter === "completed") return "完了Todo";
    if (filter === "active") return "未完了Todo";
    return "すべてのTodo";
  }, [filter]);

  return (
    <main className="page">
      <section className="panel">
        <h1>DDD学習用 Todo</h1>
        <p className="description">Domain-Driven Design のユースケース接続を確認するための最小画面</p>

        <form className="create-form" onSubmit={(event) => void onCreate(event)}>
          <input
            value={newTitle}
            onChange={(event) => setNewTitle(event.target.value)}
            placeholder="Todoタイトル"
            aria-label="todo-title"
          />
          <button type="submit">追加</button>
        </form>

        <div className="toolbar">
          <label htmlFor="filter">表示:</label>
          <select
            id="filter"
            value={filter}
            onChange={(event) => {
              const next = event.target.value as Filter;
              setFilter(next);
              void loadTodos(next);
            }}
          >
            <option value="all">すべて</option>
            <option value="active">未完了</option>
            <option value="completed">完了</option>
          </select>
          <span>{title}</span>
        </div>

        {isLoading ? <p>読み込み中...</p> : null}
        {error ? <p className="error">{error}</p> : null}

        <ul className="list">
          {todos.map((todo) => (
            <li key={todo.id} className="item">
              <button type="button" onClick={() => void toggleComplete(todo)}>
                {todo.isCompleted ? "未完了に戻す" : "完了にする"}
              </button>

              {editingId === todo.id ? (
                <>
                  <input
                    value={editingTitle}
                    onChange={(event) => setEditingTitle(event.target.value)}
                    aria-label={`edit-${todo.id}`}
                  />
                  <button type="button" onClick={() => void saveEdit(todo.id)}>
                    保存
                  </button>
                  <button type="button" onClick={() => setEditingId(null)}>
                    キャンセル
                  </button>
                </>
              ) : (
                <>
                  <span className={todo.isCompleted ? "completed" : ""}>{todo.title}</span>
                  <button type="button" onClick={() => startEdit(todo)}>
                    編集
                  </button>
                </>
              )}

              <button type="button" onClick={() => void onDelete(todo.id)}>
                削除
              </button>
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
