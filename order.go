package wego

type Order interface {
	Unify(m Map) Map
	Close(no string) Map
	//Query(Map) Map
	QueryByTransactionId(id string) Map
	QueryByOutTradeNumber(no string) Map
}

type order struct {
	Config
	Payment
}

func NewOrder(application Application, config Config) Order {
	return &order{
		//app:    application,
		Payment: application.Payment(),
		Config:  config,
	}
}

func (o *order) Unify(m Map) Map {
	if !m.Has("spbill_create_ip") {
		if m.Get("trade_type") == "NATIVE" {
			m.Set("spbill_create_ip", GetServerIp())
		}
		//TODO: getclientip with request
	}

	m.Set("appid", o.Get("app_id"))
	//$params['notify_url'] = $params['notify_url'] ?? $this->app['config']['notify_url'];
	if !m.Has("notify_url") {
		m.Set("notify_url", o.Get("notify_url"))
	}
	return o.Request(UNIFIEDORDER_URL_SUFFIX, m)
}

//func (o *order) request(url string, m Map) Map {
//	return o.Request(o.Link(url), m, "post", nil)
//}

/**
* 作用：关闭订单
* 场景：公共号支付、扫码支付、APP支付
* @param data 向wxpay post的请求数据
* @return PayData, error API返回数据
 */
func (o *order) Close(no string) Map {
	m := make(Map)
	m.Set("appid", o.Get("app_id"))
	m.Set("out_trade_no", no)
	return o.Request(CLOSEORDER_URL_SUFFIX, m)
}

/** QueryOrder
* 作用：查询订单
* 场景：刷卡支付、公共号支付、扫码支付、APP支付
* @param data 向wxpay post的请求数据
* @param connectTimeoutMs 连接超时时间，单位是毫秒
* @param readTimeoutMs 读超时时间，单位是毫秒
* @return PayData, error API返回数据
 */
func (o *order) query(m Map) Map {
	m.Set("appid", o.Get("app_id"))
	return o.Request(ORDERQUERY_URL_SUFFIX, m)
}

func (o *order) QueryByTransactionId(id string) Map {
	m := make(Map)
	m.Set("transaction_id", id)
	return o.query(m)
}

func (o *order) QueryByOutTradeNumber(no string) Map {
	m := make(Map)
	m.Set("out_trade_no", no)
	return o.query(m)
}
