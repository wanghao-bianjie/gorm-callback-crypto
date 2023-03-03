package callback

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func BeforeCreate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}: //[]map
				dest := db.Statement.Dest.([]map[string]interface{})
				slice, err := encryptDestMapSlice(db.Statement.Schema.Table, dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = slice
				return
			case *[]map[string]interface{}: //*[]map
				dest := db.Statement.Dest.(*[]map[string]interface{})
				slice, err := encryptDestMapSlice(db.Statement.Schema.Table, *dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = slice
				return
			case map[string]interface{}: //map
				dest := db.Statement.Dest.(map[string]interface{})
				destMap, err := encryptDestMap(db.Statement.Schema.Table, dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = destMap
				return
			case *map[string]interface{}: //*map
				dest := db.Statement.Dest.(*map[string]interface{})
				destMap, err := encryptDestMap(db.Statement.Schema.Table, *dest)
				if err != nil {
					db.AddError(err)
					return
				}
				db.Statement.Dest = destMap
				return
			}
		}

		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for _, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					index := destReflectValue.Index(i)
					if destReflectValue.Index(i).Kind() != reflect.Struct {
						continue
					}
					if fieldValue, isZero := field.ValueOf(index); !isZero { // 从字段中获取数值
						if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								encrypt, err := _options.beforeFn(s)
								if err != nil {
									db.AddError(err)
									return
								}
								_ = db.AddError(field.Set(index, encrypt))
							} else {
								db.AddError(fmt.Errorf("encrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
								return
							}
						}
					}
				}
			case reflect.Struct: //struct
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero {
					if dbTableColumnHandleMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
						if s, ok := fieldValue.(string); ok {
							encrypt, err := _options.beforeFn(s)
							if err != nil {
								db.AddError(err)
								return
							}
							_ = db.AddError(field.Set(destReflectValue, encrypt))
						} else {
							db.AddError(fmt.Errorf("encrypt table:%s,column:%s,but not string", db.Statement.Schema.Table, field.DBName))
							return
						}
					}
				}
			}
		}
	}
}

func AfterCreate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}, map[string]interface{}, *[]map[string]interface{}, *map[string]interface{}:
				return
			}
		}

		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for _, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					// 从字段中获取数值
					if destReflectValue.Index(i).Kind() != reflect.Struct {
						continue
					}
					if fieldValue, isZero := field.ValueOf(destReflectValue.Index(i)); !isZero {
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
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero {
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
