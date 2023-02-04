package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
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
	log.Printf("request - Address:%s\tDuration:%s\trequestId:%s\tError:%v\n",
		address,
		time.Since(start),
		requestId,
		err)

	return h, err
}
