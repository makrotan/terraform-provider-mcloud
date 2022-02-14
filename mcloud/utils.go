package mcloud

import "log"

func debug(string2 string) {
	log.Printf("\033[35m" + string2 + "\033[0m")
}
