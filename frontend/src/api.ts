export type Todo = {
  id: string;
  title: string;
  isCompleted: boolean;
};

type JsonValue = Record<string, unknown> | Array<unknown>;

async function request<T extends JsonValue | void>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...(options?.headers ?? {})
    },
    ...options
  });

  if (!response.ok) {
    let message = `HTTP ${response.status}`;
    try {
      const body = (await response.json()) as { error?: string };
      if (typeof body.error === "string" && body.error.length > 0) {
        message = body.error;
      }
    } catch {
      // エラー時にJSONでないレスポンスでも最低限のメッセージを返す。
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export async function fetchTodos(completed?: boolean): Promise<Todo[]> {
  const query = completed === undefined ? "" : `?completed=${completed}`;
  return await request<Todo[]>(`/todos${query}`);
}

export async function createTodo(title: string): Promise<Todo> {
  return await request<Todo>("/todos", {
    method: "POST",
    body: JSON.stringify({ title })
  });
}

export async function completeTodo(id: string): Promise<Todo> {
  return await request<Todo>(`/todos/${id}/complete`, { method: "PATCH" });
}

export async function reopenTodo(id: string): Promise<Todo> {
  return await request<Todo>(`/todos/${id}/reopen`, { method: "PATCH" });
}

export async function updateTodoTitle(id: string, title: string): Promise<Todo> {
  return await request<Todo>(`/todos/${id}/title`, {
    method: "PATCH",
    body: JSON.stringify({ title })
  });
}

export async function deleteTodo(id: string): Promise<void> {
  await request<void>(`/todos/${id}`, { method: "DELETE" });
}
