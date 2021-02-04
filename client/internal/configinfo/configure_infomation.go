package configinfo

import "strings"

//ConfigureInformation config info
type ConfigureInformation struct {
	DataId        string
	Group         string
	ConfigureInfo string
}

//NewConfigureInformation new
func NewConfigureInformation() *ConfigureInformation {
	c := &ConfigureInformation{}
	return c
}

//String ConfigureInformation model to string
func (c *ConfigureInformation) String() string {
	builder := strings.Builder{}
	builder.WriteString("DataID: ")
	builder.WriteString(c.DataId)
	builder.WriteString(", Group: ")
	builder.WriteString(c.Group)
	builder.WriteString(", ConfigureInformation: ")
	builder.WriteString(c.ConfigureInfo)
	return builder.String()
}
