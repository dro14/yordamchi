package functions

import "bytes"

func MarkdownV2(text string) string {

	bts := []byte(text)
	escapeChars := []byte{'\\', '_', '[', ']', '(', ')', '~', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'}

	char := []byte{' '}
	escape := []byte{'\\', ' '}
	for i := range escapeChars {
		char[0] = escapeChars[i]
		escape[1] = escapeChars[i]
		bts = bytes.ReplaceAll(bts, char, escape)
	}

	block := []byte{'`', '`', '`'}
	blocks := bytes.Count(bts, block)
	if blocks%2 != 0 {
		bts = append(bts, block...)
	}

	var (
		found  bool
		before []byte
		buffer bytes.Buffer
	)

	for {
		before, bts, found = bytes.Cut(bts, block)
		switch {
		case bytes.Count(before, []byte{'`'})%2 != 0:
			before = bytes.ReplaceAll(before, []byte{'`'}, []byte{'\\', '`'})
		case bytes.Count(before, []byte{'*', '*'})%2 != 0:
			before = bytes.ReplaceAll(before, []byte{'*'}, []byte{'\\', '*'})
		}
		buffer.Write(before)
		if !found {
			break
		}
		buffer.Write(block)

		before, bts, _ = bytes.Cut(bts, block)
		before = bytes.ReplaceAll(before, []byte{'`'}, []byte{'\\', '`'})
		before = bytes.ReplaceAll(before, []byte{'*'}, []byte{'\\', '*'})
		buffer.Write(before)
		buffer.Write(block)
	}

	return buffer.String()
}
