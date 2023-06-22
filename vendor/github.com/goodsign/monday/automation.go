package monday

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

var (
	wordsRx        = regexp.MustCompile("(\\p{L}+)")
	debugLayoutDef = false
)

// An InvalidTypeError indicates that data was parsed incorrectly as a result
// of a type mismatch.
type InvalidTypeError struct {
	error
}

// An InvalidLengthError is returned when an item's length was longer or
// shorter than expected for a particular token.
type InvalidLengthError struct {
	error
}

// NewInvalidTypeError instantiates an InvalidTypeError.
func NewInvalidTypeError() InvalidTypeError {
	return InvalidTypeError{error: errors.New("invalid type for token")}
}

// NewInvalidLengthError instantiates an InvalidLengthError.
func NewInvalidLengthError() InvalidLengthError {
	return InvalidLengthError{error: errors.New("invalid length for token")}
}

type layoutSpanI interface {
	scanInt(s *scanner.Scanner) (int, error)
	scanString(s *scanner.Scanner) (string, error)
	isString() bool
	isDelimiter() bool
}

type lengthLimitSpan struct {
	minLength int
	maxLength int
}

func (lls lengthLimitSpan) scanInt(s *scanner.Scanner) (int, error) {
	return -1, NewInvalidTypeError()
}

func (lls lengthLimitSpan) scanString(s *scanner.Scanner) (string, error) {
	return "", NewInvalidTypeError()
}

func (lls lengthLimitSpan) isString() bool    { return false }
func (lls lengthLimitSpan) isDelimiter() bool { return false }

func initLengthLimitSpan(min, max int) lengthLimitSpan {
	return lengthLimitSpan{
		minLength: min,
		maxLength: max,
	}
}

type limitedStringSpan struct {
	lengthLimitSpan
}

func initLimitedStringSpan(minLength, maxLength int) limitedStringSpan {
	return limitedStringSpan{lengthLimitSpan: initLengthLimitSpan(minLength, maxLength)}
}

func (lss limitedStringSpan) scanString(s *scanner.Scanner) (string, error) {
	tok := s.Scan()
	if tok != scanner.EOF && tok == -2 {
		return s.TokenText(), nil
	}
	return "", NewInvalidTypeError()
}

func (lss limitedStringSpan) isString() bool { return true }
func (lss limitedStringSpan) String() string {
	return fmt.Sprintf("[limitedStringSpan:%v]", lss.lengthLimitSpan)
}

type rangeIntSpan struct {
	lengthLimitSpan
	min int
	max int
}

func initRangeIntSpan(minValue, maxValue, minLength, maxLength int) rangeIntSpan {
	return rangeIntSpan{
		lengthLimitSpan: initLengthLimitSpan(minLength, maxLength),
		min:             minValue,
		max:             maxValue,
	}
}

func (rs rangeIntSpan) scanInt(s *scanner.Scanner) (int, error) {
	var tok = s.Scan()
	var negative bool
	if tok == 45 {
		negative = true
		if debugLayoutDef {
			fmt.Printf("scan negative:'%s'\n", s.TokenText())
		}
		tok = s.Scan()
	} else if tok == 43 { // positive
		tok = s.Scan()
	}
	if tok == -3 {
		str := s.TokenText()
		i, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		if negative {
			i = i * -1
		}
		return i, nil
	}

	if debugLayoutDef {
		fmt.Printf("invalid tok: %v '%s'\n", tok, s.TokenText())
	}

	return 0, NewInvalidTypeError()
}

func (rs rangeIntSpan) String() string {
	return fmt.Sprintf("[rangeIntSpan:%v]", rs.lengthLimitSpan)
}

type delimiterSpan struct {
	lengthLimitSpan
	character string
}

func initDelimiterSpan(character string, minLength, maxLength int) delimiterSpan {
	return delimiterSpan{
		lengthLimitSpan: initLengthLimitSpan(minLength, maxLength),
		character:       character,
	}
}

func (ds delimiterSpan) scanString(s *scanner.Scanner) (string, error) {
	tok := s.Scan()
	if tok != scanner.EOF && tok != -2 && tok != 45 && tok != -3 {
		return s.TokenText(), nil
	}
	if debugLayoutDef {
		fmt.Printf("expected tok:=!(-2,-3,45), received:%d ('%s')\n", tok, s.TokenText())
	}

	return "", NewInvalidTypeError()
}

func (ds delimiterSpan) isString() bool    { return false }
func (ds delimiterSpan) isDelimiter() bool { return true }
func (ds delimiterSpan) String() string {
	return fmt.Sprintf("[delimiterSpan '%s':%v]", ds.character, ds.lengthLimitSpan)
}

type layoutDef struct {
	spans         []layoutSpanI
	errorPosition int
}

