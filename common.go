package easyregexp

import (
	"regexp"
	"sort"
	"sync"
)

const MaxCompiledRegex = 1024

var mutex sync.Mutex
var compiledRegex map[string]*EasyRegexp

func init() {
	compiledRegex = map[string]*EasyRegexp{}
}

type EasyRegexp struct {
	pattern string
	r       *regexp.Regexp
	count   uint
}

func NewEasyRegexp(pattern string) *EasyRegexp {

	er, ok := compiledRegex[pattern]
	if !ok {
		// 检查是否超过 MaxCompiledRegex，如果超过，就把调用次数最少的清理掉
		mutex.Lock()
		defer mutex.Unlock()
		if len(compiledRegex) > MaxCompiledRegex {
			ers := make([]*EasyRegexp, len(compiledRegex))
			i := 0
			for _, v := range compiledRegex {
				ers[i] = v
				i++
			}
			sort.Slice(ers, func(i int, j int) bool {
				return ers[i].count < ers[j].count
			})

			delete(compiledRegex, ers[0].pattern)
		}

		er = &EasyRegexp{
			pattern, regexp.MustCompile(pattern), 0,
		}
		compiledRegex[pattern] = er
	}
	return er
}

func Match(pattern string, i interface{}) bool {
	return NewEasyRegexp(pattern).Match(i)
}

func (er *EasyRegexp) Match(i interface{}) bool {
	er.count++

	s := toString(i)

	return er.r.MatchString(s)
}

func Catch(pattern string, i interface{}) []string {
	return NewEasyRegexp(pattern).Catch(i)
}

func (er *EasyRegexp) Catch(i interface{}) []string {
	er.count++

	s := toString(i)

	var result []string
	catchedString := er.r.FindStringSubmatch(s)

	if len(catchedString) == 0 {
		result = catchedString
	} else {
		result = catchedString[1:]
	}

	return result
}

func CatchAll(pattern string, i interface{}) [][]string {
	return NewEasyRegexp(pattern).CatchAll(i)
}

func (er *EasyRegexp) CatchAll(i interface{}) [][]string {
	er.count++

	s := toString(i)

	var result [][]string
	allCatchedString := er.r.FindAllStringSubmatch(s, -1)

	if len(allCatchedString) == 0 {
		result = allCatchedString
	} else {
		result = make([][]string, len(allCatchedString))
		for i, v := range allCatchedString {
			result[i] = v[1:]
		}
	}

	return result
}

func ReplaceAll(pattern string, i interface{}, repl string) string {
	return NewEasyRegexp(pattern).ReplaceAll(i, repl)
}

func (er *EasyRegexp) ReplaceAll(i interface{}, repl string) string {
	er.count++

	s := toString(i)

	return er.r.ReplaceAllString(s, repl)
}

func Split(pattern string, i interface{}) []string {
	return NewEasyRegexp(pattern).Split(i)
}

func (er *EasyRegexp) Split(i interface{}) []string {
	er.count++

	s := toString(i)

	return er.r.Split(s, -1)
}

func toString(i interface{}) string {
	var result string

	switch a := i.(type) {
	case string:
		result = a
	case []byte:
		result = string(a)
	default:
		panic("only support string and []byte")
	}

	return result
}
