// Code generated by "stringer -type=Error"; DO NOT EDIT.

package pduerror

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoError-0]
	_ = x[TooBig-1]
	_ = x[NoSuchName-2]
	_ = x[BadValue-3]
	_ = x[ReadOnly-4]
	_ = x[GenError-5]
	_ = x[NoAccess-6]
	_ = x[WrongType-7]
	_ = x[WrongLength-8]
	_ = x[WrongEncoding-9]
	_ = x[WrongValue-10]
	_ = x[NoCreation-11]
	_ = x[InconsistentValue-12]
	_ = x[ResourceUnavailable-13]
	_ = x[CommitFailed-14]
	_ = x[UndoFailed-15]
	_ = x[AuthorizationError-16]
	_ = x[NotWritable-17]
	_ = x[InconsistentName-18]
	_ = x[OpenFailed-256]
	_ = x[NotOpened-257]
	_ = x[IndexWrongType-258]
	_ = x[IndexAlreadyAllocated-259]
	_ = x[IndexNonAvailable-260]
	_ = x[IndexNotAllocated-261]
	_ = x[UnsupportedContext-262]
	_ = x[DuplicateRegistration-263]
	_ = x[UnknownRegistration-264]
	_ = x[UnknownAgentCaps-265]
	_ = x[ParseError-266]
	_ = x[RequestDenied-267]
	_ = x[ProcessingError-268]
}

const (
	_Error_name_0 = "NoErrorTooBigNoSuchNameBadValueReadOnlyGenErrorNoAccessWrongTypeWrongLengthWrongEncodingWrongValueNoCreationInconsistentValueResourceUnavailableCommitFailedUndoFailedAuthorizationErrorNotWritableInconsistentName"
	_Error_name_1 = "OpenFailedNotOpenedIndexWrongTypeIndexAlreadyAllocatedIndexNonAvailableIndexNotAllocatedUnsupportedContextDuplicateRegistrationUnknownRegistrationUnknownAgentCapsParseErrorRequestDeniedProcessingError"
)

var (
	_Error_index_0 = [...]uint8{0, 7, 13, 23, 31, 39, 47, 55, 64, 75, 88, 98, 108, 125, 144, 156, 166, 184, 195, 211}
	_Error_index_1 = [...]uint8{0, 10, 19, 33, 54, 71, 88, 106, 127, 146, 162, 172, 185, 200}
)

func (i Error) String() string {
	switch {
	case i <= 18:
		return _Error_name_0[_Error_index_0[i]:_Error_index_0[i+1]]
	case 256 <= i && i <= 268:
		i -= 256
		return _Error_name_1[_Error_index_1[i]:_Error_index_1[i+1]]
	default:
		return "Error(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}