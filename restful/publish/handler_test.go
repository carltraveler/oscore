package publish

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/restful/api/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublishAPI_OUPUTJSON(t *testing.T) {
	rp := &common.ReqParam{}

	var x core.PublishAPI
	var Tags tables.Tag
	var ErrorCodes core.PublishErrorCode
	var Params tables.RequestParam
	var Specs tables.Specifications
	x.Tags = append(x.Tags, Tags)
	x.ErrorCodes = append(x.ErrorCodes, ErrorCodes)
	x.Params = append(x.Params, Params)
	x.Specs = append(x.Specs, Specs)
	rp.Params = &x

	s, err := json.Marshal(rp)
	assert.Nil(t, err)
	fmt.Printf("%s\n", s)
}
