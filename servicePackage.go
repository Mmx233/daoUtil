package daoUtil

import (
	"container/list"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"time"
)

func init() {
	go serviceLocks.worker()
}

type serviceLock struct {
	locks     sync.Map
	renewChan chan string
	times     sync.Map
	timeStack list.List
	listLock  sync.Mutex
}

var serviceLocks = serviceLock{
	renewChan: make(chan string, 1),
}

type delLockInfo struct {
	Name string
	Time time.Time
}

func (a *serviceLock) worker() {
	go func() {
		for {
			n := <-a.renewChan
			var e *list.Element
			a.listLock.Lock()
			if i, ok := a.times.Load(n); !ok {
				e = a.timeStack.PushBack(&delLockInfo{Name: n, Time: time.Now()})
			} else {
				e = i.(*list.Element)
				a.timeStack.Remove(e)
				e = a.timeStack.PushBack(&delLockInfo{Name: n, Time: time.Now()})
			}
			a.listLock.Unlock()
			a.times.Store(n, e)
		}
	}()

	go func() {
		var s time.Duration
		for {
			a.listLock.Lock()
			e := a.timeStack.Front()
			if e == nil {
				//GC
				a.locks = sync.Map{}
				a.times = sync.Map{}
				a.listLock.Unlock()
				time.Sleep(time.Minute * 5)
				continue
			}
			info := e.Value.(*delLockInfo)
			s = time.Minute*5 - time.Since(info.Time)
			if s > 0 {
				a.listLock.Unlock()
				time.Sleep(s)
			} else {
				a.timeStack.Remove(e)
				a.times.Delete(info.Name)
				a.listLock.Unlock()
			}
		}
	}()
}

func (a *serviceLock) Lock(n string) {
	a.renewChan <- n
	t, _ := a.locks.LoadOrStore(n, &sync.Mutex{})
	t.(*sync.Mutex).Lock()
}

func (a *serviceLock) UnLock(n string) {
	t, _ := a.locks.Load(n)
	t.(*sync.Mutex).Unlock()
}

type ServicePackage struct {
	Tx        *gorm.DB
	name      string
	committed bool
}

func (ServicePackage) Begin(a interface{}, key string) ServicePackage {
	var name string
	if a != nil {
		name = reflect.TypeOf(a).Elem().Name() + "-" + key
		serviceLocks.Lock(name)
	}
	tx := Begin()
	return ServicePackage{
		Tx:   tx,
		name: name,
	}
}

func (ServicePackage) BeginWith(tx *gorm.DB) ServicePackage {
	if tx == nil {
		tx = c.DB
	}
	return ServicePackage{
		Tx: tx,
	}
}

func (a *ServicePackage) end(e func() *gorm.DB) {
	if a.committed {
		return
	}
	if a.name != "" {
		serviceLocks.UnLock(a.name)
	}
	a.committed = true
	e()
}

// RollBack 回滚，使用行锁时必须defer
func (a *ServicePackage) RollBack() {
	a.end(a.Tx.Rollback)
}

func (a *ServicePackage) Commit() {
	a.end(a.Tx.Commit)
}
