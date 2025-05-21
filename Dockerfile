FROM golang:1.24.3

WORKDIR /mcp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOPATH=/mcp/.go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

RUN go mod tidy && \
    go build -o goreq ./cmd/goreq && \
    go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5101.0 install --with-deps
    
EXPOSE 9000
ENTRYPOINT ["/mcp/goreq"]
CMD []