# ER図

```mermaid
erDiagram

users {
  char(36) id PK
  varchar(24) name
  varchar(60) password
  datetime(6) created_at
  datetime(6) updated_at
  datetime(6) deleted_at
}

user_tokens {
  char(36) user_id PK, FK
  char(32) token
  datetime(6) expires_at
}
```

# テーブル

## users
**ユーザーテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | --- | --- |
| char(36) | id | PK | | ID |
| varchar(24) | name | UQ | | ユーザー名 |
| varchar(60) | password | | | パスワード |
| datetime(6) | created_at | | | 作成日 |
| datetime(6) | updated_at | | | 更新日 |
| datetime(6) | deleted_at | | TRUE | 削除日 |

## user_tokens
**ユーザートークンテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | --- | --- |
| char(36) | user_id | PK / FK | | ユーザーID |
| char(32) | token | UQ | | トークン |
| datetime(6) | expires_at | | | 有効期限 |
