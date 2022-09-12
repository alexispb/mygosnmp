package agentx

import "strconv"

type CloseReason byte

const (
	CloseReasonOther         CloseReason = 1
	CloseReasonParseError    CloseReason = 2
	CloseReasonProtocolError CloseReason = 3
	CloseReasonTimeouts      CloseReason = 4
	CloseReasonShutdown      CloseReason = 5
	CloseReasonByManager     CloseReason = 6
)

// String returns CloseReason string representation
func (r CloseReason) String() string {
	switch r {
	case CloseReasonOther:
		return "Other"
	case CloseReasonParseError:
		return "ParseError"
	case CloseReasonProtocolError:
		return "ProtocolError"
	case CloseReasonTimeouts:
		return "Timeouts"
	case CloseReasonShutdown:
		return "Shutdown"
	case CloseReasonByManager:
		return "ByManager"
	default:
		return "?" + strconv.FormatInt(int64(r), 10)
	}
}
