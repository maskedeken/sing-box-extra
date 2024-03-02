package boxapi

import (
	"context"
	"github.com/matsuridayo/sing-box-extra/boxbox"
	"github.com/sagernet/sing-box/common/dialer"
	"github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"net"
)

func DialContext(ctx context.Context, box *boxbox.Box, network, addr string) (net.Conn, error) {
	router := box.Router()
	conn, err := dialer.NewRouter(router).DialContext(ctx, network, metadata.ParseSocksaddr(addr))
	if err != nil {
		return nil, err
	}
	if vs := router.V2RayServer(); vs != nil {
		if ss, ok := vs.StatsService().(*SbStatsService); ok {
			outbound, err := router.DefaultOutbound(N.NetworkName(network))
			if err != nil {
				return nil, err
			}
			conn = ss.RoutedConnectionInternal("", outbound.Tag(), "", conn, false)
		}
	}
	return conn, nil
}

func DialUDP(ctx context.Context, box *boxbox.Box) (net.PacketConn, error) {
	router := box.Router()
	return dialer.NewRouter(router).ListenPacket(ctx, metadata.Socksaddr{})
}
