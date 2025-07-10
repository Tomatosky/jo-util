package poolUtil

import "time"

type IPool interface {
	Submit(task func())
	SubmitWithId(id int32, task func())
	Shutdown(timeout time.Duration) (isTimeout bool)
}
