package main

import "os"

func main() {
	gameDir, err := GetDefaultGameDir()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = Reconcile(gameDir, wd, false)
	if err != nil {
		panic(err)
	}
	err = Reconcile(wd, gameDir, false)
	if err != nil {
		panic(err)
	}
}
