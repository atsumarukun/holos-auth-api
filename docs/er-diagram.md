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
