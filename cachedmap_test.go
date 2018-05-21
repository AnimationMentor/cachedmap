package cachedmap

import (
	"testing"
	"time"
)

func Test_test1(t *testing.T) {
	cm := NewCachedMap(1*time.Second, 10*time.Second, nil)

	cm.Set("1", 150)

	v, ok := cm.Get("1")

	if !ok {
		t.Errorf("test %d: unable to get value I just set", 1)
	}

	i, ok := v.(int)

	if !ok {
		t.Errorf("test %d: unable to cast value to int (is %T)", 2, v)
	}

	if i != 150 {
		t.Errorf("test %d: value mismatch", 3)
	}

	time.Sleep(2 * time.Second)

	_, ok = cm.Get("1")
	if ok {
		t.Errorf("test %d: that value should have timed out but hasn't", 4)
	}
}
