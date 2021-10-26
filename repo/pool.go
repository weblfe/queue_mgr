package repo

import (
	"github.com/panjf2000/ants/v2"
	"os"
	"strconv"
	"sync"
)

const DefaultPoolSize = 5000

type (
	processorRepository struct {
		size uint
		pool *ants.Pool
	}
	Runner func()

	Parallel struct {
		waitGroup sync.WaitGroup
		result    []interface{}
		errs      []error
		pool      *ants.Pool
		locker    sync.RWMutex
		size      int
		count     int
	}

)

func NewParallel(size ...int) *Parallel {
	size = append(size, 10)
	var parallel = new(Parallel)
	if size[0] > 0 {
		parallel.size = size[0]
	} else {
		parallel.size = 10
	}
	parallel.locker = sync.RWMutex{}
	parallel.waitGroup = sync.WaitGroup{}
	parallel.pool, _ = ants.NewPool(parallel.size)
	return parallel
}

func (p *Parallel) Add(caller func() interface{}) bool {
	if err := p.pool.Submit(func() {
		defer p.waitGroup.Done()
		var v = caller()
		p.append(v)
	}); err != nil {
		return false
	}
	p.incr()
	return true
}

func (p *Parallel) incr() {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.count++
	p.waitGroup.Add(1)
}

func (p *Parallel) AddTask(caller func() (interface{}, error)) bool {
	if err := p.pool.Submit(func() {
		defer p.waitGroup.Done()
		var v, err = caller()
		if err != nil {
			p.append(err)
		} else {
			p.append(v)
		}
	}); err != nil {
		return false
	}
	p.incr()
	return true
}

func (p *Parallel) append(v interface{}) {
	if v == nil {
		return
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	switch v.(type) {
	case error:
		p.errs = append(p.errs, v.(error))
	default:
		p.result = append(p.result, v)
	}
	p.count--
}

func (p *Parallel) Wait() ([]interface{}, []error) {
	p.waitGroup.Wait()
	return p.result, p.errs
}

func (p *Parallel) Empty() bool {
		if p.count > 0 {
				return false
		}
		return true
}

var defaultPool = newProcessorRepository()

func newProcessorRepository() *processorRepository {
	var rep = new(processorRepository)
	if err := rep.init(); err != nil {
		panic(err)
	}
	return rep
}

func (repo *processorRepository) init() error {
	if repo.size <= 0 {
		repo.size = repo.getDefaultSize()
	}
	var pool, err = ants.NewPool(int(repo.size))
	if err != nil {
		return err
	}
	repo.pool = pool
	return nil
}

func (repo *processorRepository) getDefaultSize() uint {
	var v = os.Getenv("POOL_SIZE")
	if v == "" {
		return DefaultPoolSize
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return DefaultPoolSize
	}
	return uint(n)
}

func (repo *processorRepository) Add(task func()) error {
	return repo.pool.Submit(task)
}

func (repo *processorRepository) GetPool() *ants.Pool {
	return repo.pool
}

// RegisterProcessor 注册一个携程处理器
func RegisterProcessor(runner Runner) error {
	return defaultPool.Add(runner)
}

func GetPoolRepo() *processorRepository {
	return defaultPool
}
