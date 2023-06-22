package monday

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// dateStringLayoutItem represents one word or set of delimiters between words.
// This is an abstraction level above date raw character string of date representation.
//
// Example: "1  February / 2013" ->
//           dateStringLayoutItem { item: "1",        isWord: true }
//           dateStringLayoutItem { item: "  ",       isWord: false }
//           dateStringLayoutItem { item: "February", isWord: true }
//           dateStringLayoutItem { item: " / ",      isWord: false }
//           dateStringLayoutItem { item: "2013",     isWord: true }
type dateStringLayoutItem struct {
	item    string
	isWord  bool // true if this is a sequence of letters/digits (as opposed to a sequence of non-letters like delimiters)
	isDigit bool // true if this is a sequence only containing digits
}

// extractLetterSequence extracts first word (sequence of letters ending with a non-letter)
// starting with the specified index and wraps it to dateStringLayoutItem according to the type
// of the word.
func extractLetterSequence(originalStr string, index int) (it dateStringLayoutItem) {
	letters := &strings.Builder{}

	bytesToParse := []byte(originalStr[index:])
	runeCount := utf8.RuneCount(bytesToParse)

	var isWord bool
	var isDigit bool

	letters.Grow(runeCount)
	for i := 0; i < runeCount; i++ {
		rne, runeSize := utf8.DecodeRune(bytesToParse)
		bytesToParse = bytesToParse[runeSize:]

		if i == 0 {
			isWord = unicode.IsLetter(rne)
			isDigit = unicode.IsDigit(rne)
		} else {
			if (isWord && (!unicode.IsLetter(rne) && !unicode.IsDigit(rne))) ||
				(isDigit && !unicode.IsDigit(rne)) ||
				(!isWord && unicode.IsLetter(rne)) ||
				(!isDigit && unicode.IsDigit(rne)) {
				break
			}
		}

		letters.WriteRune(rne)
	}

	it.item = letters.String()
	it.isWord = isWord
	it.isDigit = isDigit
	return
}

// stringToLayoutItems transforms raw date string (like "2 Mar 2012") into
// a set of dateStringLayoutItems, which are more convenient to work with
// in other analysis modules.
func stringToLayoutItems(dateStr string) (seqs []dateStringLayoutItem) {
	i := 0

	for i < len(dateStr) {
		seq := extractLetterSequence(dateStr, i)
		i += len(seq.item)
		seqs = append(seqs, seq)
	}

	return
}

func layoutToString(li []dateStringLayoutItem) string {
	// This function is expensive enough to be worth counting
	// bytes and allocating all in one go.
	numChars := 0
	for _, v := range li {
		numChars += len(v.item)
	}

	sb := &strings.Builder{}
	sb.Grow(numChars)
	for _, v := range li {
		sb.WriteString(v.item)
	}

	return sb.String()
}
