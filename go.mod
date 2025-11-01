module github.com/sagernet/sing-box

go 1.25

replace github.com/sagernet/sing-vmess => github.com/domaingts/sing-vmess v1.0.3

require (
	github.com/caddyserver/certmagic v0.23.0
	github.com/domaingts/electricity v1.1.0
	github.com/domaingts/venc v0.0.1
	github.com/go-chi/chi/v5 v5.2.2
	github.com/go-chi/render v1.0.3
	github.com/godbus/dbus/v5 v5.1.1-0.20230522191255-76236955d466
	github.com/gofrs/uuid/v5 v5.3.2
	github.com/insomniacslk/dhcp v0.0.0-20250417080101-5f8cf70e8c5f
	github.com/libdns/alidns v1.0.5-libdns.v1.beta1
	github.com/libdns/cloudflare v0.2.2-0.20250708034226-c574dccb31a6
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/metacubex/tfo-go v0.0.0-20250921095601-b102db4216c0
	github.com/metacubex/utls v1.8.3
	github.com/mholt/acmez/v3 v3.1.2
	github.com/miekg/dns v1.1.67
	github.com/oschwald/maxminddb-golang v1.13.1
	github.com/sagernet/asc-go v0.0.0-20241217030726-d563060fe4e1
	github.com/sagernet/bbolt v0.0.0-20231014093535-ea5cb2fe9f0a
	github.com/sagernet/cors v1.2.1
	github.com/sagernet/fswatch v0.1.1
	github.com/sagernet/gomobile v0.1.8
	github.com/sagernet/quic-go v0.52.0-sing-box-mod.3
	github.com/sagernet/sing v0.7.13
	github.com/sagernet/sing-mux v0.3.3
	github.com/sagernet/sing-quic v0.5.2-0.20250909083218-00a55617c0fb
	github.com/sagernet/sing-shadowsocks v0.2.8
	github.com/sagernet/sing-shadowsocks2 v0.2.1
	github.com/sagernet/sing-tun v0.7.3
	github.com/sagernet/sing-vmess v0.2.7
	github.com/sagernet/ws v0.0.0-20231204124109-acfe8907c854
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	github.com/vishvananda/netns v0.0.5
	go.uber.org/zap v1.27.0
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba
	golang.org/x/crypto v0.42.0
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6
	golang.org/x/mod v0.27.0
	golang.org/x/net v0.43.0
	golang.org/x/sys v0.36.0
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20241231184526-a9ab2273dd10
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
	howett.net/plist v1.0.1
)

//replace github.com/sagernet/sing => ../sing

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/caddyserver/zerossl v0.1.3 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/libdns/libdns v1.1.0 // indirect
	github.com/mdlayher/netlink v1.7.3-0.20250113171957-fbb4dce95f42 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/refraction-networking/utls v1.8.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sagernet/gvisor v0.0.0-20250325023245-7a9c0f5725fb // indirect
	github.com/sagernet/netlink v0.0.0-20240612041022-b9a21c07ac6a // indirect
	github.com/sagernet/nftables v0.3.0-beta.4 // indirect
	github.com/sagernet/smux v1.5.34-mod.2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
)
