package configinfo

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	c := &ConfigureInformation{Group: "DEFAULT_GROUP", DataId: "linname", ConfigureInfo: getConfigInfo("linname", "DEFAULT_GROUP", "I hate linnana too")}
	assert.NotNil(t, c)
	println(c.ConfigureInfo)
	println(c.String())
}

func getConfigInfo(dataId, group, content string) string {
	builder := strings.Builder{}
	builder.WriteString("dataId=")
	builder.WriteString(dataId)
	builder.WriteString(" ,group=")
	builder.WriteString(group)
	builder.WriteString(" ,content=")
	builder.WriteString(content)
	return builder.String()
}
