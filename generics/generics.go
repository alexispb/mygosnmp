package generics

func Min[T Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Copied from: https://cs.opensource.google/go/x/exp/+/39d4317d:constraints/constraints.go
// The tilde (~) token is used in the form ~T to denote
// the set of types whose underlying type is T.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type Integer interface {
	Signed | Unsigned
}
type Float interface {
	~float32 | ~float64
}
type Complex interface {
	~complex64 | ~complex128
}
type Ordered interface {
	Integer | Float | ~string
}
