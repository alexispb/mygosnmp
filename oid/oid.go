package oid

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexispb/mygosnmp/generics"
)

// Clone returns new instance of oid.
func Clone(id []uint32) (idnew []uint32) {
	idnew = make([]uint32, len(id))
	copy(idnew, id)
	return
}

// Cat appends subids to oid and returns result
// as a new oid instance.
func Cat(id []uint32, subids ...uint32) (idnew []uint32) {
	idnew = make([]uint32, len(id)+len(subids))
	copy(idnew, id)
	copy(idnew[len(id):], subids)
	return
}

// HasPrefix tests whether oid begins with subids.
func HasPrefix(id []uint32, subids ...uint32) bool {
	if len(id) < len(subids) {
		return false
	}
	for i := 0; i < len(subids); i++ {
		if id[i] != subids[i] {
			return false
		}
	}
	return true
}

// Compare compares two oida lexicographically.
// The result is -1, 0, or +1.
func Compare(id1, id2 []uint32) int {
	len1, len2 := len(id1), len(id2)
	lmin := generics.Min(len1, len2)
	for i := 0; i < lmin; i++ {
		if id1[i] < id2[i] {
			return -1
		}
		if id1[i] > id2[i] {
			return 1
		}
	}
	if len1 < len2 {
		return -1
	}
	if len1 > len2 {
		return 1
	}
	return 0
}

// Lt tests whether id1 is less than id2.
func Lt(id1, id2 []uint32) bool {
	return Compare(id1, id2) == -1
}

// Le tests whether id1 is less or equal to id2.
func Le(id1, id2 []uint32) bool {
	cmp := Compare(id1, id2)
	return cmp == -1 || cmp == 0
}

// Eq tests whether id1 is equal to id2.
func Eq(id1, id2 []uint32) bool {
	return Compare(id1, id2) == 0
}

// Ne tests whether id1 is not equal to id2.
func Ne(id1, id2 []uint32) bool {
	return Compare(id1, id2) != 0
}

// Ge tests whether id1 is greater or equal to id2.
func Ge(id1, id2 []uint32) bool {
	cmp := Compare(id1, id2)
	return cmp == 0 || cmp == 1
}

// Gt tests whether id1 is greater than id2.
func Gt(id1, id2 []uint32) bool {
	return Compare(id1, id2) == 1
}

// Fprint appends oid string representation to the builder.
// E.g. string representation of []uint32{1,2,3} is 1.2.3.
func Fprint(sb *strings.Builder, id []uint32) {
	for i, subid := range id {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(strconv.FormatUint(uint64(subid), 10))
	}
}

// String return oid string representation.
// E.g. string representation of []uint32{1,2,3} is 1.2.3.
func String(id []uint32) string {
	var sb strings.Builder
	sb.Grow(11 * len(id))
	Fprint(&sb, id)
	return sb.String()
}

// Parse interprets a string s as oid string representation
// and returns the corresponding oid. It panics if the string
// can not be interpreted as oid string representation.
// (Note the string may start with a leading dot which is
// ignored in this case).
func Parse(s string) (id []uint32) {
	var err error
	defer func() {
		if err != nil {
			panic(err.Error())
		}
	}()

	if s = strings.TrimPrefix(s, "."); len(s) == 0 {
		return
	}
	nsubids := strings.Count(s, ".") + 1
	if nsubids > 128 {
		err = fmt.Errorf("ParseOid: oid length %d > 128", nsubids)
		return
	}

	id = make([]uint32, nsubids)
	for i, ind := 0, 0; i < nsubids; i++ {
		s = s[ind:]
		if ind = strings.IndexByte(s, '.'); ind == -1 {
			ind = len(s)
		}
		var subid uint64
		if subid, err = strconv.ParseUint(s[:ind], 10, 32); err != nil {
			return
		}
		id[i], ind = uint32(subid), ind+1
	}
	return
}
