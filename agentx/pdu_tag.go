package agentx

import "strconv"

type PduTag byte

const (
	pduTagMin                 = TagOpen
	TagOpen            PduTag = 1
	TagClose           PduTag = 2
	TagRegister        PduTag = 3
	TagUnregister      PduTag = 4
	TagGet             PduTag = 5
	TagGetNext         PduTag = 6
	TagGetBulk         PduTag = 7
	TagTestSet         PduTag = 8
	TagCommitSet       PduTag = 9
	TagUndoSet         PduTag = 10
	TagCleanupSet      PduTag = 11
	TagNotify          PduTag = 12
	TagPing            PduTag = 13
	TagIndexAllocate   PduTag = 14
	TagIndexDeallocate PduTag = 15
	TagAddAgentCaps    PduTag = 16
	TagRemoveAgentCaps PduTag = 17
	TagResponse        PduTag = 18
	pduTagMax                 = TagResponse
)

func (tag PduTag) IsKnown() bool {
	return pduTagMin <= tag && tag <= pduTagMax
}

func (tag PduTag) String() string {
	if !tag.IsKnown() {
		return "?" + strconv.FormatInt(int64(tag), 10)
	}
	return pduTable[tag].tagString
}

func (tag PduTag) NotAllowedFlags(f Flags) Flags {
	return f &^ pduTable[tag].allowedFlags
}

type funcParseParams func(d decoder, data []byte) (p PayloadParams, next []byte)
type funcParseParamsDbg func(d decoderDbg, startpos int) (p PayloadParams, nextpos int)

type pduEntry struct {
	// tagString is the string representation of pdu tag.
	tagString string
	// allowedFlags are flags allowed for the pdu.
	allowedFlags Flags
	// includesRanges defines whether the pdu can include search ranges.
	includesRanges bool
	// includesVarbinds defines whether the pdu can include varbinds.
	includesVarbinds bool
	// isApplicableParams defines whether the PayloadParams
	// argument is applicable to the pdu type.
	isApplicableParams func(PayloadParams) bool
	// parseParams is the parse-function applicable for the pdu type
	parseParams funcParseParams
	// parseParamsDbg is the parseDbg-function applicable for the pdu type
	parseParamsDbg funcParseParamsDbg
}

var pduTable = [pduTagMax + 1]pduEntry{
	{},
	{
		tagString:          "Open",
		allowedFlags:       FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(OpenParams); return ok },
		parseParams:        parseOpenParams,
		parseParamsDbg:     parseOpenParamsDbg,
	},
	{
		tagString:          "Close",
		allowedFlags:       FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(CloseParams); return ok },
		parseParams:        parseCloseParams,
		parseParamsDbg:     parseCloseParamsDbg,
	},
	{
		tagString:          "Register",
		allowedFlags:       FlagInstanceRegistration | FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(RegisterParams); return ok },
		parseParams:        parseRegisterParams,
		parseParamsDbg:     parseRegisterParamsDbg,
	},
	{
		tagString:          "Unregister",
		allowedFlags:       FlagInstanceRegistration | FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(UnregisterParams); return ok },
		parseParams:        parseUnregisterParams,
		parseParamsDbg:     parseUnregisterParamsDbg,
	},
	{
		tagString:          "Get",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     true,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "GetNext",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     true,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "GetBulk",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     true,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(GetBulkParams); return ok },
		parseParams:        parseGetBulkParams,
		parseParamsDbg:     parseGetBulkParamsDbg,
	},
	{
		tagString:          "TestSet",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   true,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "CommitSet",
		allowedFlags:       FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "UndoSet",
		allowedFlags:       FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "CleanupSet",
		allowedFlags:       FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "Notify",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   true,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "Ping",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "IndexAllocate",
		allowedFlags:       FlagNewIndex | FlagAnyIndex | FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   true,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "IndexDeallocate",
		allowedFlags:       FlagNewIndex | FlagAnyIndex | FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   true,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(NoParams); return ok },
		parseParams:        parseNoParams,
		parseParamsDbg:     parseNoParamsDbg,
	},
	{
		tagString:          "AddAgentCaps",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(AddAgentCapsParams); return ok },
		parseParams:        parseAddAgentCapsParams,
		parseParamsDbg:     parseAddAgentCapsParamsDbg,
	},
	{
		tagString:          "RemoveAgentCaps",
		allowedFlags:       FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   false,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(RemoveAgentCapsParams); return ok },
		parseParams:        parseRemoveAgentCapsParams,
		parseParamsDbg:     parseRemoveAgentCapsParamsDbg,
	},
	{
		tagString:          "Response",
		allowedFlags:       FlagInstanceRegistration | FlagNewIndex | FlagAnyIndex | FlagNonDefaultContext | FlagNetworkByteOrder,
		includesRanges:     false,
		includesVarbinds:   true,
		isApplicableParams: func(p PayloadParams) bool { _, ok := p.(ResponseParams); return ok },
		parseParams:        parseResponseParams,
		parseParamsDbg:     parseResponseParamsDbg,
	},
}
