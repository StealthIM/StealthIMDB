package errorcode

const (
	// DBGatewayRedisServiceError Redis数据库错误
	DBGatewayRedisServiceError int32 = 2000 + iota
	// DBGatewaySQLTransactionError SQL事务错误
	DBGatewaySQLTransactionError
	// DBGatewaySQLExecuteError SQL执行错误
	DBGatewaySQLExecuteError
	// DBGatewaySQLUnknownTypeError SQL数据类型未知错误
	DBGatewaySQLUnknownTypeError
)
