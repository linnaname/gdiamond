package configinfo

import "strings"

type ConfigureInfomation struct {
	DataId        string
	Group         string
	ConfigureInfo string
}

func NewConfigureInfomation() *ConfigureInfomation {
	c := &ConfigureInfomation{}
	return c
}

func (c *ConfigureInfomation) String() string {
	builder := strings.Builder{}
	builder.WriteString("DataID: ")
	builder.WriteString(c.DataId)
	builder.WriteString(", Group: ")
	builder.WriteString(c.Group)
	builder.WriteString(", ConfigureInfomation: ")
	builder.WriteString(c.ConfigureInfo)
	return builder.String()
}
