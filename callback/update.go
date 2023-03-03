package callback

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func BeforeUpdate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && !db.Statement.SkipHooks {
		destIsMap := false
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case map[string]interface{}: //map
				dest := db.Statement.Dest.(map[string]interface{})
				destMap, err := encryptDestMap(db.Statement.Table, dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = destMap
				destIsMap = true
			case *map[string]interface{}: //*map
				dest := db.Statement.Dest.(*map[string]interface{})
				destMap, err := encryptDestMap(db.Statement.Table, *dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = destMap
				destIsMap = true
			}
		}

		if db.Statement.Schema == nil {
			return
		}
		fn := func(reflectValue reflect.Value) {
			for _, field := range db.Statement.Schema.Fields {
				switch reflectValue.Kind() {
				case reflect.Struct: //struct
					if fieldValue, isZero := field.ValueOf(reflectValue); !isZero { // 从字段中获取数值
						if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								encrypt, err := _options.beforeFn(s)
								if err != nil {
									db.AddError(err)
									return
								}
								if reflectValue.CanSet() {
									_ = db.AddError(field.Set(reflectValue, encrypt))
								} else {
									newStruct := reflect.New(reflectValue.Type())
									for i := 0; i < reflectValue.Type().NumField(); i++ {
										if reflectValue.Type().Field(i).Name == field.Name {
											newStruct.Elem().Field(i).SetString(encrypt)
										} else {
											newStruct.Elem().Field(i).Set(reflectValue.Field(i))
										}
									}
									db.Statement.Dest = newStruct.Elem().Interface()
								}
							} else {
								db.AddError(fmt.Errorf("encrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
								return
							}
						}
					}
				}
			}
		}
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		if !destIsMap {
			fn(destReflectValue)
		}
		if destReflectValue != db.Statement.ReflectValue {
			fn(db.Statement.ReflectValue)
		}
	}
}

func AfterUpdate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && !db.Statement.SkipHooks && db.Statement.Schema != nil {
		destIsMap := false
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}, map[string]interface{}, *[]map[string]interface{}, *map[string]interface{}:
				destIsMap = true
			}
		}

		fn := func(reflectValue reflect.Value) {
			for _, field := range db.Statement.Schema.Fields {
				switch reflectValue.Kind() {
				case reflect.Struct: //struct
					if fieldValue, isZero := field.ValueOf(reflectValue); !isZero { // 从字段中获取数值
						if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								if reflectValue.CanSet() {
									decrypt, err := _options.afterFn(s)
									if err != nil {
										db.AddError(err)
										return
									}
									_ = db.AddError(field.Set(reflectValue, decrypt))
								}
							} else {
								db.AddError(fmt.Errorf("decrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
								return
							}
						}
					}
				}
			}
		}
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		if !destIsMap {
			fn(destReflectValue)
		}
		if destReflectValue != db.Statement.ReflectValue {
			fn(db.Statement.ReflectValue)
		}
	}
}
