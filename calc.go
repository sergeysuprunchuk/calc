package calc

import (
	"context"
)

func Calc(program string) any {
	return newParser(program).
		parse().
		exec(context.Background())
}
