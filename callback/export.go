package callback

type ICryptoModel interface {
	TableName() string
	CryptoColumns() []string
}
