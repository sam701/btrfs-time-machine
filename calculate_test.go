package main

import (
	"testing"
	"time"
)

func TestCalculation(t *testing.T) {
	today = time.Date(2016, time.June, 1, 0, 0, 0, 0, time.UTC)
	config = &Config{
		Retentions: []Retention{
			Retention{
				EveryNDays:      1,
				SnapshotsToKeep: 3,
			},
			Retention{
				EveryNDays:      5,
				SnapshotsToKeep: 3,
			},
		},
	}

	toCheck := getDatesToDelete([]time.Time{
		today,
		today.AddDate(0, 0, -1),
		today.AddDate(0, 0, -2),
		today.AddDate(0, 0, -3),
		today.AddDate(0, 0, -4),
		today.AddDate(0, 0, -5),
		today.AddDate(0, 0, -6),
		today.AddDate(0, 0, -7),
		today.AddDate(0, 0, -11),
		today.AddDate(0, 0, -12),
		today.AddDate(0, 0, -16),
		today.AddDate(0, 0, -23),
		today.AddDate(0, 0, -25),
	})

	expected := []time.Time{
		today.AddDate(0, 0, -3),
		today.AddDate(0, 0, -4),
		today.AddDate(0, 0, -5),
		today.AddDate(0, 0, -6),
		today.AddDate(0, 0, -11),
		today.AddDate(0, 0, -23),
	}

	if !equalDateSlices(toCheck, expected) {
		t.Fatal("expected", expected, "but was", toCheck)
	}
}

func equalDateSlices(a, b []time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	for i, ax := range a {
		if !ax.Equal(b[i]) {
			return false
		}
	}

	return true
}
