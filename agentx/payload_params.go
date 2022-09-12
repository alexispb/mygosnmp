package agentx

import (
	"github.com/alexispb/mygosnmp/pduerror"
)

type PayloadParams interface {
	// encodingSize returns number of bytes required
	// to encode PayloadParams.
	encodingSize() int
	// append appends result of encoding PayloadParams.
	append(e encoder, data []byte) []byte
	// appendDbg appends the result of encoding PayloadParams
	// and additionally logs the process of encoding
	appendDbg(e encoderDbg, data []byte) []byte
}

// NoParams is the implementation of PayloadParams applicable to
// pdu which require no payload params, i.e. Get-, GetNext-, TestSet-,
// CommitSet-, UndoSet-, CleanupSet-, Notify-, Ping-, IndexAllocate-,
// and IndexDeallocate-payloads.
type NoParams struct{}

func (NoParams) encodingSize() int {
	return 0
}
func (NoParams) append(e encoder, data []byte) []byte {
	return data
}

func parseNoParams(d decoder, data []byte) (PayloadParams, []byte) {
	return NoParams{}, data
}

func (NoParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending NoParams: nothing to append")
	return data
}

func parseNoParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing NoParams: nothing to parse")
	return NoParams{}, pos
}

// OpenParams is the implementation of PayloadParams applicable to
// Open-pdu.
type OpenParams struct {
	Timeout     byte
	Oid         []uint32
	Description string
}

func (p OpenParams) encodingSize() int {
	return 4 + objectIdEncodingSize(p.Oid) + octetStringEncodingSize(p.Description)
}

func (p OpenParams) append(e encoder, data []byte) []byte {
	data = append(data, p.Timeout, 0, 0, 0)
	data = e.appendObjectId(data, p.Oid, 0)
	return e.appendOctetString(data, p.Description)
}

func parseOpenParams(d decoder, data []byte) (PayloadParams, []byte) {
	p := OpenParams{}
	if data[1] != 0 || data[2] != 0 || data[3] != 0 {
		panic(ErrEncoding)
	}
	p.Timeout = data[0]
	p.Oid, _, data = d.parseObjectId(data[4:])
	p.Description, data = d.parseOctetString(data)
	return p, data
}

func (p OpenParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending OpenParams")
	e.log.Writef("appending Timeout: %d, and reserved bytes: 0, 0, 0",
		p.Timeout)
	data = append(data, p.Timeout, 0, 0, 0)
	e.log.Write("appending Oid")
	data = e.appendObjectId(data, p.Oid, 0)
	e.log.Write("appending Description")
	data = e.appendOctetString(data, p.Description)
	return data
}

func parseOpenParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing OpenParams")
	p := OpenParams{}

	d.log.Write("parsing first 4 bytes")
	data, pos := d.nextChunk(pos, 4)
	p.Timeout = data[0]
	d.log.Writef("parsed Timeout: %d", p.Timeout)

	if data[1] != 0 || data[2] != 0 || data[3] != 0 {
		panic("non-zero reserved byte")
	}

	d.log.Write("parsing Oid")
	p.Oid, _, pos = d.parseObjectId(pos)
	d.log.Write("parsing Description")
	p.Description, pos = d.parseOctetString(pos)
	return p, pos
}

// CloseParams is the implementation of PayloadParams applicable to
// Close-pdu.
type CloseParams struct {
	Reason CloseReason
}

func (CloseParams) encodingSize() int {
	return 4
}

func (p CloseParams) append(e encoder, data []byte) []byte {
	return append(data, byte(p.Reason), 0, 0, 0)
}

func parseCloseParams(d decoder, data []byte) (PayloadParams, []byte) {
	if data[1] != 0 || data[2] != 0 || data[3] != 0 {
		panic(ErrEncoding)
	}
	return CloseParams{Reason: CloseReason(data[0])}, data[4:]
}

