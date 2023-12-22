package scheduler

import (
	"testing"
	"time"
)

// TestNormal 测试会每一秒打印一个ok，不会超过5个ok
func TestNormal(t *testing.T) {
	s := NewScheduler()

	ch := make(chan struct{})
	s.Register(NewJob(func() error {
		close(ch)
		return nil
	}, WithOnce(time.Now().Add(5*time.Second))))

	s.Register(NewJob(func() error {
		println("ok")
		return nil
	}, WithInterval(1*time.Second)))

	go s.run()

	<-ch
}
