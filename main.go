package main

import (
	"encoding/binary"

	"cassandra-filter/decoder"

	"github.com/datastax/go-cassandra-native-protocol/frame"
	"github.com/datastax/go-cassandra-native-protocol/primitive"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{
		statements_insert: proxywasm.DefineCounterMetric("cassandra_filter.statements_insert"),
		statements_delete: proxywasm.DefineCounterMetric("cassandra_filter.statements_delete"),
		statements_update: proxywasm.DefineCounterMetric("cassandra_filter.statements_update"),
		statements_select: proxywasm.DefineCounterMetric("cassandra_filter.statements_select"),
		statements_other:  proxywasm.DefineCounterMetric("cassandra_filter.statements_other"),
	}
}

type pluginContext struct {
	types.DefaultPluginContext
	statements_insert proxywasm.MetricCounter
	statements_delete proxywasm.MetricCounter
	statements_update proxywasm.MetricCounter
	statements_select proxywasm.MetricCounter
	statements_other  proxywasm.MetricCounter
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewTcpContext(contextID uint32) types.TcpContext {

	callbackCounter := map[string]func(count uint64){
		"INSERT": func(count uint64) {
			ctx.statements_insert.Increment(count)
		},
		"DELETE": func(count uint64) {
			ctx.statements_delete.Increment(count)
		},
		"UPDATE": func(count uint64) {
			ctx.statements_update.Increment(count)
		},
		"SELECT": func(count uint64) {
			ctx.statements_select.Increment(count)
		},
		"OTHER": func(count uint64) {
			ctx.statements_other.Increment(count)
		},
	}
	return &networkContext{
		callback_map: callbackCounter,
		codec:        frame.NewRawCodec(),
	}
}

type networkContext struct {
	types.DefaultTcpContext
	callback_map map[string]func(count uint64)
	codec        frame.RawCodec
}

// Override types.DefaultTcpContext.
func (ctx *networkContext) OnNewConnection() types.Action {
	return types.ActionContinue
}

// Override types.DefaultTcpContext.
func (ctx *networkContext) OnDownstreamData(dataSize int, endOfStream bool) types.Action {
	if dataSize == 0 {
		return types.ActionContinue
	}

	data, err := proxywasm.GetDownstreamData(0, dataSize)
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("failed to get downstream data: %v", err)
		return types.ActionContinue
	}
	requestLen := binary.BigEndian.Uint32(data[5:9])
	dataMissing := (primitive.ProtocolVersion3.FrameHeaderLengthInBytes() + int(requestLen)) - len(data)

	if requestLen == 0 {
		return types.ActionContinue

	}

	if dataMissing > 0 {
		// full header received, but only partial request
		proxywasm.LogInfof("Hdr received, but need %d more bytes of request", dataMissing)
		return types.ActionContinue
	}

	errResult := decoder.CassandraParseRequest(data, ctx.codec, ctx.callback_map)

	if errResult != nil {
		proxywasm.LogInfof("Error", errResult.Error())
	}
	return types.ActionContinue
}

// Override types.DefaultTcpContext.
func (ctx *networkContext) OnDownstreamClose(types.PeerType) {
	proxywasm.LogInfo("downstream connection close!")
	return
}

// Override types.DefaultTcpContext.
func (ctx *networkContext) OnUpstreamData(dataSize int, endOfStream bool) types.Action {
	return types.ActionContinue
}

// Override types.DefaultTcpContext.
func (ctx *networkContext) OnStreamDone() {
	proxywasm.LogInfo("connection complete!")
}
