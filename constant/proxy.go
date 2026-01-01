package constant

const (
	TypeTun          = "tun"
	TypeRedirect     = "redirect"
	TypeTProxy       = "tproxy"
	TypeDirect       = "direct"
	TypeBlock        = "block"
	TypeDNS          = "dns"
	TypeSOCKS        = "socks"
	TypeMixed        = "mixed"
	TypeSSH          = "ssh"
	TypeVLESS        = "vless"
)

const (
	TypeSelector = "selector"
	TypeURLTest  = "urltest"
)

func ProxyDisplayName(proxyType string) string {
	switch proxyType {
	case TypeTun:
		return "TUN"
	case TypeRedirect:
		return "Redirect"
	case TypeTProxy:
		return "TProxy"
	case TypeDirect:
		return "Direct"
	case TypeBlock:
		return "Block"
	case TypeDNS:
		return "DNS"
	case TypeSOCKS:
		return "SOCKS"
	case TypeMixed:
		return "Mixed"
	case TypeSSH:
		return "SSH"
	case TypeVLESS:
		return "VLESS"
	case TypeSelector:
		return "Selector"
	case TypeURLTest:
		return "URLTest"
	default:
		return "Unknown"
	}
}
