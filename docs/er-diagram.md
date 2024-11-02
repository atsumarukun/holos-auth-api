# ER図

```mermaid
erDiagram

users {
  binary(16) id PK
  varchar(255) name
  varchar(60) password
  datetime(6) created_at
  datetime(6) updated_at
  datetime(6) deleted_at
}
```

# テーブル

## users
**ユーザーテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | --- | --- |
| binary(16) | id | PK | | ID |
| varchar(255) | name | UQ | | ユーザー名 |
| varchar(60) | password | | | パスワード |
| datetime(6) | created_at | | | 作成日 |
| datetime(6) | updated_at | | | 更新日 |
| datetime(6) | deleted_at | | TRUE | 削除日 |
