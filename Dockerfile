FROM golang:alpine as builder

WORKDIR /GoCommentSvc

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/

FROM alpine

WORKDIR /GoCommentSvc

COPY --from=builder /GoCommentSvc/main /GoCommentSvc/main
COPY --from=builder /GoCommentSvc/pkg/config/envs/*.env /GoCommentSvc/

RUN chmod +x /GoCommentSvc/main

CMD ["./main"]
