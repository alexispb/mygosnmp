package demo

import (
	"fmt"
	"net"
	"testing"

	"github.com/alexispb/mygosnmp/agentx"
	"github.com/alexispb/mygosnmp/asn"
	"github.com/alexispb/mygosnmp/logger"
	"github.com/alexispb/mygosnmp/oid"
	"github.com/alexispb/mygosnmp/pduerror"
)

func init() { oid.Name["myVar"] = []uint32{1, 3, 6, 1, 4, 1, 999, 1, 0} }

func TestDemo(t *testing.T) {
	log := logger.Console()
	runDemoMaster(1071, log)
	chanFinished := runDemoSubagent(1071, log)
	<-chanFinished
}

// runDemoMaster (simplified for demo purpose)
// sends three requests to subagent:
// (1) valid Get-request
// (2) invalid Get-request
// (3) valid Close-request
func runDemoMaster(port int, log logger.Log) {
	ln, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	go func() {
		var (
			conn, _ = ln.Accept()
			data    []byte
			hdata   = make([]byte, agentx.PduHeaderSize)
			pdata   []byte
		)

		pdu := agentx.Pdu{
			Flags:     agentx.FlagNetworkByteOrder,
			SessionId: 7,
			PacketId:  101,
		}

		log.Write("Master.Sending valid Get-request")
		pdu.Tag = agentx.TagGet
		pdu.Params = agentx.NoParams{}
		pdu.Ranges = []agentx.SearchRange{
			{StartOid: oid.Name["myVar"]},
		}
		data, _ = agentx.EncodePdu(pdu)
		conn.Write(data)

		// reading response
		conn.Read(hdata)
		pdu, _ = agentx.DecodePduHeader(hdata)
		pdata = make([]byte, pdu.PayloadSize)
		conn.Read(pdata)
		_ = agentx.DecodePduPayload(&pdu, pdata)

		log.Write("Master.Sending invalid Get-request")
		pdu.Tag = agentx.TagGet
		pdu.Params = agentx.NoParams{}
		pdu.Ranges = []agentx.SearchRange{
			{StartOid: oid.Name["myVar"]},
		}
		data, _ = agentx.EncodePdu(pdu)
		// damage data: set reserve byte to non-zero value
		data[23] = 0xFF
		conn.Write(data)

		// reading response
		conn.Read(hdata)
		pdu, _ = agentx.DecodePduHeader(hdata)
		pdata = make([]byte, pdu.PayloadSize)
		conn.Read(pdata)
		_ = agentx.DecodePduPayload(&pdu, pdata)

		log.Write("Master.Sending Close-request")
		pdu.Tag = agentx.TagClose
		pdu.Params = agentx.CloseParams{
			Reason: agentx.CloseReasonShutdown,
		}
		data, _ = agentx.EncodePdu(pdu)
		conn.Write(data)
	}()
}

// runDemoSubagent (simplified for demo purpose)
// reads Get- and Close-requests.
func runDemoSubagent(port int, log logger.Log) (chanFinished chan struct{}) {
	chanFinished = make(chan struct{})

	go func() {
		var (
			conn, _ = net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
			ok      bool
			pdu     agentx.Pdu
			hdata   = make([]byte, agentx.PduHeaderSize)
			pdata   []byte
			data    []byte
		)
	loop:
		for {
			log.Write("Subagent. Waiting for request")
			conn.Read(hdata)
			if pdu, ok = agentx.DecodePduHeader(hdata); !ok {
				log.Write("Subagent. Error decoding header")
				log.Write("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Log Begin")
				agentx.DecodePduHeaderDbg(hdata, log)
				log.Write("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Log End")
				log.Write("Subagent. Decoding header error is considered critical")
				break loop
			}

			pdata = make([]byte, pdu.PayloadSize)
			conn.Read(pdata)
			if ok = agentx.DecodePduPayload(&pdu, pdata); !ok {
				log.Writef("Subagent. Error decoding payload")
				pdu.Tag = agentx.TagResponse
				pdu.Params = agentx.ResponseParams{
					Error: pduerror.ParseError,
				}
				log.Writef("Subagent. Sending response notifying ParseError:\n%s\n",
					pdu.String(1))
				data, _ = agentx.EncodePdu(pdu)
				conn.Write(data)
				log.Write("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Log Begin")
				pdu, _ = agentx.DecodePduHeaderDbg(hdata, log)
				_ = agentx.DecodePduPayloadDbg(&pdu, pdata, log)
				log.Write("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Log End")
				continue loop
			}

			log.Writef("Subagent. Received request:\n%s\n",
				pdu.String(1))

			switch pdu.Tag {
			case agentx.TagClose:
				chanFinished <- struct{}{}
				break loop

			case agentx.TagGet:
				pdu.Tag = agentx.TagResponse
				pdu.Params = agentx.ResponseParams{}
				pdu.Varbinds = []asn.Varbind{
					{Oid: oid.Name["myVar"], Tag: asn.TagInteger32, Value: int32(123)},
				}
				pdu.Ranges = nil
				log.Writef("Subagent. Sending response:\n%s\n",
					pdu.String(1))
				data, _ = agentx.EncodePdu(pdu)
				conn.Write(data)
			}
		}
	}()

	return
}
