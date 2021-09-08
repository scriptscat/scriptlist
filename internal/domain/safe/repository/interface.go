package repository

type RateRule struct {
	Name     string
	Interval int64
}

type RateUserInfo struct {
	Uid int64
}

type Rate interface {
	GetLastOpTime(user string, operation string) (int64, error)
	SetLastOpTime(user string, operation string, t int64) error
}
