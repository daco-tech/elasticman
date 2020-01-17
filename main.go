package main

import (
	"elasticman/cmd"
	"log"
)

func main() {

	log.SetFlags(0)
	log.Println("   ____ __           __   _       __  ___          ")
	log.Println("  / __// /___ _ ___ / /_ (_)____ /  |/  /___ _ ___ ")
	log.Println(" / _/ / // _ `/(_-</ __// // __// /|_/ // _ \\`/ _ \\")
	log.Println("/___//_/ \\_,_//___/\\__//_/ \\__//_/  /_/ \\_,_//_//_/")
	log.Println("")
	log.SetFlags(1)
	cmd.Execute()

}
