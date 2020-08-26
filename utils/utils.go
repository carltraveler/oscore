package utils

import (
	"errors"
	"fmt"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/models/tables"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func ToStringByPrecise(bigNum *big.Int, decimals uint64) string {
	result := ""
	destStr := bigNum.String()
	destLen := uint64(len(destStr))
	if decimals >= destLen { // add "0.000..." at former of destStr
		var i uint64 = 0
		prefix := "0."
		for ; i < decimals-destLen; i++ {
			prefix += "0"
		}
		if destLen > 12 {
			destStr = destStr[:12]
		}
		result = prefix + destStr
	} else { // add "."
		pointIndex := destLen - decimals
		if len(destStr[pointIndex:]) > 12 {
			destStr = destStr[:pointIndex+12]
		}
		result = destStr[0:pointIndex] + "." + destStr[pointIndex:]
	}
	// delete no need "0" at last of result
	i := len(result) - 1
	for ; i >= 0; i-- {
		if result[i] != '0' {
			break
		}
	}
	result = result[:i+1]
	// delete "." at last of result
	if result[len(result)-1] == '.' {
		result = result[:len(result)-1]
	}
	return result
}

func ToIntByPrecise(str string, decimals uint64) *big.Int {
	result := new(big.Int)
	splits := strings.Split(str, ".")
	if len(splits) == 1 { // doesn't contain "."
		var i uint64 = 0
		for ; i < decimals; i++ {
			str += "0"
		}
		intValue, ok := new(big.Int).SetString(str, 10)
		if ok {
			result.Set(intValue)
		}
	} else if len(splits) == 2 {
		value := new(big.Int)
		ok := false
		floatLen := uint64(len(splits[1]))
		if floatLen <= decimals { // add "0" at last of str
			parseString := strings.Replace(str, ".", "", 1)
			var i uint64 = 0
			for ; i < decimals-floatLen; i++ {
				parseString += "0"
			}
			value, ok = value.SetString(parseString, 10)
		} else { // remove redundant digits after "."
			splits[1] = splits[1][:decimals]
			parseString := splits[0] + splits[1]
			value, ok = value.SetString(parseString, 10)
		}
		if ok {
			result.Set(value)
		}
	}

	return result
}

func CalOutDateByMonth(key *tables.APIKey, month int32) (int64, error) {
	if month == 0 {
		return 0, nil
	}

	if month > 120 || month < 0 {
		return 0, fmt.Errorf("can not over ten years or less than zero.")
	}

	if key.ApiKeyType != tables.API_KEY_TYPE_DURATION {
		return 0, fmt.Errorf("error spec type %d.", key.ApiKeyType)
	}

	var outDate int64
	// less than current is outOfDate.
	if key.OutDate < time.Now().Unix() || key.OutDate == 0 {
		// if out of date. or init. use current time.
		log.Debugf("CalOutDateByMonth: Y.0.0 %v %d", *key, time.Now().Unix())
		outDate = time.Now().AddDate(0, int(month), 0).Unix()
		log.Debugf("CalOutDateByMonth: %s Y.0.1 %d", key.ApiKey, outDate)
	} else {
		// if not out of date. add onto the OutDate.
		log.Debugf("CalOutDateByMonth: %s Y.1 %d", key.ApiKey, outDate)
		outDate = time.Unix(key.OutDate, 0).AddDate(0, int(month), 0).Unix()
		log.Debugf("CalOutDateByMonth: %s Y.2 %d", key.ApiKey, outDate)
	}

	return outDate, nil
}

func GetOutDateByCurrAddDuration(duration string) (int64, error) {
	if duration == "" {
		return 0, nil
	}

	arr := strings.Split(duration, ":")
	if len(arr) != 3 {
		return 0, errors.New("wrong duration type.")
	}

	year, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, fmt.Errorf("%s", err)
	}
	month, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, fmt.Errorf("%s", err)
	}
	day, err := strconv.Atoi(arr[2])
	if err != nil {
		return 0, fmt.Errorf("%s", err)
	}

	if year > 10 || month > 120 || day > 3600 {
		return 0, fmt.Errorf("can not over ten years each.")
	}

	outDate := time.Now().AddDate(year, month, day).Unix()

	return outDate, nil
}
