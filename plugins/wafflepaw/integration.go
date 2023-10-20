package wafflepaw

import "github.com/infinitybotlist/sysmanage-web/plugins/wafflepaw/scope"

type Integration struct {
	// Display name of the integration
	Name string

	// Description of the integration
	Run func(s *scope.Scope) error
}

var IntegrationList = make(map[string]*Integration)

func AddIntegration(i *Integration) bool {
	_, ok := IntegrationList[i.Name]

	if ok {
		return false
	}

	IntegrationList[i.Name] = i
	return true
}
