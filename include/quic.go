//go:build with_quic

package include

import (
	"github.com/sagernet/sing-box/adapter/inbound"
	"github.com/sagernet/sing-box/adapter/outbound"
	"github.com/sagernet/sing-box/protocol/hysteria2"
	_ "github.com/sagernet/sing-box/protocol/naive/quic"
	"github.com/sagernet/sing-box/protocol/tuic"
	_ "github.com/sagernet/sing-box/transport/v2rayquic"
	_ "github.com/sagernet/sing-dns/quic"
)

func registerQUICInbounds(registry *inbound.Registry) {
	tuic.RegisterInbound(registry)
	hysteria2.RegisterInbound(registry)
}

func registerQUICOutbounds(registry *outbound.Registry) {
	tuic.RegisterOutbound(registry)
	hysteria2.RegisterOutbound(registry)
}
