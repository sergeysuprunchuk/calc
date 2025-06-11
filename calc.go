package calc

func Calc(program string, namespace Namespace) any {
	return newParser(program).
		parse().
		exec(namespace)
}
