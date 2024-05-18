FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download
RUN make install-tool-deps
RUN go generate ./...
RUN go build ./cmd/blog

CMD [ "./blog" ]