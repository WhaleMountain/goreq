# GoReq

PlaywrightでURLにアクセスしてコンテンツをMarkdown形式で返却するMCPサーバー

## ビルド

```
docker build -t mcp/goreq .
```

## 実行方法

* sse
```
docker run -i --rm -p 9000:9000 mcp/goreq -t sse
```

* stdio
```
docker run -i --rm mcp/goreq -t stdio
```