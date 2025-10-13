package main

import "testing"

func TestManually(t *testing.T) {
	src := "C:\\Users\\joghourt\\Documents\\My games\\Tabletop Simulator\\Mods"
	dst := "C:\\Users\\joghourt\\home\\dest"

	err := Reconcile(src, dst, false)
	if err != nil {
		t.Fatal(err)
	}
}
