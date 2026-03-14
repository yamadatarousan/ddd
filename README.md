# DDD学習用 Todo アプリ

Domain-Driven Design（ドメイン駆動設計）を学ぶための最小構成アプリです。

- Frontend: TypeScript + React + Vite
- Backend: Go + Gin

## ディレクトリ
- `backend/domain`: ドメインモデル
- `backend/application`: ユースケース
- `backend/infrastructure`: リポジトリ実装
- `backend/interfaces/http`: HTTPルーター
- `backend/acceptance`: 受け入れテスト
- `cmd/server`: サーバー起動エントリ
- `frontend`: 学習用最小画面
- `docs`: API仕様と設計メモ

## セットアップ
### 1. Backend
```bash
make test
make backend-run
```

### 2. Frontend
```bash
make frontend-install
make frontend-run
```

ブラウザで `http://localhost:5173` を開くと利用できます。

## 実行コマンド
- 全テスト: `make test`
- 受け入れテスト: `make test-acceptance`
- Backend起動: `make backend-run`
- Frontend起動: `make frontend-run`

## 学習の進め方
1. `backend/domain` のテストを読む
2. `backend/application` のユースケースを読む
3. `backend/interfaces/http/router.go` でHTTP接続を見る
4. `backend/acceptance/todo_flow_test.go` で全体フローを確認する
5. `docs/ddd.md` と `docs/flow.md` で依存方向を再確認する

## API
主要エンドポイントは `docs/api.md` を参照してください。
