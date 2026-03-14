# API仕様

## 概要
- Base URL: `http://localhost:8080`
- Content-Type: `application/json`

## エンドポイント

### 1. ヘルスチェック
- `GET /health`
- Response `200`
```json
{ "status": "ok" }
```

### 2. Todo作成
- `POST /todos`
- Request
```json
{ "title": "牛乳を買う" }
```
- Response `201`
```json
{ "id": "todo-1", "title": "牛乳を買う", "isCompleted": false }
```

### 3. Todo一覧
- `GET /todos`
- Query（任意）
  - `completed=true`
  - `completed=false`
- Response `200`
```json
[
  { "id": "todo-1", "title": "牛乳を買う", "isCompleted": false }
]
```
- `completed` が不正値の場合は `400`

### 4. Todo完了
- `PATCH /todos/:id/complete`
- Response `200`
```json
{ "id": "todo-1", "title": "牛乳を買う", "isCompleted": true }
```

### 5. Todo再開（未完了に戻す）
- `PATCH /todos/:id/reopen`
- Response `200`
```json
{ "id": "todo-1", "title": "牛乳を買う", "isCompleted": false }
```

### 6. タイトル変更
- `PATCH /todos/:id/title`
- Request
```json
{ "title": "パンを買う" }
```
- Response `200`
```json
{ "id": "todo-1", "title": "パンを買う", "isCompleted": false }
```

### 7. Todo削除
- `DELETE /todos/:id`
- Response `204`

## エラー方針
- 入力不正: `400`
- 対象なし: `404`
- その他: `500`（一覧取得のみ内部エラーを想定）
