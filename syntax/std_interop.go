package syntax

import "regexp/syntax"

// FromStd will convert a stdlib syntax.Regexp into the forks version of
// Regexp.
//
// Note: we share Rune with the re to avoid allocations.
//
// Note: Care does need to be taken to ensure your version of go hasn't
// diverged too much. Luckily this API has only changed once in 2014,
// otherwise has been stable since 2011.
//
// Note: the performance improvements in this specific package are all done on
// Prog so you won't miss out on them.
func FromStd(re *syntax.Regexp) *Regexp {
	re2 := &Regexp{
		Op:    Op(re.Op),
		Flags: Flags(re.Flags),
		Rune:  re.Rune,
		Rune0: re.Rune0,
		Min:   re.Min,
		Max:   re.Max,
		Cap:   re.Cap,
		Name:  re.Name,
	}

	// Avoid reference to re if we can use short storage
	if len(re2.Rune) <= len(re2.Rune0) {
		src := re2.Rune
		dst := re2.Rune0[:len(src)]
		copy(dst, src)
		re2.Rune = dst
	}

	// Use short storage if we can
	if len(re.Sub) <= len(re2.Sub0) {
		re2.Sub = re2.Sub0[:len(re.Sub)]
	} else {
		re2.Sub = make([]*Regexp, len(re.Sub))
	}

	// Recurse into children
	for i, child := range re.Sub {
		re2.Sub[i] = FromStd(child)
	}

	return re2
}
