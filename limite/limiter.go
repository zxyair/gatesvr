package limite

type Limiter interface {
	GetToken() bool
}
