package primitives

type StringSet map[string]bool

func NewStringSet(s []string) StringSet {
	r := StringSet{}
	for _, i := range s {
		r.Add(i)
	}
	return r
}

func (m StringSet) Add(s string) {
	m[s] = true
}

func (m StringSet) Remove(s string) {
	delete(m, s)
}

func (m StringSet) Contains(s string) bool {
	_, ok := m[s]
	return ok
}

func (m StringSet) Strings() Strings {
	r := make(Strings, len(m))
	i := 0
	for k, _ := range m {
		r[i] = k
		i++
	}
	return r
}
