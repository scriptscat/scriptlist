package errs

import (
	"fmt"

	"github.com/scriptscat/scriptweb/internal/domain/safe/repository"
)

func NewOperationTimeToShort(rule *repository.RateRule) error {
	return NewBadRequestError(4001, fmt.Sprintf("两次操作时间过断,请%d秒后重试", rule.Interval))
}
