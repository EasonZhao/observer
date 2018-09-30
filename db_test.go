package main

import (
	"testing"
)

func TestGetLastHeight(t *testing.T) {
	const path = "./database"
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	height, err := db.GetLastHeight()
	if err != nil {
		t.Fatal(err)
	}
	if height != 0 {
		t.Fatalf("height(%d) must 0", height)
	}
}
