package virtualaccounts

import (
	"testing"
	"time"
)

type Event struct {
	Name      string
	StartDate time.Time
	EndDate   time.Time
}

func TestActiveEvents(t *testing.T) {

	// ini dibikin skenario, kalau inputnya begini, dan hasilnya oke, nanti akan keluar oke.
	res := ActiveEvents(le, time.Date(2020, 07, 03, 0, 0, 0, 0, time.UTC))
	if len(res) != 2 {
		t.Fatalf("Expect 2 Events, got: %v", len(res))
	}
	if res[0].Name != "Training REST" {
		t.Fatalf("Expected event : Training REST, got : %v", res[0].Name)
	}
	if res[1].Name != "Training Android" {
		t.Fatalf("Expected event : Training Android, got : %v", res[1].Name)
	}

}
