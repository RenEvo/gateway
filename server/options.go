package server

import "github.com/renevo/gateway/server/static"

type Option func(*Server)

func MountSite(path string) Option {
	return func(s *Server) {
		s.site = static.New(path)
	}
}
