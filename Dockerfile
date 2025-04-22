FROM golang:1.24.2 AS builder

COPY ./ /mcp/
WORKDIR /mcp/

ENV GOPATH=/mcp/.go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN/bin:$PATH

RUN go mod tidy && \
    go build && \
    go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
    
EXPOSE 9000
ENTRYPOINT ["/mcp/goreq"]
CMD []