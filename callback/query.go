package callback

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

//AfterQuery 只支持查询返回的结果是单表model的结构体，map、自定义的其他结构体（联表、部分字段）暂不支持
func AfterQuery(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && db.RowsAffected > 0 && !db.Statement.SkipHooks {
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for fieldIndex, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					index := destReflectValue.Index(i)
					if !(index.Kind() == reflect.Struct || (index.Kind() == reflect.Ptr && index.Elem().Kind() == reflect.Struct)) {
						continue
					}
					if index.Kind() == reflect.Ptr {
						index = index.Elem()
					}
					if index.NumField() != len(db.Statement.Schema.Fields) {
						return
					}
					if fieldValue, isZero := field.ValueOf(index); !isZero { // 从字段中获取数值
						if index.Type().Field(fieldIndex).Name != field.Name || index.Type().Field(fieldIndex).Type.Kind() != field.FieldType.Kind() {
							return
						}
						if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								decrypt, err := _options.afterFn(s)
								if err != nil {
									db.AddError(err)
									return
								}
								_ = db.AddError(field.Set(destReflectValue.Index(i), decrypt))
							} else {
								db.AddError(fmt.Errorf("decrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
								return
							}
						}
					}
				}
			case reflect.Struct: //struct
				if destReflectValue.NumField() != len(db.Statement.Schema.Fields) {
					return
				}
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero { // 从字段中获取数值
					if destReflectValue.Type().Field(fieldIndex).Name != field.Name {
						return
					}
					if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
						if s, ok := fieldValue.(string); ok {
							decrypt, err := _options.afterFn(s)
							if err != nil {
								db.AddError(err)
								return
							}
							_ = db.AddError(field.Set(destReflectValue, decrypt))
						} else {
							db.AddError(fmt.Errorf("decrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
							return
						}
					}
				}
			}
		}
	}
}
