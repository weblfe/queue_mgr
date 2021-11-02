package domain

import (
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/models"
	"sync"
)

type (
	serverDomain struct {
		// 队列状态树
		trees map[string]*stateTree
		// 队列信息存储
		queues map[string]*QueueInfo
	}

	stateTree struct {
		safe    sync.RWMutex
		items   []*QueueInfo
		indexes map[string]int
		state   entity.QueueState
	}

	QueueInfo struct {
		advisor  facede.Advisor
		Base     *models.QueueInfo
		Bind     *models.QueryBindInfo
		consumer *models.ConsumerInfo
	}
)

func NewStateTree(state entity.QueueState) *stateTree {
	var tree = new(stateTree)
	tree.state = state
	tree.indexes = make(map[string]int)
	return tree
}

func (tree *stateTree) Check() bool {
	if tree.state.Check() {
		return true
	}
	return false
}

func (tree *stateTree) Name() string {
	return tree.state.String()
}

func (tree *stateTree) Len() int {
	return len(tree.items)
}

func (tree *stateTree) Add(item *QueueInfo) int {
	if item == nil {
		return 0
	}
	tree.safe.Lock()
	defer tree.safe.Unlock()
	var name = item.Queue()
	if len(tree.items) > 0 {
		if n, ok := tree.indexes[name]; ok {
			return n
		}
	}
	tree.items = append(tree.items, item)
	var index = len(tree.items)
	tree.indexes[name] = index
	return index
}

func (tree *stateTree) Remove(queue string) bool {
	tree.safe.Lock()
	defer tree.safe.Unlock()
	var size = len(tree.items)
	if size > 0 {
		if n, ok := tree.indexes[queue]; ok {
			if tree.delete(n) {
				delete(tree.indexes, queue)
				return true
			}
		}
	}
	return false
}

func (tree *stateTree) Clear() bool {
	tree.safe.Lock()
	defer tree.safe.Unlock()
	if tree.items == nil {
		return true
	}
	tree.items = nil
	tree.indexes = make(map[string]int)
	return true
}

func (tree *stateTree) IndexOf(index int) (*QueueInfo, bool) {
	tree.safe.Lock()
	defer tree.safe.Unlock()
	var size = len(tree.items)
	if size > 0 && index <= size && index >= 1 {
		return tree.items[index-1], true
	}
	return nil, false
}

func (tree *stateTree) Index(queue string) int {
	var size = len(tree.items)
	if size <= 0 || queue == "" {
		return -1
	}
	tree.safe.Lock()
	defer tree.safe.Unlock()
	if index, ok := tree.indexes[queue]; ok && index <= size && index >= 1 {
		return index
	}
	return -1
}

func (tree *stateTree) Get(queue string) (*QueueInfo, bool) {
	tree.safe.Lock()
	defer tree.safe.Unlock()
	var (
		size      = len(tree.items)
		index, ok = tree.indexes[queue]
	)
	if size > 0 && ok && index <= size {
		return tree.items[index], true
	}
	return nil, false
}

func (tree *stateTree) delete(index int) bool {
	var sizeOf = len(tree.items)
	if sizeOf >= index && sizeOf != 0 {
		defer tree.reIndex()
		if index == 1 {
			tree.items = tree.items[index:]
			return true
		}
		if index == sizeOf {
			tree.items = tree.items[:index]
			return true
		}
		tree.items = append(tree.items[:index-1], tree.items[index:]...)
		return true
	}
	return false
}

func (tree *stateTree) reIndex() {
	for i, v := range tree.items {
		if v == nil {
			continue
		}
		tree.indexes[v.Queue()] = i + 1
	}
}

func (tree *stateTree) Exists(queue string) bool {
	tree.safe.Lock()
	defer tree.safe.Unlock()
	if _, ok := tree.indexes[queue]; ok {
		return ok
	}
	return false
}

func (tree *stateTree) ForEach(each func(i int, info *QueueInfo)) {
	if each == nil {
		return
	}
	for i, v := range tree.items {
		each(i, v)
	}
}

func (info *QueueInfo) Queue() string {
	if info == nil {
		return ""
	}
	if info.Base == nil {
		return ""
	}
	return info.Base.Name
}
