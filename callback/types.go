package callback

import "gorm.io/gorm"

var (
	_options = options{
		defaultAesFnKey: []byte("AesKeyLengthIs16"),
		beforeFn:        aesEncryptToBase64,
		afterFn:         aesDecryptFromBase64,
	}
)

type options struct {
	defaultAesFnKey []byte
	beforeFn        func(origin string) (store string, err error)
	afterFn         func(store string) (origin string, err error)
}

type Option interface {
	apply(*options)
}

type defaultAesFnKeyOption []byte

func (o defaultAesFnKeyOption) apply(opt *options) {
	opt.defaultAesFnKey = o
}

type beforeFnOption func(origin string) (store string, err error)

func (o beforeFnOption) apply(opt *options) {
	opt.beforeFn = o
}

type afterFnOption func(origin string) (store string, err error)

func (o afterFnOption) apply(opt *options) {
	opt.beforeFn = o
}

func WithDefaultAesFnKey(aesKey []byte) Option {
	return defaultAesFnKeyOption(aesKey)
}

func WithBeforeHandleFn(fn func(string) (string, error)) Option {
	return beforeFnOption(fn)
}

func WithAfterHandleFn(fn func(string) (string, error)) Option {
	return afterFnOption(fn)
}

func Register(db *gorm.DB, cryptoModels []ICryptoModel, opts ...Option) error {
	if err := db.Callback().Create().Before("gorm:create").After("gorm:before_create").Register("customize:before_create", BeforeCreate); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Before("gorm:after_create").Register("customize:after_create", AfterCreate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").After("gorm:before_update").Register("customize:before_update", BeforeUpdate); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Before("gorm:after_update").Register("customize:after_update", AfterUpdate); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Before("gorm:after_query").Register("customize:after_query", AfterQuery); err != nil {
		return err
	}
	for _, o := range opts {
		o.apply(&_options)
	}
	registerAesTableColumns(cryptoModels)
	return nil
}
