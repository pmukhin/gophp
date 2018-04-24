package object

import "testing"

func TestLocalStorage_Get(t *testing.T) {
	st := newLocalStorage()

	st.Set("var1", &IntegerObject{Value: 1})
	st.Set("var2", &IntegerObject{Value: 2})

	intSt := newLocalStorage()
	intSt.SetParent(st)

	intSt.Set("var2", &IntegerObject{Value: 5})

	v, e := intSt.Get("var2")
	if e != nil {
		t.Error(e)
	}
	if v.(*IntegerObject).Value != 5 {
		t.Errorf("expected value from current context")
	}
	v, e = intSt.Get("var1")
	if e != nil {
		t.Error(e)
	}
	if v.(*IntegerObject).Value != 1 {
		t.Errorf("expected value from parent context")
	}
}
