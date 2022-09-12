package asn

type Tag byte

const (
	TagInteger32      Tag = 2
	TagOctetString    Tag = 4
	TagNull           Tag = 5
	TagObjectId       Tag = 6
	TagSequence       Tag = 48
	TagIpAddress      Tag = 64
	TagCounter32      Tag = 65
	TagGauge32        Tag = 66
	TagTimeTicks      Tag = 67
	TagOpaque         Tag = 68
	TagCounter64      Tag = 70
	TagNoSuchObject   Tag = 128
	TagNoSuchInstance Tag = 129
	TagEndOfMibView   Tag = 130
)
