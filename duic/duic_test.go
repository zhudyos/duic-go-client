package duic

import "testing"

func init() {
	BaseUri = "https://duic.zhudy.io/api/v1"
	Name = "unit-test"
	Profiles = "test"

	Init()
}

func TestBool(t *testing.T) {
	v, err := Bool("test.bool")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.bool: %v\n", v)
	}
}

func TestBool2(t *testing.T) {
	v := Bool2("test.bool", false)
	t.Logf("test.bool: %v\n", v)
}

func TestBool2_2(t *testing.T) {
	v := Bool2("test.bool.xxxxx", true)
	if !v {
		t.FailNow()
	} else {
		t.Logf("test.bool: %v\n", v)
	}
}

func TestInt(t *testing.T) {
	v, err := Int("test.int")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.int: %v\n", v)
	}
}

func TestInt2(t *testing.T) {
	var defVar = -55555
	v := Int2("test.int.xxxxx", defVar)
	if v != defVar {
		t.FailNow()
	} else {
		t.Logf("test.int.xxxxx: %v\n", v)
	}
}

func TestInt64(t *testing.T) {
	v, err := Int64("test.int")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.int: %v\n", v)
	}
}

func TestInt642(t *testing.T) {
	var defVar int64 = -99999
	v := Int642("test.int", defVar)

	if v == defVar {
		t.FailNow()
	} else {
		t.Logf("test.int: %v\n", v)
	}
}

func TestInt642_2(t *testing.T) {
	var defVar int64 = -99999
	v := Int642("test.int.xxxxx", defVar)
	if defVar != v {
		t.FailNow()
	} else {
		t.Logf("test.int.xxxxx: %v\n", defVar)
	}
}

func TestFloat64(t *testing.T) {
	v, err := Float64("test.float")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.float: %v\n", v)
	}
}

func TestFloat642(t *testing.T) {
	defVar := 9.9999999
	v := Float642("test.float", defVar)
	if v == defVar {
		t.FailNow()
	} else {
		t.Logf("test.float: %v\n", defVar)
	}
}

func TestFloat642_2(t *testing.T) {
	defVar := 9.9999999
	v := Float642("test.float.xxxxx", defVar)
	if v != defVar {
		t.FailNow()
	} else {
		t.Logf("test.float.xxxxx: %v\n", defVar)
	}
}

func TestString(t *testing.T) {
	v, err := String("test.string")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.string: %v\n", v)
	}
}

func TestString2(t *testing.T) {
	defVar := "xyz"
	v := String2("test.string", defVar)
	if v == defVar {
		t.FailNow()
	} else {
		t.Logf("test.string: %v\n", v)
	}
}

func TestString2_2(t *testing.T) {
	defVar := "xyz"
	v := String2("test.string.xxxxx", defVar)
	if v != defVar {
		t.FailNow()
	} else {
		t.Logf("test.string.xxxxx: %v\n", v)
	}
}

func TestArray(t *testing.T) {
	v, err := Array("test.array")

	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.array: %v\n", v)
	}
}

func TestObject(t *testing.T) {
	v, err := Object("test.object")
	if err != nil {
		t.FailNow()
	} else {
		t.Logf("test.object: %v\n", v)
	}
}
