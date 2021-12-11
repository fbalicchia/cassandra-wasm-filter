package decoder

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/datastax/go-cassandra-native-protocol/frame"
	"github.com/datastax/go-cassandra-native-protocol/message"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

func CassandraParseRequest(data []byte, codec frame.RawCodec, callBack map[string]func(count uint64)) error {

	frame, err := codec.DecodeFrame(bytes.NewReader(data))

	if err != nil {
		return nil
	}

	switch messageIn := frame.Body.Message.(type) {
	case *message.Query:
		proxywasm.LogInfof("Keyspace %s query %s", messageIn.Options.Keyspace, messageIn.Query)
		statement := queryStatement(messageIn.Query)
		_, ok := callBack[statement]
		if ok {
			callBack[statement](1)
		}

	case *message.Batch:
		proxywasm.LogInfof("Keyspace %s Batch %s", messageIn.Keyspace, messageIn.String())

	default:
		fmt.Println("value null")
	}

	return nil
}

//Find a more robust ways
func queryStatement(query string) string {
	querySlitted := strings.Split(query, " ")
	result := strings.ReplaceAll(querySlitted[0], " ", "")
	if result == "" {
		return "OTHER"
	}
	return strings.ToUpper(result)
}
