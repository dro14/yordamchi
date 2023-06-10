package functions

import (
	"log"
	"time"
)

func LanguageCode(lang string) string {
	if lang == "" {
		lang = "uz"
	} else if lang != "uz" && lang != "ru" {
		lang = "en"
	}
	return lang
}

func Sleep(retryDelay *time.Duration) {
	log.Printf("retrying request after %v", *retryDelay)
	time.Sleep(*retryDelay)
	*retryDelay *= 2
}

//func MarkdownV2(text string) string {
//
//	var (
//		found       bool
//		before      []byte
//		buffer      bytes.Buffer
//		char        = []byte{0}
//		bts         = []byte(text)
//		escape      = []byte{'\\', 0}
//		block       = []byte{'`', '`', '`'}
//		blocks      = bytes.Count(bts, block)
//		escapeChars = []byte{'\\', '_', '*', '[', ']', '(', ')', '~', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!', '`'}
//	)
//
//	for i := range escapeChars {
//		char[0] = escapeChars[i]
//		escape[1] = escapeChars[i]
//		if i != 18 {
//			bts = bytes.ReplaceAll(bts, char, escape)
//		}
//	}
//
//	if blocks%2 != 0 {
//		bts = append(bts, block...)
//	}
//
//	for {
//		before, bts, found = bytes.Cut(bts, block)
//		if bytes.Count(before, char)%2 != 0 {
//			before = bytes.ReplaceAll(before, char, escape)
//		}
//		buffer.Write(before)
//		if !found {
//			break
//		}
//		buffer.Write(block)
//
//		before, bts, _ = bytes.Cut(bts, block)
//		before = bytes.ReplaceAll(before, char, escape)
//		buffer.Write(before)
//		buffer.Write(block)
//	}
//
//	return buffer.String()
//}
