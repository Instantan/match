# match
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/Instantan/match)
[![Actions Status](https://github.com/Instantan/match/workflows/Tests/badge.svg)](https://github.com/Instantan/match/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Instantan/match)](https://goreportcard.com/report/github.com/Instantan/match)

> a matching language for strings with focus on performance

it tries to achieve these performance goals by generating a cartesian product for every "or" combination in the query and extracting the static pre and suffixes from it.

The pre and suffix checks are simple string comparisons that can reduce the amount of times the wildmatch has to run.

It also automatically detects if a match is advanced or simple and optimizes its matching based on that.

The generated matchers get sorted based on the pattern and pre/suffix complexity wich can reduce the amount of checks.

All that optimization happen during the pattern compilation, to improve the performance of the matching itself. 

The focus of that library is not on the compile + match time but rather on the match time itself. That means you shouldnt compile your query all the time but rather store the return matcher and reuse it.

## Example 

```go
import "github.com/Instantan/match"

func main() {
    m, err := match.Compile("namespace.[ real | virtual ].[ root* ].value")
	if err != nil {
        panic(err)
    }
    if m.Matches("namespace.real.root.path.value") {
        println("It's a match!")
    }
}
```
