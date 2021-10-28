package domain

type pluginDomainImpl struct{}

var (
	pluginDomain *pluginDomainImpl
)

func NewPluginDomain() *pluginDomainImpl {
	var domain = new(pluginDomainImpl)
	return domain
}

func GetPluginDomain() *pluginDomainImpl {
	if pluginDomain == nil {
		pluginDomain = NewPluginDomain()
	}
	return pluginDomain
}
