package callback

import (
	"fmt"
	"reflect"

	"github.com/wanghao-bianjie/gorm-callback-crypto/util/aes"
)

var dbTableColumnHandleMap = map[[2]string]bool{
	//example
	//[2]string{"user", "name"}:  false,
	//[2]string{"user", "id_no"}: true,
}

func registerAesTableColumns(cryptoModels []ICryptoModel) {
	for _, cryptoModel := range cryptoModels {
		for _, column := range cryptoModel.CryptoColumns() {
			dbTableColumnHandleMap[[2]string{cryptoModel.TableName(), column}] = true
		}
	}
}

func getReflectValueElem(i interface{}) reflect.Value {
	value := reflect.ValueOf(i)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func aesEncryptToBase64(str string) (string, error) {
	if str == "" {
		return str, nil
	}
	//if _, err := aes.CBCPKCS7DecryptFromBase64(str, _options.defaultAesFnKey); err == nil {
	//	return str, nil
	//}
	return aes.CBCPKCS7EncryptToBase64([]byte(str), _options.defaultAesFnKey)
}

func aesDecryptFromBase64(str string) (string, error) {
	if str == "" {
		return str, nil
	}
	return aes.CBCPKCS7DecryptFromBase64(str, _options.defaultAesFnKey)
}

func encryptDestMap(table string, dest map[string]interface{}) (map[string]interface{}, error) {
	var newDest = make(map[string]interface{}, len(dest))
	for field, value := range dest {
		if dbTableColumnHandleMap[[2]string{table, field}] {
			if s, ok := value.(string); ok {
				encrypt, err := _options.beforeFn(s)
				if err != nil {
					return nil, err
				}
				newDest[field] = encrypt
			} else {
				return nil, fmt.Errorf("encrypt table:%s,column:%s,but not string", table, field)
			}
		} else {
			newDest[field] = value
		}
	}
	return newDest, nil
}

func encryptDestMapSlice(table string, dest []map[string]interface{}) ([]map[string]interface{}, error) {
	var newDest = make([]map[string]interface{}, 0, len(dest))
	for _, m := range dest {
		destMap, err := encryptDestMap(table, m)
		if err != nil {
			return nil, err
		}
		newDest = append(newDest, destMap)
	}
	return newDest, nil
}