func (ld *layoutDef) validate(value string) bool {
	s := &scanner.Scanner{}
	s.Init(strings.NewReader(value))
	s.Whitespace = 0
	for _, span := range ld.spans {
		if span.isString() || span.isDelimiter() {
			if _, err := span.scanString(s); err != nil {
				ld.errorPosition = s.Pos().Offset
				if debugLayoutDef {
					fmt.Printf("error at pos: %d: %s (span=%+v) - expected string or delimiter\n", s.Pos().Offset, err.Error(), span)
				}
				return false
			}
		} else if _, err := span.scanInt(s); err != nil {
			if debugLayoutDef {
				fmt.Printf("error at pos: %d: %s (span=%+v) - expected integer\n", s.Pos().Offset, err.Error(), span)
			}
			ld.errorPosition = s.Pos().Offset
			return false
		}
	}
	ld.errorPosition = s.Pos().Offset
	return s.Pos().Offset == len(value)
}

// A LocaleDetector parses time.Time values by using various heuristics and
// techniques to determine which locale should be used to parse the
// time.Time value. As not all possible locales and formats are supported,
// this process can be somewhat lossy and inaccurate.
type LocaleDetector struct {
	localeMap         map[string]*set
	lastLocale        Locale
	layoutsMap        map[string]layoutDef
	lastErrorPosition int
}

func (ld *LocaleDetector) prepareLayout(layout string) layoutDef {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(layout))
	s.Whitespace = 0
	result := make([]layoutSpanI, 0)
	var tok rune
	// var pos int = 0
	var span layoutSpanI
	var sign bool
	//	var neg bool = false
	for tok != scanner.EOF {
		tok = s.Scan()
		switch tok {
		case -2: // text
			span = initLimitedStringSpan(1, -1)
		case -3: // digit
			span = initRangeIntSpan(-1, -1, 1, -1)
			if sign {
				sign = false
			}
		case 45: // negative sign
			sign = true
			// neg = s.TokenText() == "-"
			continue
		case 43: // positive sign
			sign = true
			continue
		case scanner.EOF:
			continue
		default: // fixed character
			span = initDelimiterSpan(s.TokenText(), 1, 1)
		}
		result = append(result, span)
		// length := s.Pos().Offset - pos
		// pos = s.Pos().Offset
		// fmt.Printf("tok'%s' [%d %d] length=%d\n", s.TokenText(), pos, s.Pos().Offset, length)

	}
	if debugLayoutDef {
		fmt.Printf("layout:'%s'\n", layout)
		fmt.Printf("layout:%v\n", result)
	}
	ret := layoutDef{spans: result}
	ld.layoutsMap[layout] = ret
	return ret
}

func (ld *LocaleDetector) validateValue(layout string, value string) bool {
	l, ok := ld.layoutsMap[layout]
	if !ok {
		l = ld.prepareLayout(layout)
	}
	result := l.validate(value)
	ld.lastErrorPosition = l.errorPosition
	return result
}

func (ld *LocaleDetector) errorPosition() int { return ld.lastErrorPosition }

func (ld *LocaleDetector) addWords(words []string, v Locale) {
	for _, w := range words {
		l := strings.ToLower(w)
		if _, ok := ld.localeMap[w]; !ok {
			ld.localeMap[w] = newSet(v)
			if l != w {
				ld.localeMap[l] = newSet(v)
			}
		} else {
			ld.localeMap[w].Add(v)
			if l != w {
				ld.localeMap[l].Add(v)
			}
		}
	}
}

// NewLocaleDetector instances a LocaleDetector instance.
func NewLocaleDetector() *LocaleDetector {
	ld := &LocaleDetector{localeMap: make(map[string]*set), lastLocale: LocaleEnGB, layoutsMap: make(map[string]layoutDef)}
	for _, v := range ListLocales() {
		days := GetShortDays(v)
		ld.addWords(days, v)
		days = GetLongDays(v)
		ld.addWords(days, v)
		months := GetShortMonths(v)
		ld.addWords(months, v)
		months = GetLongMonths(v)
		ld.addWords(months, v)
	}
	return ld
}

// Parse will attempt to parse a time.Time struct from a layout (format) and a
// value to parse from.
//
// If no locale can be determined, this method will return an error and an
// empty time object.
func (ld *LocaleDetector) Parse(layout, value string) (time.Time, error) {
	if ld.validateValue(layout, value) {
		ld.lastLocale = ld.detectLocale(value)
		return ParseInLocation(layout, value, time.UTC, ld.lastLocale)
	}
	return time.Time{}, &time.ParseError{
		Value:   value,
		Layout:  layout,
		Message: fmt.Sprintf("'%s' not matches to '%s' last error position = %d\n", value, layout, ld.lastErrorPosition),
	}
}

func (ld *LocaleDetector) detectLocale(value string) Locale {
	localesMap := make(map[Locale]int)
	for _, v := range wordsRx.FindAllStringSubmatchIndex(value, -1) {
		word := strings.ToLower(value[v[0]:v[1]])

		if localesSet, ok := ld.localeMap[word]; ok {
			localesSet.Each(func(loc Locale) bool {
				if _, ok := localesMap[loc]; !ok {
					localesMap[loc] = 1
				} else {
					localesMap[loc]++
				}
				return true
			})
		}
	}
	var result Locale = LocaleEnUS
	frequency := 0
	for key, counter := range localesMap {
		if counter > frequency {
			frequency = counter
			result = key
		}
	}
	return result
}
