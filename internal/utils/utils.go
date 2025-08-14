package utils

import "log"

func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
