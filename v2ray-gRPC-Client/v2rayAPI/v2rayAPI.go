package v2rayAPI

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	statsService "v2ray.com/core/app/stats/command"
)

func CallStatsService(ctx context.Context, conn *grpc.ClientConn, method string, request string) (string, error) {
	client := statsService.NewStatsServiceClient(conn)

	switch strings.ToLower(method) {
	case "getstats":
		r := &statsService.GetStatsRequest{}
		if err := proto.UnmarshalText(request, r); err != nil {
			return "", err
		}
		resp, err := client.GetStats(ctx, r)
		if err != nil {
			return "", err
		}
		return proto.MarshalTextString(resp), nil
	case "querystats":
		r := &statsService.QueryStatsRequest{}
		if err := proto.UnmarshalText(request, r); err != nil {
			return "", err
		}
		resp, err := client.QueryStats(ctx, r)
		if err != nil {
			return "", err
		}
		return proto.MarshalTextString(resp), nil
	case "getsysstats":
		// SysStatsRequest is an empty message
		r := &statsService.SysStatsRequest{}
		resp, err := client.GetSysStats(ctx, r)
		if err != nil {
			return "", err
		}
		return proto.MarshalTextString(resp), nil
	default:
		return "", errors.New("Unknown method: " + method)
	}
}

func CheckIPv4(address string) bool {
	num_string := strings.Split(address, ".")
	if len(num_string) != 4 {
		return false
	} else {
		for i := 0; i < 4; i++ {
			num, _ := strconv.Atoi(num_string[i])
			if num > 255 || num < 0 {
				return false
			}
		}
	}
	return true
}
