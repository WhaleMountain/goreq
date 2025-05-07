FROM golang:1.24.3

COPY ./ /mcp/
WORKDIR /mcp/

ENV GOPATH=/mcp/.go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

RUN go mod tidy && \
    go build && \
    go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5101.0 install --with-deps
    
EXPOSE 9000
ENTRYPOINT ["/mcp/goreq"]
CMD []