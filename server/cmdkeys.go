package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &keysCmd{} })
}

type keysCmd struct {
	baseCmd
	pattern string
}

func (cmd *keysCmd) cmdName() string {
	return "keys"
}

func (cmd *keysCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.pattern)},
	}}}
}

func (cmd *keysCmd) arity() int {
	return 2
}

func (cmd *keysCmd) readonly() bool {
	return true
}

func (cmd *keysCmd) distributed() bool {
	return true
}

func (cmd *keysCmd) exec(shard *shard) (cmdResult, error) {
	items := make([]string, 0)
	for key := range shard.items {
		if match(navigablePattern{c: []rune(cmd.pattern)}, navigablePattern{c: []rune(key)}) {
			items = append(items, key)
		}
	}
	return &stringsResult{items: items}, nil
}

func (cmd *keysCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.pattern, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*keysCmd)(nil)

// this function is ported from original redis source to copy the behavior exactly
// see https://github.com/antirez/redis/blob/f97efe0cac1653fbadf02679056cba3e5317aad2/src/util.c#L47
func match(pattern navigablePattern, string navigablePattern) bool {
	for pattern.len() > 0 {
		switch pattern.c[0] {
		case '*':
			for pattern.check(1, '*') {
				pattern = pattern.next()
			}
			if pattern.len() == 1 {
				return true
			}
			for string.len() > 0 {
				if match(pattern.next(), string) {
					return true
				}
				string = string.next()
			}
			return false
		case '?':
			if string.len() == 0 {
				return false
			}
			string = string.next()
		case '[':
			var not, match bool
			pattern = pattern.next()
			not = pattern.check(0, '^')
			if not {
				pattern = pattern.next()
			}
			match = false
			for {
				if pattern.check(0, '\\') && pattern.len() >= 2 {
					pattern = pattern.next()
					if pattern.same(0, 0, string) {
						match = true
					}
				} else if pattern.check(0, ']') {
					break
				} else if pattern.len() == 0 {
					pattern = pattern.prev()
					break
				} else if pattern.check(1, '-') && pattern.len() >= 3 {
					if string.between(0, pattern.c[0], pattern.c[2]) {
						match = true
					}
					pattern = pattern.next().next()
				} else {
					if pattern.same(0, 0, string) {
						match = true
					}
				}
				pattern = pattern.next()
			}
			if not {
				match = !match
			}
			if !match {
				return false
			}
			string = string.next()
		case '\\':
			if pattern.len() >= 2 {
				pattern = pattern.next()
			}
			fallthrough
		default:
			if !pattern.same(0, 0, string) {
				return false
			}
			string = string.next()
		}
		pattern = pattern.next()
		if string.len() == 0 {
			for pattern.check(0, '*') {
				pattern = pattern.next()
			}
			break
		}
	}
	if pattern.len() == 0 && string.len() == 0 {
		return true
	}
	return false
}

type navigablePattern struct {
	c []rune
	p *navigablePattern
}

func (pattern navigablePattern) next() navigablePattern {
	return navigablePattern{pattern.c[1:], &pattern}
}

func (pattern navigablePattern) prev() navigablePattern {
	return *pattern.p
}

func (pattern navigablePattern) check(i int, r rune) bool {
	return pattern.len() > i && pattern.c[i] == r
}

func (pattern navigablePattern) between(i int, r1 rune, r2 rune) bool {
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return pattern.len() > i && pattern.c[i] >= r1 && pattern.c[i] <= r2
}

func (pattern navigablePattern) same(i, j int, other navigablePattern) bool {
	return pattern.len() > i && other.len() > j && pattern.c[i] == other.c[j]
}

func (pattern navigablePattern) len() int {
	return len(pattern.c)
}
