package worker

import (
	"fmt"
	"time"

	"github.com/Bauka07/AP2/internal/server"
)

func StartWorker(s *server.Server) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer close(s.WorkerDone)

		for {
			select {
			case <- ticker.C:
				fmt.Printf("|Stats| requests=%v keys=%v\n", s.ReqCounter.Load(), s.St.Len())
			case <- s.WorkerStop:
				fmt.Println("Worker stopping...")
				return
			}
		}
	}()
}

func StopWorker(s *server.Server) {
	close(s.WorkerStop)
	<-s.WorkerDone
}
