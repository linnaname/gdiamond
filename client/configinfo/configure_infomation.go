package configinfo

import "strings"

type ConfigureInfomation struct {
	DataId        string
	Group         string
	ConfigureInfo string
}

func (c *ConfigureInfomation) String() string {
	builder := strings.Builder{}
	builder.WriteString("DataId: ")
	builder.WriteString(c.DataId)
	builder.WriteString(", Group: ")
	builder.WriteString(c.Group)
	builder.WriteString(", ConfigureInfomation: ")
	builder.WriteString(c.ConfigureInfo)
	return builder.String()
}
