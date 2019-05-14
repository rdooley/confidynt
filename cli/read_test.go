package cli

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rdooley/confidynt/service"
	"github.com/rdooley/confidynt/types"
)

func TestRead(t *testing.T) {
	t.Log("Testing Read")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	table := "table"
	key := "key"
	val := "val"

	conf := types.Config{}
	conf[key] = val
	conf["other_key"] = "other_val"

	expected := "key=val\n"
	expected += "other_key=other_val\n"

	buf := new(bytes.Buffer)

	mockDynamo := service.NewMockDynamo(mockCtrl)
	mockDynamo.EXPECT().Read(table, key, val).Return(conf, nil)
	Read(table, key, val, mockDynamo, buf)
	assert.Equal(t, buf.String(), expected)
}
