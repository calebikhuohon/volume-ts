package service

import (
	"context"
	"volume-ts/pkg/req"
)

type Service struct {
	web req.Requests
}

func NewService() *Service {
	return &Service{
		web: req.NewRequests(),
	}
}

func (s *Service) SortFlightPaths(ctx context.Context, user string) ([]string, error) {
	paths, err := s.web.GetFlightPaths(ctx, user)
	if err != nil {
		return nil, err
	}

	var (
		StartAirport = make(map[string]struct{})
		endAirport   = make(map[string]struct{})

		result = make([]string, 0)
	)

	for _, path := range paths {
		StartAirport[path[0]] = struct{}{}
		endAirport[path[1]] = struct{}{}
	}

	for _, path := range paths {
		_, startOk := endAirport[path[0]]
		_, endOk := StartAirport[path[1]]

		if startOk {
			delete(endAirport, path[0])
		}

		if endOk {
			delete(StartAirport, path[1])
		}

	}

	for k := range StartAirport {
		result = append(result, k)
	}

	for k := range endAirport {
		result = append(result, k)
	}

	return result, nil
}
