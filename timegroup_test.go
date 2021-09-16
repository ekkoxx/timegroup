package timegroup

import (
	"errors"
	"testing"
	"time"
)

func Test_Success(t *testing.T) {
	t.Parallel()
	v1, v2, v3 := 1, 1, 1
	g := New()
	g.Go(func() error {
		v1++
		return nil
	})
	g.Go(func() error {
		v2++
		return nil
	})
	g.Go(func() error {
		v3++
		return nil
	})
	err := g.Wait()
	if err != nil {
		t.Errorf("wait err:%s", err.Error())
	}
	if v1 != 2 || v2 != 2 || v3 != 2 {
		t.Errorf("exec err.v1:%d;v2:%d;v3:%d", v1, v2, v3)
	}
}

func Test_FailTimeOut(t *testing.T) {
	t.Parallel()
	var v1err = errors.New("v1 err")
	v1, v2, v3 := 1, 1, 1
	g := New()
	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		v1++
		return v1err
	})
	g.Go(func() error {
		time.Sleep(time.Millisecond * 100)
		v2++
		return nil
	})
	g.Go(func() error {
		v3++
		return nil
	})
	err := g.WaitTimeout(time.Second)
	if err != v1err {
		t.Errorf("wait err:%s", err.Error())
	}
	if v1 != 2 || v2 != 1 || v3 != 2 {
		t.Errorf("exec err.v1:%d;v2:%d;v3:%d", v1, v2, v3)
	}
}
