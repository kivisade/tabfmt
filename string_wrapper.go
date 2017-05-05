package tabfmt

import (
	"unicode/utf8"
	"regexp"
	"strings"
)

func RuneSplit(source string, position int) (left string, right string) {
	if position > utf8.RuneCountInString(source) {
		return source, ""
	}
	return string([]rune(source)[:position]), string([]rune(source)[position:])
}

func StringBreak(source string, width int, r *regexp.Regexp, lineBreak string) (out []string) {
	var (
		buf string
		idx []int
	)

	out = []string{}

	for {
		if utf8.RuneCountInString(source) <= width {
			out = append(out, source)
			break
		}
		buf, source = RuneSplit(source, width)

		if brx := strings.Index(buf, lineBreak); brx >= 0 {
			out, source = append(out, buf[:brx]), buf[brx + len(lineBreak):]+source
			continue
		}

		if found := r.FindAllStringIndex(buf, -1); len(found) == 0 {
			out = append(out, buf)
			continue
		} else {
			idx = found[len(found)-1]
		}

		out, source = append(out, buf[:idx[0]]), buf[idx[1]:]+source
	}

	return
}

func StringBreakSimple(source string, width int) (out []string) {
	r, _ := regexp.Compile(`\s`)
	return StringBreak(source, width, r, "\n")
}

func StringWrap(source string, width int, r *regexp.Regexp, lineBreak string) (out string, lines int) {
	splitLines := StringBreak(source, width, r, lineBreak)
	return strings.Join(splitLines, lineBreak), len(splitLines)
}

func StringWrapSimple(source string, width int) (out string, lines int) {
	r, _ := regexp.Compile(`\s`)
	return StringWrap(source, width, r, "\n")
}