func (p CloseParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending CloseParams")
	e.log.Writef("appending Reason: %s, and reserved bytes: 0, 0, 0",
		p.Reason.String())
	return append(data, byte(p.Reason), 0, 0, 0)
}

func parseCloseParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing CloseParams")
	p := CloseParams{}

	d.log.Write("parsing first 4 bytes")
	data, pos := d.nextChunk(pos, 4)
	p.Reason = CloseReason(data[0])
	d.log.Writef("parsed Reason: %d", p.Reason.String())
	if data[1] != 0 || data[2] != 0 || data[3] != 0 {
		panic("non-zero reserved byte")
	}
	return p, pos
}

// RegisterParams is the implementation of PayloadParams applicable to
// Register-pdu.
type RegisterParams struct {
	Timeout    byte
	Priority   byte
	RangeSubid byte
	Subtree    []uint32
	UpperBound uint32
}

func (p RegisterParams) encodingSize() (size int) {
	size += 4 + objectIdEncodingSize(p.Subtree)
	if p.RangeSubid != 0 {
		size += 4
	}
	return
}

func (p RegisterParams) append(e encoder, data []byte) []byte {
	data = append(data, p.Timeout, p.Priority, p.RangeSubid, 0)
	data = e.appendObjectId(data, p.Subtree, 0)
	if p.RangeSubid != 0 {
		data = e.appendUint32(data, p.UpperBound)
	}
	return data
}

func parseRegisterParams(d decoder, data []byte) (PayloadParams, []byte) {
	if data[3] != 0 {
		panic(ErrEncoding)
	}
	p := RegisterParams{Timeout: data[0], Priority: data[1], RangeSubid: data[2]}
	p.Subtree, _, data = d.parseObjectId(data[4:])
	if p.RangeSubid != 0 {
		p.UpperBound, data = d.parseUint32(data)
	}
	return p, data
}

func (p RegisterParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending RegisterParams")
	e.log.Writef("appending Timeout: %d, Priority: %d, RangeSubid: %d, reserved byte: 0",
		p.Timeout, p.Priority, p.RangeSubid)
	data = append(data, p.Timeout, p.Priority, p.RangeSubid, 0)

	e.log.Write("appending Subtree")
	data = e.appendObjectId(data, p.Subtree, 0)

	if p.RangeSubid != 0 {
		e.log.Write("appending UpperBound")
		data = e.appendUint32(data, p.UpperBound)
	}
	return data
}

func parseRegisterParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing RegisterParams")
	d.log.Write("parsing first 4 bytes")
	data, pos := d.nextChunk(pos, 4)
	p := RegisterParams{Timeout: data[0], Priority: data[1], RangeSubid: data[2]}
	d.log.Writef("parsed Timeout: %d, Priority: %d, RangeSubid: %d",
		p.Timeout, p.Priority, p.RangeSubid)
	if data[3] != 0 {
		panic("non-zero reserved byte")
	}
	d.log.Write("parsing Subtree")
	p.Subtree, _, pos = d.parseObjectId(pos)
	if p.RangeSubid != 0 {
		d.log.Write("parsing UpperBound")
		p.UpperBound, pos = d.parseUint32(pos)
	}
	return p, pos
}

// UnregisterParams is the implementation of PayloadParams applicable to
// Unregister-pdu.
type UnregisterParams struct {
	Priority   byte
	RangeSubid byte
	Subtree    []uint32
	UpperBound uint32
}

func (p UnregisterParams) encodingSize() (size int) {
	size += 4 + objectIdEncodingSize(p.Subtree)
	if p.RangeSubid != 0 {
		size += 4
	}
	return
}

func (p UnregisterParams) append(e encoder, data []byte) []byte {
	data = append(data, 0, p.Priority, p.RangeSubid, 0)
	data = e.appendObjectId(data, p.Subtree, 0)
	if p.RangeSubid != 0 {
		data = e.appendUint32(data, p.UpperBound)
	}
	return data
}

