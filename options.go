package gateway

import "github.com/renevo/gateway/static"

type Option func(*Server)

func MountSite(path string) Option {
	return func(s *Server) {
		s.site = static.New(path)
	}
}
