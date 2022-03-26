# create-gin-app
Gin(Go)のWebアプリケーション雛形作成スクリプト。
ディレクトリ構成 + Signup/Login/Logout 機能を画面およびサーバプログラム自動生成。

インストール後 ~ /ceate-gin-app/bin にPATHを通し、下記コマンド実行。
```
create-gin-app <appname>
```

スクリプト実行途中にJWT認証用の秘密鍵の入力が要求されるため、
ランダム文字列を入力。( .envファイルに設定される ）
```
Please enter JWT_SECRET_KEY:
```
