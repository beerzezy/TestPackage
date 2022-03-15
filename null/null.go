// Package null provides a way to handle null without using pointers for
// all of the golang types plus Dec2, Dec5 and Dec8
package null

import (
	"bytes"
)

// Nuller supports identifying values as nullable
type Nuller interface {
	Null() bool
}

// Zeroer supports identifying values as the go zero value
type Zeroer interface {
	Zero() bool
}

// Constants that are used by StripNullJSON
const minJSONSizeWithNull = 7

var (
	// null is a string of bytes containing null in UTF8
	null = []byte(":null")
	// quote contains a quote in a byte value in UTF8
	quote = []byte("\"")[0]
	// comma contains a comma in a byte value in UTF8
	comma = []byte(",")[0]
	// comma contains a comma in a byte value in UTF8
	leftCurly = []byte("{")[0]
	// minimum size of json with a null = "":null

	// backslash contains a backslash in a byte value in UTF8
	backslash = []byte(`\`)[0]
)

// removeIndex is a structure that holds the byte slice start and end index of
// bytes that should be removed from a target byte slice.  It is used by
// StripNullJSON
type removeIndex struct {
	Start int
	End   int
}

// StripNullJSON scans through the passed in byte slice containing json
// and strips out null values and returns a byte slice.  The byte slice returned
// is the same byte slice if the passed in one is null or empty or if there were
// no instances of null.  Otherwise, it is a new slice.
func StripNullJSON(b []byte) []byte {
	if nil == b || len(b) <= minJSONSizeWithNull {
		return b
	}
	// Hold all of the removals
	rems := make([]removeIndex, 0, 20)
	// find all of the nulls

	for i := 0; i < len(b)-len(null); i++ {
		window := b[i : i+len(null)]

		// We found a null
		if bytes.Compare(window, null) == 0 {
			// start next window past the null
			quoteCount := 0
			// Now search backwards for the quoted identifier
			for j := i; j >= 0; j-- {
				// eat two quotes
				if b[j] == quote {
					quoteCount = quoteCount + 1
					if quoteCount == 2 {
						start := j
						end := i + len(null) - 1

						// If there was a backslash before the quote escaping it, because we are inside a json object embedded in a
						// string property, eat the backslash as leaving it means we'd have malformed json
						if j-1 >= 0 && b[j-1] == backslash {
							start = j - 1
						}

						// If there is a comma before us, eat it
						if j-1 >= 0 && b[j-1] == comma {
							start = j - 1
						}
						// If there is a { before us, eat the comma after us if it is there
						if j-1 >= 0 && b[j-1] == leftCurly {
							if end+1 < len(b) && b[end+1] == comma {
								end = end + 1
							}
						}
						ri := removeIndex{Start: start, End: end}
						rems = append(rems, ri)
						break
					}
				}
			}
			i = i + len(null)
		}
	}
	// We did not find any
	if len(rems) == 0 {
		return b
	}
	// Build up the stripped down json with nulls removed by looping through
	// the remove indexes and take all the data in between
	res := make([]byte, 0, len(b))
	index := 0
	for _, rem := range rems {
		// Check to see if the removals are contiguous and if so, then just
		// ignore adding anything in between them
		if index < rem.Start {
			// remove {, combinations
			if len(res) > 0 && res[len(res)-1] == leftCurly && b[index] == comma {
				index = index + 1
			}
			res = append(res, b[index:rem.Start]...)
		}
		index = rem.End + 1
	}
	res = append(res, b[index:len(b)]...)
	// final hack to remove {, if there are any.  It will loop until all are gone
	upper := len(res) - 3
	i := 0
	found := false
	for i < upper {
		res, found = stripLBracketComma(res, i, upper)
		if found {
			i = i + 1
		} else {
			break
		}
	}
	return res
}

// stripLBracketComma looks for combinations of }, that could be there due
// to some hole in the logic.  This was found in testing.  This function will
// return another slice so that it can then be operated on again if need be.
func stripLBracketComma(b []byte, index int, upper int) ([]byte, bool) {
	var res []byte
	found := false
	for i := index; i < upper; i++ {
		if b[i] == leftCurly && b[i+1] == comma {
			res = append(b[:i+1], b[i+2:]...)
			found = true
			break
		}
	}
	if !found {
		res = b
	}
	return res, found
}
