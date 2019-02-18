[![Go Report Card](https://goreportcard.com/badge/github.com/AlaxLee/easyregexp)](https://goreportcard.com/report/github.com/AlaxLee/easyregexp)


# easyregexp
以前用perl使用正则表达式习惯固化了，打算封装一下标准包，提供一些能快速处理文本的方法

支持以下方法：
- Match
- Catch
- CatchAll
- ReplaceAll
- Split

输入参数可以是  string / \[\]byte / io.Reader