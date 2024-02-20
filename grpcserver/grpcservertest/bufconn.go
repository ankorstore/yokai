package grpcservertest

import "google.golang.org/grpc/test/bufconn"

func NewBufconnListener(size int) *bufconn.Listener {
	return bufconn.Listen(size)
}
