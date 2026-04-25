package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run  <file_name>.go <input_file> <output_file>")
		os.Exit(1)
	}

	content, err := readFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: could not read input file —", err)
		os.Exit(1)
	}

	// Apply all text transformations in order
	content = processText(content)

	if err := writeFile(os.Args[2], content); err != nil {
		fmt.Println("Error: could not write to output file —", err)
		os.Exit(1)
	}
}

// readFile reads the entire content of the file at the given path and returns it as a string.
func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	return string(data), err
}

func processText(content string) string {
	content = processModifiers(content)
	content = fixPunctuation(content)
	content = fixQuotes(content)
	content = fixArticles(content)

	return content
}

// writeFile writes the given content string to a file at the specified path,
// creating or overwriting it with permissions 0644.
func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0o644)
}

// hexToDecimal converts a hexadecimal string to its decimal representation.
// Returns the original string if parsing fails.
func hexToDecimal(hexStr string) string {
	num, err := strconv.ParseInt(hexStr, 16, 0)
	if err != nil {
		fmt.Println("Error: The entered Number is not an HexaDecimal")
		os.Exit(1)
	}
	return strconv.Itoa(int(num))
}

// binToDecimal converts a binary string to its decimal representation.
// Returns the original string if parsing fails.
func binToDecimal(binStr string) string {
	num, err := strconv.ParseInt(binStr, 2, 0)
	if err != nil {
		fmt.Println("Error: The entered Number is not an Binary Number")
		os.Exit(1)
	}
	return strconv.Itoa(int(num))
}

// processModifiers scans the text for inline modifier tags such as (hex), (bin),
// (up), (low), (cap), and their multi-word variants (up, N), (low, N), (cap, N).
// Each modifier transforms the word(s) immediately preceding it.
func processModifiers(text string) string {
	tokens := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case "(hex)":
			if len(result) == 0 {
				break
			}
			result[len(result)-1] = hexToDecimal(result[len(result)-1])
		case "(bin)":
			if len(result) == 0 {
				break
			}
			result[len(result)-1] = binToDecimal(result[len(result)-1])
		case "(up)":
			if len(result) == 0 {
				break
			}
			result[len(result)-1] = toUpperCase(result[len(result)-1])
		case "(low)":
			if len(result) == 0 {
				break
			}
			result[len(result)-1] = toLowerCase(result[len(result)-1])
		case "(cap)":
			if len(result) == 0 {
				break
			}
			result[len(result)-1] = capitalizeFirst(result[len(result)-1])
		case "(up,":
			count := parseWordCount(tokens[i+1], len(result))
			for n := 0; n < count; n++ {
				idx := len(result) - count + n
				if idx >= 0 {
					result[idx] = toUpperCase(result[idx])
				}
			}
			i++ // skip the number token
		case "(low,":
			count := parseWordCount(tokens[i+1], len(result))
			for n := 0; n < count; n++ {
				idx := len(result) - count + n
				if idx >= 0 {
					result[idx] = toLowerCase(result[idx])
				}
			}
			i++ // skip the number token
		case "(cap,":
			count := parseWordCount(tokens[i+1], len(result))
			for n := 0; n < count; n++ {
				idx := len(result) - count + n
				if idx >= 0 {
					result[idx] = capitalizeFirst(result[idx])
				}
			}
			i++ // skip the number token
		default:
			result = append(result, tokens[i])
		}
	}

	return strings.Join(result, " ")
}

// fixPunctuation ensures that punctuation marks (.,!?:;) are attached to the
// preceding word with no space before them, and a single space after.
func fixPunctuation(text string) string {
	// Add spaces around each punctuation mark so Fields can split on them
	var spaced strings.Builder
	for _, ch := range text {
		if strings.ContainsRune(".,!?:;", ch) {
			spaced.WriteString(" ")
			spaced.WriteRune(ch)
			spaced.WriteString(" ")
		} else {
			spaced.WriteRune(ch)
		}
	}

	tokens := strings.Fields(spaced.String())
	result := []string{}

	for i, token := range tokens {
		isPunct := token == "." || token == "," || token == "!" ||
			token == "?" || token == ":" || token == ";"

		if i == 0 || isPunct {
			result = append(result, token)
		} else {
			result = append(result, " "+token)
		}
	}

	return strings.Join(result, "")
}

