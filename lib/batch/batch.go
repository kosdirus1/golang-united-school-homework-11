package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	res = make([]user, 0, n)
	var wg sync.WaitGroup
	sem := make(chan struct{}, pool)
	users := make(chan user)
	done := make(chan struct{})

	go fillSlice(users, done, &res)

	for i := 0; i < int(n); i++ {
		wg.Add(1)
		go func(i int) {
			sem <- struct{}{}
			users <- getOne(int64(i))
			<-sem
			wg.Done()
		}(i)
	}

	wg.Wait()
	done <- struct{}{}
	close(sem)
	close(users)
	close(done)
	return res
}

func fillSlice(users <-chan user, done <-chan struct{}, res *[]user) {
	for {
		select {
		case u := <-users:
			*res = append(*res, u)
		case <-done:
			return
		}
	}
}
