## installation

```
$ brew cask install kindlegen
```

config/subscribe.yml
```
n_codes: [string] # Nコードの配列
email: string,
send_to_kindle_email: string,
smtp_user_name: string # gmailの場合、emailと同じ
smtp_password: string
smtp_host: string # smtp.gmail.com
smtp_port: int # 587
```

emailがGMailのとき、smtp_passwordに普段使っているパスワードを入れると警告が出る。
https://security.google.com/settings/security/apppasswords
ここから専用のパスワードを発行する。

起動するとlocalhost:1323でサーバが立っているので、そこから操作できる
