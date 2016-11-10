package game

import "testing"

func TestFlavors(t *testing.T) {
	ctx := createContext("eng")

	styles := Flavors(ctx)
	AssertEquals(t, 4, len(styles))
}
