package main

import (
	"github.com/go-playground/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name, want string
		tm         time.Time
	}{
		{name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "2024-03-17 10:15 AM"},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "2024-03-17 10:15 AM",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, humanDateTime(tt.tm), tt.want)
	}

}
