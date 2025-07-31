package poolUtil

import "time"

type IPool interface {
	Submit(task func())
	SubmitWithId(id any, task func())
	Shutdown(timeout time.Duration) (isTimeout bool)
}