// toUpperCase returns the word converted to all uppercase letters.
func toUpperCase(word string) string {
	return strings.ToUpper(word)
}

// toLowerCase returns the word converted to all lowercase letters.
func toLowerCase(word string) string {
	return strings.ToLower(word)
}

// capitalizeFirst returns the word with only its first letter uppercased
// and the rest lowercased.
func capitalizeFirst(word string) string {
	words := []rune(word)
	return strings.ToTitle(string(words[0])) + strings.ToLower(string(words[1:]))
}

// parseWordCount strips the closing parenthesis from a token like "3)" and
// returns the integer value. Falls back to 0 on error, or clamps to the
// number of available words if the requested count exceeds it.
func parseWordCount(token string, available int) int {
	trimmed := strings.TrimSuffix(token, ")")
	count, err := strconv.Atoi(trimmed)
	if err != nil {
		fmt.Println("Error: expected a number in modifier, got:", token)
		os.Exit(1)
	}
	if count > available {
		fmt.Printf("Warning: modifier count (%d) exceeds available words (%d); applying to all available words.\n", count, available)
	}
	return count
}

// fixQuotes ensures that single and double quotation marks have no space
// between them and the text they enclose:
//
//	' hello world '  →  'hello world'
func fixQuotes(text string) string {
	// Isolate each quote character so Fields can split on it
	var spaced strings.Builder
	for _, ch := range text {
		if strings.ContainsRune("'\"", ch) {
			spaced.WriteString(" ")
			spaced.WriteRune(ch)
			spaced.WriteString(" ")
		} else {
			spaced.WriteRune(ch)
		}
	}

	tokens := strings.Fields(spaced.String())
	openQuotes := []string{} // stack of currently open quote characters
	result := []string{}
	noSpaceBefore := false // true immediately after an opening quote

	for i, token := range tokens {
		isQuote := token == "'" || token == "\""

		if isQuote {
			if i == 0 {
				// Opening quote at the very start of the text
				result = append(result, token)
				openQuotes = append(openQuotes, token)
				noSpaceBefore = true
				continue
			}

			if len(openQuotes) == 0 {
				// No open quote yet → this is an opening quote
				openQuotes = append(openQuotes, token)
				result = append(result, " "+token)
				noSpaceBefore = true
			} else if token == openQuotes[len(openQuotes)-1] {
				// Matches the most recent open quote → closing quote (no leading space)
				result = append(result, token)
				openQuotes = slices.Delete(openQuotes, len(openQuotes)-1, len(openQuotes))
			} else {
				// Different quote type → nested opening quote
				openQuotes = append(openQuotes, token)
				if noSpaceBefore {
					result = append(result, token)
				} else {
					result = append(result, " "+token)
				}
				noSpaceBefore = true
			}
		} else if i == 0 {
			result = append(result, token)
		} else if noSpaceBefore {
			// Word immediately after an opening quote — no leading space
			result = append(result, token)
			noSpaceBefore = false
		} else {
			result = append(result, " "+token)
		}
	}

	return strings.Join(result, "")
}

// fixArticles corrects "a"/"A" to "an"/"An" (and vice-versa) based on whether
// the following word begins with a vowel sound (a, e, i, o, u, h).
func fixArticles(text string) string {
	words := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(words); i++ {
		// Last word: nothing follows it, so no correction needed
		if i == len(words)-1 {
			result = append(result, words[i])
			break
		}

		isArticle := words[i] == "a" || words[i] == "A" ||
			words[i] == "an" || words[i] == "An"

		if isArticle {
			nextFirstChar := []rune(words[i+1])[0]
			vowelSound := strings.ContainsRune("aeiouhAEIOUH", nextFirstChar)

			isLower := words[i] == "a" || words[i] == "an"
			if vowelSound {
				if isLower {
					result = append(result, "an")
				} else {
					result = append(result, "An")
				}
			} else {
				if isLower {
					result = append(result, "a")
				} else {
					result = append(result, "A")
				}
			}
		} else {
			result = append(result, words[i])
		}
	}

	return strings.Join(result, " ")
}
