package internal

const (
	tabsize  = 4
	maxdepth = 8
)

var maxprefix [tabsize * maxdepth]byte

func init() {
	for i := range maxprefix {
		maxprefix[i] = ' '
	}
}

func Prefix(depth int) string {
	if depth >= maxdepth {
		depth = maxdepth
	}
	return SliceAsString(maxprefix[:tabsize*depth])
}