func parseUnregisterParams(d decoder, data []byte) (PayloadParams, []byte) {
	if data[0] != 0 || data[3] != 0 {
		panic(ErrEncoding)
	}
	p := UnregisterParams{Priority: data[1], RangeSubid: data[2]}
	p.Subtree, _, data = d.parseObjectId(data[4:])
	if p.RangeSubid != 0 {
		p.UpperBound, data = d.parseUint32(data)
	}
	return p, data
}

func (p UnregisterParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending UnregisterParams")
	e.log.Writef("appending reserved byte: 0, Priority: %d, RangeSubid: %d, reserved byte: 0",
		p.Priority, p.RangeSubid)
	data = append(data, 0, p.Priority, p.RangeSubid, 0)
	e.log.Write("appending Subtree")
	data = e.appendObjectId(data, p.Subtree, 0)
	if p.RangeSubid != 0 {
		e.log.Write("appending UpperBound")
		data = e.appendUint32(data, p.UpperBound)
	}
	return data
}

func parseUnregisterParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing UnregisterParams")
	d.log.Write("parsing first 4 bytes")
	data, pos := d.nextChunk(pos, 4)
	p := UnregisterParams{Priority: data[1], RangeSubid: data[2]}
	d.log.Writef("Priority: %d, RangeSubid: %d",
		p.Priority, p.RangeSubid)
	if data[0] != 0 || data[3] != 0 {
		panic("non-zero reserved byte(s)")
	}
	d.log.Write("parsing Subtree")
	p.Subtree, _, pos = d.parseObjectId(pos)
	if p.RangeSubid != 0 {
		d.log.Write("parsing UpperBound")
		p.UpperBound, pos = d.parseUint32(pos)
	}
	return p, pos
}

// GetBulkParams is the implementation of PayloadParams applicable to
// GetBulk-pdu.
type GetBulkParams struct {
	NonRepeaters    int16
	MaxRepeatitions int16
}

func (GetBulkParams) encodingSize() int {
	return 4
}

func (p GetBulkParams) append(e encoder, data []byte) []byte {
	data = e.appendInt16(data, p.NonRepeaters)
	return e.appendInt16(data, p.MaxRepeatitions)
}

func parseGetBulkParams(d decoder, data []byte) (PayloadParams, []byte) {
	p := GetBulkParams{}
	p.NonRepeaters, data = d.parseInt16(data)
	p.MaxRepeatitions, data = d.parseInt16(data)
	if p.NonRepeaters < 0 || p.MaxRepeatitions < 0 {
		panic(ErrEncoding)
	}
	return p, data
}

func (p GetBulkParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending GetBulkParams")
	e.log.Write("appending NonRepeaters")
	data = e.appendInt16(data, p.NonRepeaters)
	e.log.Write("appending MaxRepeatitions")
	return e.appendInt16(data, p.MaxRepeatitions)
}

func parseGetBulkParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing GetBulkParams")
	p := GetBulkParams{}
	d.log.Write("parsing NonRepeaters")
	p.NonRepeaters, pos = d.parseInt16(pos)
	if p.NonRepeaters < 0 {
		d.log.Write("!!! Error. NonRepeaters < 0")
		panic(ErrEncoding)
	}
	d.log.Write("parsing MaxRepeatitions")
	p.MaxRepeatitions, pos = d.parseInt16(pos)
	if p.MaxRepeatitions < 0 {
		d.log.Write("!!! Error. MaxRepeatitions < 0")
		panic(ErrEncoding)
	}
	return p, pos
}

// AddAgentCapsParams is the implementation of PayloadParams applicable to
// AddAgentCaps-pdu.
type AddAgentCapsParams struct {
	Oid         []uint32
	Description string
}

func (p AddAgentCapsParams) encodingSize() int {
	return objectIdEncodingSize(p.Oid) + octetStringEncodingSize(p.Description)
}

