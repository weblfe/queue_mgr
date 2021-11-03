package domain

import (
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/models"
	"sync"
)

type (

	// 队列状态
	stateTree struct {
		safe    sync.RWMutex
		items   []*QueueInfo
		indexes map[string]int
		state   entity.QueueState
	}

	treeContainer map[State]*stateTree

	QueueInfo struct {
		advisor  facede.Advisor
		Base     *models.QueueInfo
		Bind     *models.QueryBindInfo
		consumer *models.ConsumerInfo
	}

	Cursor struct {
		index  int
		parent *stateTree
	}
)

func (container treeContainer) Register(state entity.QueueState, queue *QueueInfo) bool {
	if queue == nil {
		return false
	}
	var (
		queueName = queue.Queue()
		stateKey  = State(state.String())
	)
	if queueName == "" {
		return false
	}
	if v, ok := container[stateKey]; ok {
		return v.Add(queue) > 0
	}
	var tree = NewStateTree(state)
	tree.Add(queue)
	container[stateKey] = tree
	return true
}

func (container treeContainer) Remove(state entity.QueueState, queue string) bool {
	var stateKey = State(state.String())
	if v, ok := container[stateKey]; ok {
		return v.Remove(queue)
	}
	return false
}

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

func (tree *stateTree) State() entity.QueueState {
	return tree.state
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
		if info.Bind != nil {
			return info.Bind.Queue
		}
	}
	return info.Base.Name
}

func (tree *stateTree) Cursor(index ...int) *Cursor {
	if len(index) <= 0 || index[0] <= 0 {
		return &Cursor{
			index:  0,
			parent: tree,
		}
	}
	if tree.Len() < index[0] {
		return &Cursor{}
	}
	return &Cursor{
		index:  index[0],
		parent: tree,
	}
}

func (c *Cursor) Index() int {
	if c == nil || c.parent == nil {
		return -1
	}
	return c.index
}

func (c *Cursor) Next() (*Cursor, bool) {
	if c == nil || c.parent == nil {
		return nil, false
	}
	if c.parent.Len() > c.index && c.index >= 0 {
		index := c.index + 1
		return &Cursor{
			index:  index,
			parent: c.parent,
		}, true
	}
	return nil, false
}

func (c *Cursor) Prev() (*Cursor, bool) {
	if c == nil || c.parent == nil {
		return nil, false
	}
	if c.parent.Len() >= c.index && c.index > 1 {
		index := c.index - 1
		return &Cursor{
			index:  index,
			parent: c.parent,
		}, true
	}
	return nil, false
}

func (c *Cursor) Node() (*QueueInfo, bool) {
	if c == nil || c.parent == nil {
		return nil, false
	}
	return c.parent.IndexOf(c.index)
}
