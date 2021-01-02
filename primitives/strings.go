package primitives

import (
	"strings"
)

type Strings []string

func (s Strings) Any(f func(string) bool) bool {
	for _, i := range s {
		if f(i) {
			return true
		}
	}
	return false
}

func (s Strings) Contains(v string) bool {
	return s.Any(func(i string) bool {
		return i == v
	})
}

func (s Strings) Map(f func(string) string) Strings {
	r := make(Strings, len(s))
	for idx, i := range s {
		r[idx] = f(i)
	}
	return r
}

func (s Strings) Join(delimiter string) string {
	return strings.Join(s, delimiter)
}

func (s Strings) Select(f func(string) bool) Strings {
	r := Strings{}
	for _, i := range s {
		if f(i) {
			r = append(r, i)
		}
	}
	return r
}

func (s Strings) Uniq() Strings {
	m := map[string]bool{}
	return s.Select(func(i string) bool {
		if _, ok := m[i]; !ok {
			m[i] = true
			return true
		}
		return false
	})
}
