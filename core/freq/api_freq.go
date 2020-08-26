package freq

import (
	"fmt"
	ontSdk "github.com/ontio/ontology-go-sdk"
	ontcom "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/ontio/oscore/utils"
	"sync"
	"sync/atomic"
	"time"
)

type DBCache struct {
	UpdateFreq     chan string
	UpdateReqLimit chan string
	apiKeyCache    *sync.Map //apikey -> ApiKey
	freqLock       *sync.Mutex
	CallBackLock   *sync.Mutex
	apiFreqCache   *sync.Map //ApiID ->uint64
	Layer2Sdk      *ontSdk.OntologySdk
}

func NewDBCache() *DBCache {
	res := &DBCache{
		apiKeyCache:    new(sync.Map),
		apiFreqCache:   new(sync.Map),
		freqLock:       new(sync.Mutex),
		CallBackLock:   new(sync.Mutex),
		UpdateFreq:     make(chan string, 20),
		UpdateReqLimit: make(chan string, 20),
		Layer2Sdk:      oscoreconfig.DefOscoreConfig.Layer2Sdk,
	}

	go res.UpdateFreqDataBase()
	go res.UpdateReqLimitDataBase()
	return res
}

func (this *DBCache) UpdateFreqDataBase() {
	for {
		select {
		case apiKey := <-this.UpdateFreq:
			keyIn, ok := this.apiKeyCache.Load(apiKey)
			if !ok {
				log.Debugf("apikey cache not exist")
				continue
			}

			key := keyIn.(*tables.APIKey)
			apiId := key.ApiId

			apiCounterP, ok := this.apiFreqCache.Load(apiId)
			if !ok {
				log.Debugf("apicounter cache not exist")
				continue
			}

			counter := atomic.LoadUint64(apiCounterP.(*uint64))
			log.Debugf("UpdateFreqDataBase: apiId: %d, counterP: %v,counter: %d, usedNum:%d, apiKey: %s", apiId, apiCounterP, counter, key.UsedNum, apiKey)

			err := this.updateApiKeyInvokeFre(key, counter)
			if err != nil {
				log.Errorf("UpdateFreqDataBase n.1 %s", err)
			}

			oldTime := key.Layer2Time
			if oldTime == 0 {
				oldTime = time.Now().Unix()
			}

			outTime := time.Unix(oldTime, 0).Add(time.Second * time.Duration(oscoreconfig.DefOscoreConfig.Layer2RecordInterval)).Unix()
			if false && (time.Now().Unix() >= outTime || key.Layer2Time == 0) {
				// here do sendtx
				addr, err := ontcom.AddressFromHexString(oscoreconfig.DefOscoreConfig.Layer2Contract)
				if err != nil {
					log.Errorf("UpdateFreqDataBase N.N.0 %s", err)
					continue
				}
				txHash, err := this.Layer2Sdk.NeoVM.InvokeNeoVMContract(0, 200000, nil, oscoreconfig.DefOscoreConfig.OscoreAccount, addr, []interface{}{"StoreUsedNum", []interface{}{key.UserId, key.OrderId, key.UsedNum}})
				if err != nil {
					log.Errorf("UpdateFreqDataBase invoke store num %s", err)
				}

				_, err = common.GetLayer2EventByTxHash(txHash.ToHexString())
				if err != nil {
					log.Errorf("UpdateFreqDataBase N.N.1 %s", err)
				}
				key.Layer2Time = time.Now().Unix()
				log.Infof("UpdateFreqDataBase layer2 addr: %s, tx hash %s, UserId: %s OrderId: %s, usedNum: %d", oscoreconfig.DefOscoreConfig.Layer2Contract, txHash.ToHexString(), key.UserId, key.OrderId, key.UsedNum)
			}
		}
	}
}

