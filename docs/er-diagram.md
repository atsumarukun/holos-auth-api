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

agents {
  char(36) id PK
  char(36) user_id FK
  varchar(255) name
  datetime(6) created_at
  datetime(6) updated_at
  datetime(6) deleted_at
}

policies {
  char(36) id PK
  char(36) user_id FK
  varchar(255) name
  enum service
  varchar(255) path
  json methods
  datetime(6) created_at
  datetime(6) updated_at
  datetime(6) deleted_at
}

permissions {
  char(36) agent_id PK, FK
  char(36) policy_id PK, FK
}

users ||--o| user_tokens: ""

users ||--o{ agents: ""
agents ||--o{ permissions: ""

users ||--o{ policies: ""
policies ||--o{ permissions: ""
```

# テーブル

## users
**ユーザーテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | :---: | --- |
| char(36) | id | PK | | ID |
| varchar(24) | name | UQ | | ユーザー名 |
| varchar(60) | password | | | パスワード |
| datetime(6) | created_at | | | 作成日 |
| datetime(6) | updated_at | | | 更新日 |
| datetime(6) | deleted_at | | * | 削除日 |

## user_tokens
**ユーザートークンテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | --- | --- |
| char(36) | user_id | PK / FK | | ユーザーID |
| char(32) | token | UQ | | トークン |
| datetime(6) | expires_at | | | 有効期限 |

## agents
**エージェントテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | :---: | --- |
| char(36) | id | PK | | ID |
| char(36) | user_id | FK | | ユーザーID |
| varchar(255) | name | | | エージェント名 |
| datetime(6) | created_at | | | 作成日 |
| datetime(6) | updated_at | | | 更新日 |
| datetime(6) | deleted_at | | * | 削除日 |

## policies
**ポリシーテーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | :---: | --- |
| char(36) | id | PK | | ID |
| char(36) | user_id | FK | | ユーザーID |
| varchar(255) | name | | | ポリシー名 |
| enum("STORAGE", "CONTENT") | service | | | サービス |
| varchar(255) | path | | | パス |
| json | methods | | | メソッド |
| datetime(6) | created_at | | | 作成日 |
| datetime(6) | updated_at | | | 更新日 |
| datetime(6) | deleted_at | | * | 削除日 |

## permissions
**権限テーブル**
| type | name | key | nullable | comment |
| --- | --- | --- | :---: | --- |
| char(36) | agent_id | PK, FK | | エージェントID |
| char(36) | policy_id | PK, FK | | ポリシーID |