func (p AddAgentCapsParams) append(e encoder, data []byte) []byte {
	data = e.appendObjectId(data, p.Oid, 0)
	return e.appendOctetString(data, p.Description)
}

func parseAddAgentCapsParams(d decoder, data []byte) (PayloadParams, []byte) {
	p := AddAgentCapsParams{}
	p.Oid, _, data = d.parseObjectId(data)
	p.Description, data = d.parseOctetString(data)
	return p, data
}

func (p AddAgentCapsParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending AddAgentCapsParams")
	e.log.Write("appending Oid")
	data = e.appendObjectId(data, p.Oid, 0)
	e.log.Write("appending Description")
	return e.appendOctetString(data, p.Description)
}

func parseAddAgentCapsParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing AddAgentCapsParams")
	p := AddAgentCapsParams{}
	d.log.Write("parsing Oid")
	p.Oid, _, pos = d.parseObjectId(pos)
	d.log.Write("parsing Description")
	p.Description, pos = d.parseOctetString(pos)
	return p, pos
}

// RemoveAgentCapsParams is the implementation of PayloadParams applicable to
// RemoveAgentCaps-pdu.
type RemoveAgentCapsParams struct {
	Oid []uint32
}

func (p RemoveAgentCapsParams) encodingSize() int {
	return objectIdEncodingSize(p.Oid)
}

func (p RemoveAgentCapsParams) append(e encoder, data []byte) []byte {
	return e.appendObjectId(data, p.Oid, 0)
}

func parseRemoveAgentCapsParams(d decoder, data []byte) (PayloadParams, []byte) {
	p := RemoveAgentCapsParams{}
	p.Oid, _, data = d.parseObjectId(data)
	return p, data
}

func (p RemoveAgentCapsParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending RemoveAgentCapsParams")
	e.log.Write("appending Id")
	return e.appendObjectId(data, p.Oid, 0)
}

func parseRemoveAgentCapsParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing RemoveAgentCapsParams")
	p := RemoveAgentCapsParams{}
	d.log.Write("parsing Oid")
	p.Oid, _, pos = d.parseObjectId(pos)
	return p, pos
}

// ResponseCapsParams is the implementation of PayloadParams applicable to
// ResponseCaps-pdu.
type ResponseParams struct {
	SysUpTime uint32
	Error     pduerror.Error
	Index     int16
}

func (p ResponseParams) encodingSize() int {
	return 8
}

func (p ResponseParams) append(e encoder, data []byte) []byte {
	data = e.appendUint32(data, p.SysUpTime)
	data = e.appendInt16(data, int16(p.Error))
	return e.appendInt16(data, p.Index)
}

func parseResponseParams(d decoder, data []byte) (PayloadParams, []byte) {
	p := ResponseParams{}
	p.SysUpTime, data = d.parseUint32(data)
	err, data := d.parseInt16(data)
	p.Error = pduerror.Error(err)
	p.Index, data = d.parseInt16(data)
	return p, data
}

func (p ResponseParams) appendDbg(e encoderDbg, data []byte) []byte {
	e.log.Write("appending ResponseParams")
	e.log.Write("appending SysUpTime")
	data = e.appendUint32(data, p.SysUpTime)
	e.log.Writef("appending Error %s", p.Error.String())
	data = e.appendInt16(data, int16(p.Error))
	e.log.Write("appending Index")
	return e.appendInt16(data, p.Index)
}

func parseResponseParamsDbg(d decoderDbg, pos int) (PayloadParams, int) {
	d.log.Write("parsing ResponseParams")
	p := ResponseParams{}
	d.log.Write("parsing SysUpTime")
	p.SysUpTime, pos = d.parseUint32(pos)
	d.log.Write("parsing Error")
	pduerr, pos := d.parseInt16(pos)
	p.Error = pduerror.Error(pduerr)
	d.log.Writef("parsed Error: %s", p.Error.String())
	d.log.Write("parsing Index")
	p.Index, pos = d.parseInt16(pos)
	return p, pos
}
