package main

import (
	"log"
	"os/exec"
	scraper "reshifr/sc-scraper/pkg"
)

func main() {
	reload := 0
	onSucceed := func() {
		reload = 0
	}
	onFailed := func() {
		if reload == 0 || reload == 31 {
			exec.Command("bash", "-c", "service tor reload").Run()
			log.Println("TOR service reloading...")
		}
		reload++
	}
	scraper.Scrap(uint64(9223372036854775808), onSucceed, onFailed)
}
