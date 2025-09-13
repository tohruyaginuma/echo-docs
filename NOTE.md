# Notes

## Basic Syntax

```go
type User struct {
    // ”json:” 部分をタグと呼ぶ
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

type User struct {
  ID string `query:"id"`
}

// in the handler for /users?id=<userID>
var user User
err := c.Bind(&u  ser); if err != nil {
    return c.String(http.StatusBadRequest, "bad request")
}
```

## Binding

- ユーザーからのデータを処理するプロセスのこと
- 下記のようなデータを扱う
  - URL Path parameter
  - URL Query parameter
  - Header
  - Request body
- バインド順

1. Path parameters
2. Query parameters
3. Request body
   同じ名前のパラメータがあれば上書きする

- 基本的には DTO（ユーザーとのインタラクトが発生する struct）とドメインモデル（ビジネス構造体。内部でしか使用しない構造体）は分けて定義する。セキュリティのため
- `json:"-"`バインドしたくないフィールどはハイフンでむしできる

`e.Use`: ミドルウェアを登録するためのメソッド
`echo.HandlerFunc` : 型定義。ハンドラー関数

`Context.RealIP()` ジオベースアクセス解析、アクセスコントロール、Auditing のために使う
Auditing ：システム上でアクティビティーを計測、評価、分析すること。eg: システムログを取ること

## Routing

### Path Matching Order

1. `/user/new` // Static
2. `/user/:id` // Param
3. `/user/1/files/*` // Match any

## Middle ware

- Basic Auth
  - HTTP のベーシックな Auth を提供
  - invalid or missing で 401 - Unauthorized を返す
- Body Dump
  - Request, Response の Payload をキャプチャしてハンドラーを呼ぶ
  - デバッグやロギング目的で使用する
- Body Limit
  - Request body の最大サイズを設定できる
  - 超えた場合 413 - Request Entity Too Large を返す
- Casbin Auth
  - オープソースライブラリ for Go
  - 様々なモデルで Auth を強要する
  - アクセスコントロールを提供する
    - ACL
    - super user | user
  - ベーシックな HTTP 認証のみサポート
- Context Timeout
  - タイムアウトの時間を定義できる
- CORS
  - CORS 設定
- CSRF
  - Cross-site request forgery
  - CSFR token は Echo#COntext から利用される。templete を通じ、Context key を用いてクライアントへ渡す
  - CSFR token は CSRF cookie からアクセスされる
- Decompress
  - Content -Encording header が gzip にセットされている場合、それを解凍する
- Gzip
  - HTTP レスポンスを zip を用いてコンプレッションする
- Jaeger [Community contribution]
  - Jaeger トレースミドルウェアを用いて、Echo フレームワークでトレースのリクエスト
- JWT
  - 有効なトークンはハンドラーを実行
  - 有効でないトークンは 401Unauthorized を返す
  - Authorization header が無効な場合 400 Bad Request
- Key Auth
  - key based authentication を提供
    - クライアントがリクエストエッダーにサーバーから発行したキーを載せる
    - 一致すれば認証。
    - キーが盗まれたら終了なのでガバガバ
  - 有効なキーでハンドラを実行
  - 無効キーで 401 を返す
  - キー不在で 400 を返す
- Loger
  - 各 HTTP リクエストに関する情報をログする
  - ２種類ある
    - Logger
      - String based logger
      - easy to start
    - RequestLogger
      - Customizable function based logger
      - カスタムできる何をどうログするかを
      - ３ rd パーティーのロガーと使い方が合う
- Method Override
  - リクエストのメソッドを上書きするミドルウェア
  - クライアントなど、特定のメソッドしか受け付けていない時に利用する
- Promethues [Community contribution]
  - HTTP リクエストの統計を生成
- Proxy
  - HTTP, Web Socket のリバースプロキシを提供
  - 設定されたロードバランシング技術でその先のサーバーへリクエストを転送
- Rate Limiter
  - 特定 IP や id の期間内のリクエストの量を制限できる
  - ディフォルトでは、リクエストのトラックをインメモリーに保存
  - しかし、ディフォルトは正確性を重んじており、大きな量の連続したリクエストや複数の source からのリクエストには向かない
- Recover
  - チェーン内のパニックから救う、スタックをトレースし、制御を一元化する
- Redirect
  - https や www のリダイレクトなどをマネージする。
- Request ID
  - リクエストに対して固有の ID を付与できる
- Rewrite
  - 提供されたルールに基づき、URL path をリライトする
- Secure
  - ブラウザが危険な動きをしないようにするための安全ヘッダーをまとめて付与
  - XSS attack、content type sniffing, clickjacking, insecure connection, other code injection attacks を防ぐ。
- Session [Community contribution]
  - http セッション管理を提供
  - gorilla sessions が利用されている
  - ディフォルトではクッキーとファイルベース ° セッションストアが提供される
- Static
  - ルートディレクトリからスタティックファイルをプロバイドできる
- Trailing Slash
  - リクエスト UR にトレイリングスラッシュを付与 or 削除する
