package main

import (
	"testing"
	"time"
)

func TestHumanData(t *testing.T) {
	//// Initialize a new time.Time object and pass it to the humanDate function.
	//t.Helper()
	//tm := time.Date(2022,12,15, 10, 0, 0, 0,  time.UTC)
	//hd := humanData(tm)
	//
	//// Check that the output from the humanDate function is in the format we
	//// expect. If it isn't what we expect, use the t.Errorf() function to
	//// indicate that the test has failed and log the expected and actual
	//// values.
	//if hd != "15 Dec 2022 at 10:00" {
	//	t.Errorf("want %q; got %q", "15 Dec 2022 at 10:00", hd)
	//}

	t.Helper()
	tests := []struct{
		name string
		tm time.Time
		want string
	} {
		{
			name: "UTC",
			tm: time.Date(2022, 12, 17, 10 ,0, 0 , 0, time.UTC),
			want: "17 Dec 2022 at 10:00",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm: time.Date(2020,12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanData(tt.tm)
			if hd != tt.want {
				t.Errorf("want %v, got %v", tt.want, hd)
			}
		})
	}
}
