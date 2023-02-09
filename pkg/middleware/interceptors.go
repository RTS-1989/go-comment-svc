package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"time"
)

const (
	DefaultXRequestIDKey = "x-request-id"
	DefaultXRequestURL   = "x-service-address"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	var requestId string
	requestIdFromContext := metadata.ValueFromIncomingContext(ctx, DefaultXRequestIDKey)

	if len(requestIdFromContext) == 0 {
		requestId = ""
	} else {
		requestId = metadata.ValueFromIncomingContext(ctx, DefaultXRequestIDKey)[0]
	}

	h, err := handler(ctx, req)

	var address string
	addressFromContext := metadata.ValueFromIncomingContext(ctx, DefaultXRequestURL)
	if len(addressFromContext) == 0 {
		address = ""
	} else {
		address = metadata.ValueFromIncomingContext(ctx, DefaultXRequestURL)[0]
	}

	//logging
	f, err := os.OpenFile("go-comment-svc-request.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logger := log.New(f, "RequestLogger: ", log.LstdFlags)

	logger.Printf("request - Address:%s\tDuration:%s\trequestId:%s\tError:%v\n",
		address,
		time.Since(start),
		requestId,
		err)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Got error %v\n", err)
		}
	}(f)

	return h, err
}
