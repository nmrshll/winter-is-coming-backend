package cancellable

import "time"

func Run(timePeriod time.Duration, loopFunc func() error) (cancelFunc func()) {
	cancelChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(timePeriod):
				// err := expectRequestType(conn, playerShotHandler(currentGame))
				// if err != nil {
				// 	errors.Log(err)
				// 	return // close connection after error
				// }
				err := loopFunc()
				if err != nil {
					return
				}
			case <-cancelChan:
				return
			}
		}
	}()
	return func() {
		close(cancelChan)
	}
}
