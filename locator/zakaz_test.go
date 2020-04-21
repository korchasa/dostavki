package locator

import (
	"github.com/korchasa/kogdaeda/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_parseTimeString(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Kiev")
	tests := []struct {
		name      string
		dt        string
		str       string
		wantStart time.Time
		wantEnd   time.Time
		wantErr   error
	}{
		{"positive", "2020-04-22", "20:00 - 22:00", time.Date(2020, 4, 22, 20, 0, 0, 0, loc), time.Date(2020, 4, 22, 22, 0, 0, 0, loc), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, gotErr := parseTimeString(tt.dt, tt.str)
			assert.Equal(t, tt.wantStart, gotStart, "bad start time")
			assert.Equal(t, tt.wantEnd, gotEnd, "bad end time")
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func Test_fold(t *testing.T) {
	windows := []types.StoreWindow{
		{
			Start: time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC),
			Store: types.Store{Chain: "foo"},
		},
		{
			Start: time.Date(2020, 1, 2, 1, 1, 1, 0, time.UTC),
			Store: types.Store{Chain: "bar"},
		},
		{
			Start: time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC),
			Store: types.Store{Chain: "bar"},
		},
	}
	folded := fold(windows)
	assert.Equal(t, [][]types.Schedule{
		{{
			Start:  time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC),
			Chains: map[string]bool{"foo":true, "bar":true},
		}},
		{{
			Start:  time.Date(2020, 1, 2, 1, 1, 1, 0, time.UTC),
			Chains: map[string]bool{"bar":true},
		}},
	}, folded)
}