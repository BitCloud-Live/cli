package yb

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

const KEY_TOKEN = "x-yb-token"
const KEY_CLIENT_VERSION = "x-yb-client-version"

// oauthAccess supplies PerRPCCredentials from a given token.
type perRPCCredentials struct {
	extras           map[string]func() string
	GetToken         func() string
	GetClientVersion func() string
}

//NewPerRPC constructs the PerRPCCredentials using a given token.
func NewPerRPC(getToken func() string, getClientVersion func() string, extras map[string]func() string) credentials.PerRPCCredentials {
	return perRPCCredentials{GetToken: getToken, GetClientVersion: getClientVersion, extras: extras}
}

func (oa perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (md map[string]string, err error) {
	md = map[string]string{
		KEY_TOKEN:          oa.GetToken(),
		KEY_CLIENT_VERSION: oa.GetClientVersion(),
	}
	if oa.extras != nil {
		for k, f := range oa.extras {
			md[k] = f()
		}
	}
	return

}

func (oa perRPCCredentials) RequireTransportSecurity() bool {
	return true
}
