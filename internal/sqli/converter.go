package sqli

import (
	"encoding/hex"
	"regexp"
	"strings"
)

type function func(string) string

var FArray = []function{
	ToLower, ConvertFromSQLHex,
	ConvertFromCommented,
	ConvertFromWhiteSpace,
	ConvertQuotes,
	convertFromOutOfRangeChars,
}

func ToLower(data string) string {
	data = strings.ToLower(data)
	return data
}

// Check for comments and erases them if available
func ConvertFromCommented(data string) string {
	if matched, _ := regexp.MatchString(`(?:\<!-|-->|\/\*|\*\/|\/\/\W*\w+\s*$)|(?:--[^-]*-)`, data); matched {
		pattern := []string{
			`(?:(?:<!)(?:(?:--(?:[^-]*(?:-[^-]+)*)--\s*)*)(?:>))`,
			`(?:(?:\/\*\/*[^\/\*]*)+\*\/)`,
			`(?:--[^-]*-)`,
		}
		for _, s := range pattern {
			re, _ := regexp.Compile(s)
			data = re.ReplaceAllString(data, `;`)
		}
	}
	re, _ := regexp.Compile(`(<\w+)\/+(\w+=?)`)
	data = re.ReplaceAllString(data, `$1/$2`)
	re, _ = regexp.Compile(`[^\\\:]\/\/(.*)$`)
	data = re.ReplaceAllString(data, `/**/$1`)
	re, _ = regexp.Compile(`([^\-&])#.*[\r\n\v\f]`)
	data = re.ReplaceAllString(data, `$1`)
	re, _ = regexp.Compile(`([^&\-])#.*\n`)
	data = re.ReplaceAllString(data, `$1`)
	re, _ = regexp.Compile(`^#.*\n`)
	data = re.ReplaceAllString(data, ` `)

	return data
}

// Strip newlines
func ConvertFromWhiteSpace(data string) string {
	// check for inline linebreaks
	search := []string{"\r", "\n", "\f", "\t", "\v"}
	for _, s := range search {
		data = strings.ReplaceAll(data, s, ";")
	}
	// replace replacement characters regular spaces
	data = strings.ReplaceAll(data, "�", " ")

	// convert real linebreaks
	re, _ := regexp.Compile(`(?:\n|\r|\v)`)
	return re.ReplaceAllString(data, ` `)
}

// Normalize quotes
func ConvertQuotes(data string) string {
	// normalize different quotes to "
	search := []string{"'", "`", "´", "’", "‘"}
	for _, s := range search {
		data = strings.ReplaceAll(data, s, "\"")
	}
	//make sure harmless quoted strings don't generate false alerts
	re, _ := regexp.Compile(`^"([^"=\\!><~]+)"$`)
	return re.ReplaceAllString(data, `$1`)
}

// Converts SQLHEX to plain text
func ConvertFromSQLHex(data string) string {
	re, _ := regexp.Compile(`(?:(?:\A|[^\d])0x[a-f\d]{3,}[a-f\d]*)+`)
	arr := re.FindAllString(data, -1)
	if len(arr) > 0 {
		for _, val := range arr {
			val = strings.Trim(val, " ")
			tempArr := strings.Split(val, " ")
			for _, tempS := range tempArr {
				converted := strings.ReplaceAll(tempS, "0x", "")
				bs, _ := hex.DecodeString(converted)
				converted = string(bs)
				data = strings.Replace(data, tempS, converted, 1)
			}
		}
	}

	// take care of hex encoded ctrl chars
	re, _ = regexp.Compile(`0x\d+`)
	return re.ReplaceAllString(data, ` 1 `)
}

// Detects out of ascii chars and change tem to space
func convertFromOutOfRangeChars(data string) string {
	values := []rune(data)
	for i, item := range values {
		if item >= 127 {
			values[i] = 32
		}
	}
	return string(values)
}

// This method removes encoded sql # comments
//func convertFromUrlencodeSqlComment(data string) string {
//	re, _ := regexp.Compile(`(?:#.*?\n)`)
//	if arr := re.FindAllString(data, -1); len(arr) > 0 {
//		converted := data
//		for _, match := range arr {
//			converted = strings.Replace(converted, match, " ", 1)
//		}
//		data += "\n" + converted
//	}
//	return data
//}
//func main() {
//	type fs func(string) string
//	data := "hi"
//	f := []fs{convertQuotes}
//	for _, v := range f {
//		data = v(data)
//	}
//}
