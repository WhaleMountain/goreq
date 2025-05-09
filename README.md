# GoReq

PlaywrightでURLにアクセスしてコンテンツをMarkdown形式で返却するMCPサーバー

## ビルド

```
docker build -t mcp/goreq .
```

## MCPの設定例

Dockerイメージを作成後、VS Codeのsettings.json(もしくは`.vscode/mcp.json`)を下記のように設定する。
```json
"mcp": {
    "inputs": [],
    "servers": {
        "goreq": {
            "command": "docker",
            "args": [
                "run",
                "--rm",
                "-i",
                "mcp/goreq",
                "-t",
                "stdio"
            ]
        }
    }
}
```

## 直接実行方法

* sse
```
docker run -i --rm -p 9000:9000 mcp/goreq -t sse
```

* stdio
```
docker run -i --rm mcp/goreq -t stdio
```