package grpc

import (
	pb "StealthIMDB/StealthIM.DBGateway"
	"StealthIMDB/errorcode"
	"StealthIMDB/redis"
	"context"
	"log"
	"time"
)

func (s *server) RedisGet(ctx context.Context, in *pb.RedisGetStringRequest) (*pb.RedisGetStringResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Redis]Call Get \"%s\"\n", in.Key)
	}
	cli := redis.GetConn(int(in.DBID))
	if cli == nil {
		return &pb.RedisGetStringResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalComponentError, Msg: "Redis Not Online"}}, nil
	}
	val, err := cli.Get(ctx, in.Key).Result()
	if err != nil {
		return &pb.RedisGetStringResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalNetworkError, Msg: err.Error()}}, nil
	}
	return &pb.RedisGetStringResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}, Value: val}, nil
}
func (s *server) RedisBGet(ctx context.Context, in *pb.RedisGetBytesRequest) (*pb.RedisGetBytesResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Redis]Call BGet \"%s\"\n", in.Key)
	}
	cli := redis.GetConn(int(in.DBID))
	if cli == nil {
		return &pb.RedisGetBytesResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalComponentError, Msg: "Redis Not Online"}}, nil
	}
	val, err := cli.Get(ctx, in.Key).Bytes()
	if err != nil {
		return &pb.RedisGetBytesResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalNetworkError, Msg: err.Error()}}, nil
	}
	return &pb.RedisGetBytesResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}, Value: val}, nil
}
func (s *server) RedisSet(ctx context.Context, in *pb.RedisSetStringRequest) (*pb.RedisSetResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Redis]Call Set \"%s\"\n", in.Key)
	}
	if in.Key == "GatewayInfo" {
		return &pb.RedisSetResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
	}
	cli := redis.GetConn(int(in.DBID))
	if cli == nil {
		return &pb.RedisSetResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalNetworkError, Msg: "Redis Not Online"}}, nil
	}
	err := cli.Set(ctx, in.Key, in.Value, time.Duration(in.Ttl)*time.Second).Err()
	if err != nil {
		return &pb.RedisSetResponse{
			Result: &pb.Result{Code: errorcode.DBGatewayRedisServiceError, Msg: err.Error()}}, nil
	}
	return &pb.RedisSetResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
}
func (s *server) RedisBSet(ctx context.Context, in *pb.RedisSetBytesRequest) (*pb.RedisSetResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Redis]Call BSet \"%s\"\n", in.Key)
	}
	if in.Key == "GatewayInfo" {
		return &pb.RedisSetResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
	}
	cli := redis.GetConn(int(in.DBID))
	if cli == nil {
		return &pb.RedisSetResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalNetworkError, Msg: "Redis Not Online"}}, nil
	}
	err := cli.Set(ctx, in.Key, in.Value, time.Duration(in.Ttl)*time.Second).Err()
	if err != nil {
		return &pb.RedisSetResponse{
			Result: &pb.Result{Code: errorcode.DBGatewayRedisServiceError, Msg: err.Error()}}, nil
	}
	return &pb.RedisSetResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
}
func (s *server) RedisDel(ctx context.Context, in *pb.RedisDelRequest) (*pb.RedisDelResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Redis]Call Del \"%s\"\n", in.Key)
	}
	if in.Key == "GatewayInfo" {
		return &pb.RedisDelResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
	}
	cli := redis.GetConn(int(in.DBID))
	if cli == nil {
		return &pb.RedisDelResponse{
			Result: &pb.Result{Code: errorcode.ServerInternalNetworkError, Msg: "Redis Not Online"}}, nil
	}
	err := cli.Del(ctx, in.Key).Err()
	if err != nil {
		return &pb.RedisDelResponse{
			Result: &pb.Result{Code: errorcode.DBGatewayRedisServiceError, Msg: err.Error()}}, nil
	}
	return &pb.RedisDelResponse{Result: &pb.Result{Code: errorcode.Success, Msg: ""}}, nil
}
