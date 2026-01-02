package libbox

import (
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

type PlatformInterface interface {
	LocalDNSTransport() LocalDNSTransport
	UsePlatformAutoDetectInterfaceControl() bool
	AutoDetectInterfaceControl(fd int32) error
	OpenTun(options TunOptions) (int32, error)
	WriteLog(message string)
	UseProcFS() bool
	FindConnectionOwner(ipProtocol int32, sourceAddress string, sourcePort int32, destinationAddress string, destinationPort int32) (int32, error)
	PackageNameByUid(uid int32) (string, error)
	UIDByPackageName(packageName string) (int32, error)
	StartDefaultInterfaceMonitor(listener InterfaceUpdateListener) error
	CloseDefaultInterfaceMonitor(listener InterfaceUpdateListener) error
	GetInterfaces() (NetworkInterfaceIterator, error)
	UnderNetworkExtension() bool
	IncludeAllNetworks() bool
	SystemCertificates() StringIterator
	ClearDNSCache()
	SendNotification(notification *Notification) error
}

type TunInterface interface {
	FileDescriptor() int32
	Close() error
}

type InterfaceUpdateListener interface {
	UpdateDefaultInterface(interfaceName string, interfaceIndex int32, isExpensive bool, isConstrained bool)
}

const (
	InterfaceTypeCellular = int32(C.InterfaceTypeCellular)
	InterfaceTypeEthernet = int32(C.InterfaceTypeEthernet)
	InterfaceTypeOther    = int32(C.InterfaceTypeOther)
)

type NetworkInterface struct {
	Index     int32
	MTU       int32
	Name      string
	Addresses StringIterator
	Flags     int32

	Type      int32
	DNSServer StringIterator
	Metered   bool
}

type NetworkInterfaceIterator interface {
	Next() *NetworkInterface
	HasNext() bool
}

type Notification struct {
	Identifier string
	TypeName   string
	TypeID     int32
	Title      string
	Subtitle   string
	Body       string
	OpenURL    string
}

type OnDemandRule interface {
	Target() int32
	DNSSearchDomainMatch() StringIterator
	DNSServerAddressMatch() StringIterator
	SSIDMatch() StringIterator
	ProbeURL() string
}

type OnDemandRuleIterator interface {
	Next() OnDemandRule
	HasNext() bool
}

type onDemandRule struct {
	option.OnDemandRule
}

func (r *onDemandRule) Target() int32 {
	if r.OnDemandRule.Action == nil {
		return -1
	}
	return int32(*r.OnDemandRule.Action)
}

func (r *onDemandRule) DNSSearchDomainMatch() StringIterator {
	return newIterator(r.OnDemandRule.DNSSearchDomainMatch)
}

func (r *onDemandRule) DNSServerAddressMatch() StringIterator {
	return newIterator(r.OnDemandRule.DNSServerAddressMatch)
}

func (r *onDemandRule) SSIDMatch() StringIterator {
	return newIterator(r.OnDemandRule.SSIDMatch)
}

func (r *onDemandRule) ProbeURL() string {
	return r.OnDemandRule.ProbeURL
}
