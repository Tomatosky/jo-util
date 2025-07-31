package poolUtil

import (
	"github.com/Tomatosky/jo-util/numberUtil"
	"time"
)

type IPool[T numberUtil.Number] interface {
	Submit(task func())
	SubmitWithId(id T, task func())
	Shutdown(timeout time.Duration) (isTimeout bool)
}
