package domain

import (
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/models"
	"testing"
)

func TestStateTree_Remove(t *testing.T) {
	var tree = NewStateTree(entity.Wait)
	tree.Add(&QueueInfo{
		Base: &models.QueueInfo{
			Name: "admin",
		},
	})
	tree.Add(&QueueInfo{
		Base: &models.QueueInfo{
			Name: "user",
		},
	})
	tree.Add(&QueueInfo{
		Base: &models.QueueInfo{
			Name: "caller",
		},
	})

	if tree.Len() != 3 {
		t.Error("异常数据")
	}
	var c = tree.Cursor()
	if v, ok := c.Prev(); ok || v != nil {
		t.Error("游标异常")
	}
	for {
		p, ok := c.Node()
		if ok {
			if p == nil {
				t.Error("游标获取节点数据异常")
			}
		}
		c, ok = c.Next()
		if !ok {
			break
		}
	}
	tree.Remove("admin")
	if tree.Len() != 2 {
		t.Error("remove 异常数据")
	}
	if tree.Index("admin") != -1 {
		t.Error("索引 异常数据")
	}
	if tree.Index("caller") != 2 {
		t.Error("索引 异常数据")
	}
	tree.Clear()
}
