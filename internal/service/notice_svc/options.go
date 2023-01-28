package notice_svc

type options struct {
	// from 发送者用户id
	from int64
	// 发送参数
	params interface{}
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	ret := &options{}
	for _, v := range opts {
		v(ret)
	}
	return ret
}

func WithFrom(from int64) Option {
	return func(o *options) {
		o.from = from
	}
}

func WithParams(params interface{}) Option {
	return func(o *options) {
		o.params = params
	}
}
