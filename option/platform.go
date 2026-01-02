package option

import (
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badoption"
)

type OnDemandOptions struct {
	Enabled bool           `json:"enabled,omitempty"`
	Rules   []OnDemandRule `json:"rules,omitempty"`
}

type OnDemandRule struct {
	Action                *OnDemandRuleAction        `json:"action,omitempty"`
	DNSSearchDomainMatch  badoption.Listable[string] `json:"dns_search_domain_match,omitempty"`
	DNSServerAddressMatch badoption.Listable[string] `json:"dns_server_address_match,omitempty"`
	SSIDMatch             badoption.Listable[string] `json:"ssid_match,omitempty"`
	ProbeURL              string                     `json:"probe_url,omitempty"`
}

type OnDemandRuleAction int

func (r *OnDemandRuleAction) MarshalJSON() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	value := *r
	var actionName string
	switch value {
	case 1:
		actionName = "connect"
	case 2:
		actionName = "disconnect"
	case 3:
		actionName = "evaluate_connection"
	default:
		return nil, E.New("unknown action: ", value)
	}
	return json.Marshal(actionName)
}

func (r *OnDemandRuleAction) UnmarshalJSON(bytes []byte) error {
	var actionName string
	if err := json.Unmarshal(bytes, &actionName); err != nil {
		return err
	}
	var actionValue int
	switch actionName {
	case "connect":
		actionValue = 1
	case "disconnect":
		actionValue = 2
	case "evaluate_connection":
		actionValue = 3
	case "ignore":
		actionValue = 4
	default:
		return E.New("unknown action name: ", actionName)
	}
	*r = OnDemandRuleAction(actionValue)
	return nil
}
