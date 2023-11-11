package functions

import (
	"log"
	"time"
)

func Sleep(retryDelay *time.Duration) {
	if *retryDelay > 0 {
		log.Printf("retrying request after %v", *retryDelay)
		time.Sleep(*retryDelay)
		*retryDelay *= 2
	}
}

func Slice(completion string) []string {

	var completions []string

	for len(completion) > 4096 {
		cutIndex := 0
	Loop:
		for i := 4096; i >= 0; i-- {
			switch completion[i] {
			case ' ', '\n', '\t', '\r':
				cutIndex = i
				break Loop
			}
		}
		completions = append(completions, completion[:cutIndex])
		completion = completion[cutIndex:]
	}

	return append(completions, completion)
}
