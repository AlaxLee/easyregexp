package easyregexp

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"sync"
)

//max number of cached *EasyRegexp
const MaxCompiledRegex = 1024

var mutex sync.Mutex
var compiledRegex map[string]*EasyRegexp

func init() {
	compiledRegex = map[string]*EasyRegexp{}
}

//contain pattern string and count
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

func FMatch(pattern string, filepath string) bool {
	return NewEasyRegexp(pattern).FMatch(filepath)
}

func (er *EasyRegexp) FMatch(filepath string) bool {
	return er.Match(openfile(filepath))
}

func Match(pattern string, i interface{}) bool {
	return NewEasyRegexp(pattern).Match(i)
}

func (er *EasyRegexp) Match(i interface{}) bool {
	er.count++

	s := toString(i)

	return er.r.MatchString(s)
}

func FCatch(pattern string, filepath string) []string {
	return NewEasyRegexp(pattern).FCatch(filepath)
}

func (er *EasyRegexp) FCatch(filepath string) []string {
	return er.Catch(openfile(filepath))
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

func FCatchAll(pattern string, filepath string) [][]string {
	return NewEasyRegexp(pattern).FCatchAll(filepath)
}

func (er *EasyRegexp) FCatchAll(filepath string) [][]string {
	return er.CatchAll(openfile(filepath))
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

func FReplaceAll(pattern string, filepath string, repl string) string {
	return NewEasyRegexp(pattern).FReplaceAll(filepath, repl)
}

func (er *EasyRegexp) FReplaceAll(filepath string, repl string) string {
	return er.ReplaceAll(openfile(filepath), repl)
}

func ReplaceAll(pattern string, i interface{}, repl string) string {
	return NewEasyRegexp(pattern).ReplaceAll(i, repl)
}

func (er *EasyRegexp) ReplaceAll(i interface{}, repl string) string {
	er.count++

	s := toString(i)

	return er.r.ReplaceAllString(s, repl)
}

func FSplit(pattern string, filepath string) []string {
	return NewEasyRegexp(pattern).FSplit(filepath)
}

func (er *EasyRegexp) FSplit(filepath string) []string {
	return er.Split(openfile(filepath))
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
	case io.Reader:
		tmp, err := ioutil.ReadAll(a)
		if err != nil {
			panic(err)
		}
		result = string(tmp)
	default:
		panic("only support string/[]byte/io.Reader")
	}

	return result
}

func openfile(filepath string) io.Reader {
	r, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	return r
}
