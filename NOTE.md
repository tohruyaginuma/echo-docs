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
