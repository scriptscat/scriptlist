package repository

type RateRule struct {
	Name     string
	Interval int64
	DayMax   int64
}

type RateUserInfo struct {
	Uid int64
}

type Rate interface {
	GetLastOpTime(user, operation string) (int64, error)
	GetDayOpCnt(user, op string) (int64, error)
	SetLastOpTime(user, operation string, t int64) error
}
