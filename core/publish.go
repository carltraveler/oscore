package core

import (
	"fmt"
	"encoding/json"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/oscoreconfig"
	"io/ioutil"
	"os"
)

const (
	defaultOntId  = "did:ont:Ad4pjz2bqep4RhQrUAzMuZJkBC3qJ1tZuT"
	defaultAuthor = "admin"
)

// used to check.
func GetPulishFunctionList(funcListFile string) ([]*PublishAPI, error) {
	file, err := os.Open(funcListFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	list := make([]*PublishAPI, 0)
	err = json.Unmarshal(bs, &list)
	if err != nil {
		fmt.Printf("GetPulishFunctionList N.0 %s\n", err)
		return nil, err
	}

	return list, nil
}

func SendPublushRequest(list []*PublishAPI, reset bool, accessMode int32) error {
	log.InitLog(log.DebugLog, log.Stdout)
	var err error
	for _, l := range list {
		if reset {
			api, err := dao.DefOscoreApiDB.QueryApiBasicInfoByTitle(nil, l.Name)
			if err != nil && !dao.IsErrNoRows(err) {
				return err
			}

			if err == nil {
				tx, errl := dao.DefOscoreApiDB.DB.Beginx()
				if errl != nil {
					return errl
				}
				defer func() {
					if errl != nil {
						tx.Rollback()
					}
				}()

				errl = dao.DefOscoreApiDB.DeleteApiBasicInfoByApiId(nil, api.ApiId)
				if errl != nil {
					return errl
				}

				errl = tx.Commit()
				if errl != nil {
					return errl
				}
			}
		} else {
			_, err = dao.DefOscoreApiDB.QueryApiBasicInfoByTitle(nil, l.Name)
		}
		if reset || dao.IsErrNoRows(err) {
			err = PublishAPIHandleCore(l, defaultOntId, defaultAuthor, accessMode, oscoreconfig.AdminUserId)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return err
}
