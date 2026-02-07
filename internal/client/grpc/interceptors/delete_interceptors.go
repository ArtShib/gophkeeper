package interceptors

//
//import (
//	"time"
//
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"
//)
//
////
////import (
////	"context"
////
////	"google.golang.org/grpc"
////	"google.golang.org/grpc/metadata"
////)
////
////type TokenProvider interface {
////	GetToken() string
////}
////
////func AuthInterceptor(provider TokenProvider) grpc.UnaryClientInterceptor {
////	return func(
////		ctx context.Context,
////		method string,
////		req, reply interface{},
////		cc *grpc.ClientConn,
////		invoker grpc.UnaryInvoker,
////		opts ...grpc.CallOption,
////	) error {
////		token := provider.GetToken()
////
////		if token != "" {
////			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
////		}
////
////		return invoker(ctx, method, req, reply, cc, opts...)
////	}
////}
//
//
//package interceptors
//
//import (
//"context"
//"time"
//
//"google.golang.org/grpc"
//"google.golang.org/grpc/metadata"
//)
//
//// AuthInterceptor - добавляет токен авторизации
//func AuthInterceptor(token string) grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{},
//		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//
//		ctx = metadata.AppendToOutgoingContext(ctx,
//			"access_token", "Bearer "+token,
//			"client-timestamp", time.Now().Format(time.RFC3339),
//		)
//
//		// Выполняем вызов с обновленным контекстом
//		return invoker(ctx, method, req, reply, cc, opts...)
//	}
//}
//
//
//type TokenManager struct {
//	token   string
//	refresh func() (string, error)
//	logger  *slog.Logger
//}
//
//func (tm *TokenManager) AuthInterceptor() grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{},
//		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//
//		package interceptors
//
//import (
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/metadata"
//	"google.golang.org/grpc/status"
//)

//
//import (
//	"time"
//
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"
//)
//
////
////import (
////	"context"
////
////	"google.golang.org/grpc"
////	"google.golang.org/grpc/metadata"
////)
////
////type TokenProvider interface {
////	GetToken() string
////}
////
////func AuthInterceptor(provider TokenProvider) grpc.UnaryClientInterceptor {
////	return func(
////		ctx context.Context,
////		method string,
////		req, reply interface{},
////		cc *grpc.ClientConn,
////		invoker grpc.UnaryInvoker,
////		opts ...grpc.CallOption,
////	) error {
////		token := provider.GetToken()
////
////		if token != "" {
////			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
////		}
////
////		return invoker(ctx, method, req, reply, cc, opts...)
////	}
////}
//
//
//package interceptors
//
//import (
//"context"
//"time"
//
//"google.golang.org/grpc"
//"google.golang.org/grpc/metadata"
//)
//
//// AuthInterceptor - добавляет токен авторизации
//func AuthInterceptor(token string) grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{},
//		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//
//		ctx = metadata.AppendToOutgoingContext(ctx,
//			"access_token", "Bearer "+token,
//			"client-timestamp", time.Now().Format(time.RFC3339),
//		)
//
//		// Выполняем вызов с обновленным контекстом
//		return invoker(ctx, method, req, reply, cc, opts...)
//	}
//}
//
//
//type TokenManager struct {
//	token   string
//	refresh func() (string, error)
//	logger  *slog.Logger
//}
//
//func (tm *TokenManager) AuthInterceptor() grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{},
//		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//
//		for {
//			ctx = metadata.AppendToOutgoingContext(ctx, "access_token", tm.token)
//
//			err := invoker(ctx, method, req, reply, cc, opts...)
//
//			// Если токен истек - обновляем
//			if status.Code(err) == codes.Unauthenticated {
//				tm.token, err = tm.refresh()
//				if err != nil {
//					return err
//				}
//				continue // Повторяем запрос
//			}
//
//			return err
//		}
//	}
//}
//for {
//			ctx = metadata.AppendToOutgoingContext(ctx, "access_token", tm.token)
//
//			err := invoker(ctx, method, req, reply, cc, opts...)
//
//			// Если токен истек - обновляем
//			if status.Code(err) == codes.Unauthenticated {
//				tm.token, err = tm.refresh()
//				if err != nil {
//					return err
//				}
//				continue // Повторяем запрос
//			}
//
//			return err
//		}
//	}
//}
