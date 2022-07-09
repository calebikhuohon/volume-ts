package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
	"volume-ts/mocks"
)

func TestService_SortFlightPaths(t *testing.T) {
	rq := &mocks.Requests{}

	rq.On("GetFlightPaths", context.Background(), "user-1").Return([][]string{{"SFO", "EWR"}}, nil)
	rq.On("GetFlightPaths", context.Background(), "user-1").Return([][]string{{"ATL", "EWR"}, {"SFO", "ATL"}}, nil)
	rq.On("GetFlightPaths", context.Background(), "user-1").Return([][]string{{"IND", "EWR"}, {"SFO", "ATL"}, {"GSO", "IND"}, {"ATL", "GSO"}}, nil)

	s := Service{web: rq}

	tests := []struct {
		paths  [][]string
		output []string
	}{
		{
			paths:  [][]string{{"SFO", "EWR"}},
			output: []string{"SFO", "EWR"},
		},
		{
			paths:  [][]string{{"ATL", "EWR"}, {"SFO", "ATL"}},
			output: []string{"SFO", "EWR"},
		},
		{
			paths:  [][]string{{"IND", "EWR"}, {"SFO", "ATL"}, {"GSO", "IND"}, {"ATL", "GSO"}},
			output: []string{"SFO", "EWR"},
		},
	}

	for _, tt := range tests {
		out, err := s.SortFlightPaths(context.Background(), "user-1")
		require.NoError(t, err)
		if !slices.Equal(out, tt.output) {
			t.Errorf("want %v, got %v", tt.output, out)
		}
	}
}
