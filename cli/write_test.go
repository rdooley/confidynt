package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rdooley/confidynt/service"
	"github.com/rdooley/confidynt/types"
)

func TestWrite(t *testing.T) {
	t.Log("Testing Write")
	path := "test.conf"
	mockCtrl := gomock.NewController(t)
	defer func() {
		os.Remove(path)
		mockCtrl.Finish()
	}()

	table := "table"
	key := "key"
	val := "val"

	conf := types.Config{}
	conf[key] = val
	conf["other_key"] = "other_val"
	conf["multi_key"] = "\"multi\n #linething\""

	text := "key=val\n"
	text += "# a comment\n"
	text += "other_key=other_val\n"
	text += "multi_key=\"multi\n #linething\"\n"
	ioutil.WriteFile(path, []byte(text), 0644)

	buf := new(bytes.Buffer)

	mockDynamo := service.NewMockDynamo(mockCtrl)
	mockDynamo.EXPECT().Write(table, conf).Return(nil)

	Write(table, path, mockDynamo, buf)
	assert.Equal(t, buf.String(), "test.conf written to table\n")
}
