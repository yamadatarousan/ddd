# DDD学習メモ

## レイヤー構成
- `backend/domain`
  - エンティティと値オブジェクトの不変条件を定義
- `backend/application`
  - ユースケース単位の操作を定義
  - リポジトリはインターフェースで依存逆転
- `backend/infrastructure`
  - インターフェースの実装（現在はインメモリ）
- `backend/interfaces/http`
  - HTTP入力をユースケースコマンドに変換
  - ユースケース結果をレスポンスへ変換

## 学習ポイント
- ドメイン制約は `domain` で集中管理する
- ユースケースは「何をするか」に集中し、保存方法を知らない
- インフラ実装を差し替えてもユースケースは変更不要
- HTTP層はアプリケーション層の呼び出しに専念し、ドメイン判断を持ち込まない

## 依存方向
`interfaces/http -> application -> domain`
`infrastructure -> domain`

インフラはアプリケーションのインターフェースを満たすだけで、
ユースケース側は具体実装を知らない。