func (this *DBCache) UpdateReqLimitDataBase() {
	for {
		select {
		case apiKey := <-this.UpdateReqLimit:
			keyIn, ok := this.apiKeyCache.Load(apiKey)
			if !ok {
				log.Debugf("UpdateReqLimitDataBase error : apikey cache not exist")
				continue
			}

			key := keyIn.(*tables.APIKey)
			log.Debugf("UpdateReqLimitDataBase: %p ,%v", key, *key)
			err := dao.DefOscoreApiDB.UpdateApiKeyReqLimit(nil, key.ApiKey, key.RequestLimit, key.OutDate)
			if err != nil {
				log.Debugf("UpdateReqLimitDataBase: N.0 %s", err)
			}
		}
	}
}

func (this *DBCache) getApiIdFreqCounter(ApiId uint32) (*uint64, error) {
	apiCounterP, ok := this.apiFreqCache.Load(ApiId)
	if !ok || apiCounterP == nil {
		freq, err := dao.DefOscoreApiDB.QueryInvokeFreByApiId(nil, ApiId)
		if err != nil {
			return nil, err
		}
		this.apiFreqCache.Store(ApiId, &freq)
		return &freq, nil
	} else {
		return apiCounterP.(*uint64), nil
	}
}

// better check test api and api key id conflict.
func (this *DBCache) BeforeCheckApiKey(apiKey string, apiId uint32) (*tables.APIKey, *uint64, error) {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyCache(apiKey)
	if err != nil {
		return nil, nil, err
	}

	var apiCounterP *uint64
	if key.ApiId != apiId {
		return nil, nil, fmt.Errorf("this apikey: %s can not invoke this api", apiKey)
	}

	switch key.ApiKeyType {
	case tables.API_KEY_TYPE_DURATION:
		if time.Now().Unix() > key.OutDate {
			return nil, nil, fmt.Errorf("apiKey already out of Date.")
		}

	case tables.API_KEY_TYPE_COUNT:
		if key.UsedNum >= key.RequestLimit {
			return nil, nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.RequestLimit)
		}
	default:
		return nil, nil, fmt.Errorf("error apikey type.")
	}

	apiCounterP, err = this.getApiIdFreqCounter(apiId)
	if err != nil {
		return nil, nil, err
	}

	// record the usedNum of ApiKey no matter apikey type
	atomic.AddUint64(&key.UsedNum, 1)

	if !common.IsTestKey(apiKey) {
		atomic.AddUint64(apiCounterP, 1)
	}

	return key, apiCounterP, nil
}

func (this *DBCache) AtomicRenewApiKeyCacheCount(apiKey string, count uint64) error {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyCache(apiKey)
	if err != nil {
		return err
	}

	log.Debugf("AtomicRenewApiKeyCacheCount: %p ,%v", key, *key)
	atomic.AddUint64(&key.RequestLimit, count)
	this.UpdateReqLimit <- apiKey
	return nil
}

func (this *DBCache) AtomicRenewApiKeyCacheDate(apiKey string, month int32) error {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyCache(apiKey)
	if err != nil {
		return err
	}

	log.Debugf("AtomicRenewApiKeyCacheDate: %p ,%v", key, *key)
	outDate, err := utils.CalOutDateByMonth(key, month)
	if err != nil {
		log.Debugf("AtomicRenewApiKeyCacheDate: N.0 %s", err)
		return err
	}
	key.OutDate = outDate
	log.Debugf("AtomicRenewApiKeyCacheDate: %p ,%v", key, *key)
	this.UpdateReqLimit <- apiKey
	return nil
}

func (this *DBCache) getApiKeyCache(apiKey string) (*tables.APIKey, error) {
	keyIn, ok := this.apiKeyCache.Load(apiKey)
	if !ok || keyIn == nil {
		key, err := dao.DefOscoreApiDB.QueryApiKeyByApiKey(nil, apiKey)
		if err != nil {
			return nil, err
		}
		this.apiKeyCache.Store(apiKey, key)
		return key, nil
	} else {
		return keyIn.(*tables.APIKey), nil
	}
}

func (this *DBCache) updateApiKeyInvokeFre(key *tables.APIKey, freqCounter uint64) error {
	return dao.DefOscoreApiDB.UpdateApiKeyInvokeFre(nil, key.ApiKey, key.ApiId, key.UsedNum, freqCounter)
}
