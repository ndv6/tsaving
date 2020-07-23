package virtualaccounts

import (
	"testing"
)

func VacToMain() {

}

func TestVacToMain(t *testing.T) {

	// ini dibikin skenario, kalau inputnya begini, dan hasilnya oke, nanti akan keluar oke.
	res := VacToMain()
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
