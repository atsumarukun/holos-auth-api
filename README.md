# 認証/認可 API

## 環境構築

### Dev Containers

開発環境には[Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)を利用.

Visual Studio Code拡張機能.<br />
https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers

## ルール

### ブランチ名

```
issue-${ISSUE_NUMBER}   // issue-1
```

### コミットメッセージ

以下の表に記載するprefixを記載する.

| prefix | content |
| --- | --- |
| create | 機能作成 |
| update | 機能更新 |
| remove | 機能削除 |
| refactor | リファクタリング |
| fix | 不具合修正 |

```
create: 新規機能を作成.
```

## デプロイ

mainブランチへmergeすることでデプロイが行われる.

デプロイされるリソースは以下の通り.

| resource | destination |
| --- | --- |
| API | GitHub Container Registry |
| SwaggerUI | GitHub Pages |

### API

APIのデプロイを行う前に.github/workflows/deploy-api.yml内のapiバージョンとSwaggerUI内のapiバージョンを更新する.<br />
https://github.com/atsumarukun/holos-auth-api/pkgs/container/holos-auth-api

### SwaggerUI

https://atsumarukun.github.io/holos-auth-api/

## データベース接続

本開発環境はDBコンテナをhostにポートフォワードしていない.<br />
データベースクライアントにadminerを利用する場合、以下のdocker-compose.ymlを作成し起動する.

```yml
services:
  adminer:
    image: adminer
    restart: always
    networks:
      - nw-holos
    ports:
      - 8080:8080

networks:
  nw-holos:
    external: true

```
