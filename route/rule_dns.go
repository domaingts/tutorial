package route

import (
	"net/netip"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewDNSRule(router adapter.Router, logger log.ContextLogger, options option.DNSRule, checkServer bool) (adapter.DNSRule, error) {
	switch options.Type {
	case "", C.RuleTypeDefault:
		if !options.DefaultOptions.IsValid() {
			return nil, E.New("missing conditions")
		}
		if options.DefaultOptions.Server == "" && checkServer {
			return nil, E.New("missing server field")
		}
		return NewDefaultDNSRule(router, logger, options.DefaultOptions)
	case C.RuleTypeLogical:
		if !options.LogicalOptions.IsValid() {
			return nil, E.New("missing conditions")
		}
		if options.LogicalOptions.Server == "" && checkServer {
			return nil, E.New("missing server field")
		}
		return NewLogicalDNSRule(router, logger, options.LogicalOptions)
	default:
		return nil, E.New("unknown rule type: ", options.Type)
	}
}

var _ adapter.DNSRule = (*DefaultDNSRule)(nil)

type DefaultDNSRule struct {
	abstractDefaultRule
	disableCache bool
	rewriteTTL   *uint32
	clientSubnet *netip.Prefix
}

func NewDefaultDNSRule(router adapter.Router, logger log.ContextLogger, options option.DefaultDNSRule) (*DefaultDNSRule, error) {
	rule := &DefaultDNSRule{
		abstractDefaultRule: abstractDefaultRule{
			invert:   options.Invert,
			outbound: options.Server,
		},
		disableCache: options.DisableCache,
		rewriteTTL:   options.RewriteTTL,
		clientSubnet: (*netip.Prefix)(options.ClientSubnet),
	}
	if len(options.Inbound) > 0 {
		item := NewInboundRule(options.Inbound)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if options.IPVersion > 0 {
		switch options.IPVersion {
		case 4, 6:
			item := NewIPVersionItem(options.IPVersion == 6)
			rule.items = append(rule.items, item)
			rule.allItems = append(rule.allItems, item)
		default:
			return nil, E.New("invalid ip version: ", options.IPVersion)
		}
	}
	if len(options.QueryType) > 0 {
		item := NewQueryTypeItem(options.QueryType)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Network) > 0 {
		item := NewNetworkItem(options.Network)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.AuthUser) > 0 {
		item := NewAuthUserItem(options.AuthUser)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Protocol) > 0 {
		item := NewProtocolItem(options.Protocol)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Domain) > 0 || len(options.DomainSuffix) > 0 {
		item := NewDomainItem(options.Domain, options.DomainSuffix)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.DomainKeyword) > 0 {
		item := NewDomainKeywordItem(options.DomainKeyword)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.DomainRegex) > 0 {
		item, err := NewDomainRegexItem(options.DomainRegex)
		if err != nil {
			return nil, E.Cause(err, "domain_regex")
		}
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Geosite) > 0 {
		item := NewGeositeItem(router, logger, options.Geosite)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourceGeoIP) > 0 {
		item := NewGeoIPItem(router, logger, true, options.SourceGeoIP)
		rule.sourceAddressItems = append(rule.sourceAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.GeoIP) > 0 {
		item := NewGeoIPItem(router, logger, false, options.GeoIP)
		rule.destinationIPCIDRItems = append(rule.destinationIPCIDRItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourceIPCIDR) > 0 {
		item, err := NewIPCIDRItem(true, options.SourceIPCIDR)
		if err != nil {
			return nil, E.Cause(err, "source_ip_cidr")
		}
		rule.sourceAddressItems = append(rule.sourceAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.IPCIDR) > 0 {
		item, err := NewIPCIDRItem(false, options.IPCIDR)
		if err != nil {
			return nil, E.Cause(err, "ip_cidr")
		}
		rule.destinationIPCIDRItems = append(rule.destinationIPCIDRItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if options.SourceIPIsPrivate {
		item := NewIPIsPrivateItem(true)
		rule.sourceAddressItems = append(rule.sourceAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if options.IPIsPrivate {
		item := NewIPIsPrivateItem(false)
		rule.destinationIPCIDRItems = append(rule.destinationIPCIDRItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourcePort) > 0 {
		item := NewPortItem(true, options.SourcePort)
		rule.sourcePortItems = append(rule.sourcePortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourcePortRange) > 0 {
		item, err := NewPortRangeItem(true, options.SourcePortRange)
		if err != nil {
			return nil, E.Cause(err, "source_port_range")
		}
		rule.sourcePortItems = append(rule.sourcePortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Port) > 0 {
		item := NewPortItem(false, options.Port)
		rule.destinationPortItems = append(rule.destinationPortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.PortRange) > 0 {
		item, err := NewPortRangeItem(false, options.PortRange)
		if err != nil {
			return nil, E.Cause(err, "port_range")
		}
		rule.destinationPortItems = append(rule.destinationPortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessName) > 0 {
		item := NewProcessItem(options.ProcessName)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessPath) > 0 {
		item := NewProcessPathItem(options.ProcessPath)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessPathRegex) > 0 {
		item, err := NewProcessPathRegexItem(options.ProcessPathRegex)
		if err != nil {
			return nil, E.Cause(err, "process_path_regex")
		}
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.PackageName) > 0 {
		item := NewPackageNameItem(options.PackageName)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.User) > 0 {
		item := NewUserItem(options.User)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.UserID) > 0 {
		item := NewUserIDItem(options.UserID)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Outbound) > 0 {
		item := NewOutboundRule(options.Outbound)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if options.ClashMode != "" {
		item := NewClashModeItem(router, options.ClashMode)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.WIFISSID) > 0 {
		item := NewWIFISSIDItem(router, options.WIFISSID)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.WIFIBSSID) > 0 {
		item := NewWIFIBSSIDItem(router, options.WIFIBSSID)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.RuleSet) > 0 {
		item := NewRuleSetItem(router, options.RuleSet, options.RuleSetIPCIDRMatchSource, options.RuleSetIPCIDRAcceptEmpty)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	return rule, nil
}

func (r *DefaultDNSRule) DisableCache() bool {
	return r.disableCache
}

func (r *DefaultDNSRule) RewriteTTL() *uint32 {
	return r.rewriteTTL
}

func (r *DefaultDNSRule) ClientSubnet() *netip.Prefix {
	return r.clientSubnet
}

func (r *DefaultDNSRule) WithAddressLimit() bool {
	if len(r.destinationIPCIDRItems) > 0 {
		return true
	}
	for _, rawRule := range r.items {
		ruleSet, isRuleSet := rawRule.(*RuleSetItem)
		if !isRuleSet {
			continue
		}
		if ruleSet.ContainsDestinationIPCIDRRule() {
			return true
		}
	}
	return false
}

func (r *DefaultDNSRule) Match(metadata *adapter.InboundContext) bool {
	metadata.IgnoreDestinationIPCIDRMatch = true
	defer func() {
		metadata.IgnoreDestinationIPCIDRMatch = false
	}()
	return r.abstractDefaultRule.Match(metadata)
}

func (r *DefaultDNSRule) MatchAddressLimit(metadata *adapter.InboundContext) bool {
	return r.abstractDefaultRule.Match(metadata)
}

var _ adapter.DNSRule = (*LogicalDNSRule)(nil)

type LogicalDNSRule struct {
	abstractLogicalRule
	disableCache bool
	rewriteTTL   *uint32
	clientSubnet *netip.Prefix
}

func NewLogicalDNSRule(router adapter.Router, logger log.ContextLogger, options option.LogicalDNSRule) (*LogicalDNSRule, error) {
	r := &LogicalDNSRule{
		abstractLogicalRule: abstractLogicalRule{
			rules:    make([]adapter.HeadlessRule, len(options.Rules)),
			invert:   options.Invert,
			outbound: options.Server,
		},
		disableCache: options.DisableCache,
		rewriteTTL:   options.RewriteTTL,
		clientSubnet: (*netip.Prefix)(options.ClientSubnet),
	}
	switch options.Mode {
	case C.LogicalTypeAnd:
		r.mode = C.LogicalTypeAnd
	case C.LogicalTypeOr:
		r.mode = C.LogicalTypeOr
	default:
		return nil, E.New("unknown logical mode: ", options.Mode)
	}
	for i, subRule := range options.Rules {
		rule, err := NewDNSRule(router, logger, subRule, false)
		if err != nil {
			return nil, E.Cause(err, "sub rule[", i, "]")
		}
		r.rules[i] = rule
	}
	return r, nil
}

func (r *LogicalDNSRule) DisableCache() bool {
	return r.disableCache
}

func (r *LogicalDNSRule) RewriteTTL() *uint32 {
	return r.rewriteTTL
}

func (r *LogicalDNSRule) ClientSubnet() *netip.Prefix {
	return r.clientSubnet
}

func (r *LogicalDNSRule) WithAddressLimit() bool {
	for _, rawRule := range r.rules {
		switch rule := rawRule.(type) {
		case *DefaultDNSRule:
			if rule.WithAddressLimit() {
				return true
			}
		case *LogicalDNSRule:
			if rule.WithAddressLimit() {
				return true
			}
		}
	}
	return false
}

func (r *LogicalDNSRule) Match(metadata *adapter.InboundContext) bool {
	if r.mode == C.LogicalTypeAnd {
		return common.All(r.rules, func(it adapter.HeadlessRule) bool {
			metadata.ResetRuleCache()
			return it.(adapter.DNSRule).Match(metadata)
		}) != r.invert
	} else {
		return common.Any(r.rules, func(it adapter.HeadlessRule) bool {
			metadata.ResetRuleCache()
			return it.(adapter.DNSRule).Match(metadata)
		}) != r.invert
	}
}

func (r *LogicalDNSRule) MatchAddressLimit(metadata *adapter.InboundContext) bool {
	if r.mode == C.LogicalTypeAnd {
		return common.All(r.rules, func(it adapter.HeadlessRule) bool {
			metadata.ResetRuleCache()
			return it.(adapter.DNSRule).MatchAddressLimit(metadata)
		}) != r.invert
	} else {
		return common.Any(r.rules, func(it adapter.HeadlessRule) bool {
			metadata.ResetRuleCache()
			return it.(adapter.DNSRule).MatchAddressLimit(metadata)
		}) != r.invert
	}
}
