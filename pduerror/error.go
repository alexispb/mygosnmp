package pduerror

//go:generate stringer -type=Error
type Error uint16

const (
	NoError             Error = 0
	TooBig              Error = 1
	NoSuchName          Error = 2
	BadValue            Error = 3
	ReadOnly            Error = 4
	GenError            Error = 5
	NoAccess            Error = 6
	WrongType           Error = 7
	WrongLength         Error = 8
	WrongEncoding       Error = 9
	WrongValue          Error = 10
	NoCreation          Error = 11
	InconsistentValue   Error = 12
	ResourceUnavailable Error = 13
	CommitFailed        Error = 14
	UndoFailed          Error = 15
	AuthorizationError  Error = 16
	NotWritable         Error = 17
	InconsistentName    Error = 18
	// Specific Agentx errors
	OpenFailed            Error = 256
	NotOpened             Error = 257
	IndexWrongType        Error = 258
	IndexAlreadyAllocated Error = 259
	IndexNonAvailable     Error = 260
	IndexNotAllocated     Error = 261
	UnsupportedContext    Error = 262
	DuplicateRegistration Error = 263
	UnknownRegistration   Error = 264
	UnknownAgentCaps      Error = 265
	ParseError            Error = 266
	RequestDenied         Error = 267
	ProcessingError       Error = 268
)
