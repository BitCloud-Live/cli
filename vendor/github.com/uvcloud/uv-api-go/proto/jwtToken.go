package uv

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

const KEY_TOKEN = "x-uv-token"

// oauthAccess supplies PerRPCCredentials from a given token.
type jwtToken struct {
	GetToken func() string
}

//NewJwtAccess constructs the PerRPCCredentials using a given token.
func NewJwtAccess(getToken func() string) credentials.PerRPCCredentials {
	return jwtToken{GetToken: getToken}
}

func (oa jwtToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		KEY_TOKEN: oa.GetToken(),
	}, nil
}

func (oa jwtToken) RequireTransportSecurity() bool {
	return true
}
