package grpc

import (
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/zeddy-go/zeddy/errx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func simpleInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result any, err error) {
	if req != nil {
		var v *protovalidate.Validator
		v, err = protovalidate.New()
		if err != nil {
			err = errx.Wrap(err, "new validator failed")
			return
		}
		if err = v.Validate(req.(proto.Message)); err != nil {
			err = errx.Wrap(err, "validate failed", errx.WithCode(int(codes.InvalidArgument)))
			return
		}
	}

	// 调用被拦截的方法
	return handler(ctx, req)
}
