package agentx

/*
Package agentx provides an interface for encoding and decoding
of Agentx PDU as defined in RFC 2741. It includes three main
functions: EncodePdu, DecodePduHeader, and DecodePduPayload.

Also, it includes debug versions of this functions which are
functionally equivalent to the corresponding functions and
additionally log the process of encoding/decoding and are
intended to be used in case when the main function encounters
an error. Debug functions include additional log parameter
(see logger package), e.g.

	log := logger.Console()

The EncodePdu function encodes pdu and returns data ready for
sending across the wire:

	pdu := agentx.Pdu{...}
	data, ok := agentx.Encode(pdu)
	if !ok {
		agentx.EncodeDbg(pdu, log)
		...
	}
	conn.Write(data)

The DecodePduHeader and DecodePduPayload functions are to be
used for decoding pdu. DecodePduHeader returns pdu with header
fields set to the result of decoding header data:

	hdata := make([]byte, agentx.PduHeaderSize)
	conn.Read(hdata)
	pdu, ok := agentx.DecodePduHeader(hdata)
	if !ok {
		agentx.DecodePduHeaderDbg(hdata, log)
		...
	}

DecodePduPayload sets pdu payload fields to the result of
decoding payload data:

	pdata := make([]byte, pdu.PayloadSize)
	conn.Read(pdata)
	ok = agentx.DecodePduPayload(&pdu, pdata)
	if !ok {
		pdu, _ := agentx.DecodePduHeaderDbg(hdata, log)
		agentx.DecodePduPayloadDbg(&pdu, pdata, log)
		...
	}

Agentx PDU is represented by the unified agentx.Pdu structure
which covers different PDU types. E.g.

	pdu := agentx.Pdu{
		Tag:           agentx.TagResponse,
		Flags:         agentx.FlagNetworkByteOrder,
		SessionId:     ...,
		TransactionId: ...,
		PacketId:      ...,
		Context        ...,
		PayloadSize    ...,
		Params         agentx.ResponseParams{...}
	    Ranges         []agentx.SearchRange{...}
		Varbinds       []agentx.Varbind{...}
	}

The Params field is an interface whose implementation must corresponds
to PDU type. If PDU has no specific payload params, this field
must be set to NoParams{}.

When encoding or decoding an agentx.pdu structure the Context,
Ranges, and Varbinds fields are ignored if they are not included
in the corresponding Agentx PDU type as specified by RFC 2741.

The PayloadSize field is set in process of encoding pdu.

See also agentx/demo/demo_test.go.
*/
