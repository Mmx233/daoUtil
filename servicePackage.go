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
			a.listLock.Lock()
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
	Tx *gorm.DB
}

func (ServicePackage) Begin(a interface{}, key string) ServicePackage {
	var name string
	if a != nil {
		name = reflect.TypeOf(a).Elem().Name() + "-" + key
	}
	tx := Begin()
	tx.Set("name", name)
	return ServicePackage{
		Tx: tx,
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

func (a *ServicePackage) committed() bool {
	ed, ok := a.Tx.Get("committed")
	return ok && ed.(bool) == true
}

func (a *ServicePackage) mark() {
	a.Tx.Set("committed", true)
}

func (a *ServicePackage) name() string {
	o, _ := a.Tx.Get("name")
	return o.(string)
}

// RollBack 回滚，使用行锁时必须defer
func (a *ServicePackage) RollBack() {
	if n := a.name(); n != "" {
		serviceLocks.UnLock(n)
	}
	if !a.committed() {
		a.mark()
		a.Tx.Rollback()
	}
}

func (a *ServicePackage) Commit() {
	if !a.committed() {
		a.mark()
		a.Tx.Commit()
	}
}
