package number

import (
	"gitee.com/conero/uymas/v2/util/unit"
	"testing"
)

func TestDecT36(t *testing.T) {
	tcs := [][]string{
		[]string{"0", "0", DecT36(0)},
		[]string{"9", "9", DecT36(9)},
		[]string{"a", "10", DecT36(10)},
		[]string{"b", "11", DecT36(11)},
		[]string{"10", "36", DecT36(36)},
		[]string{"11", "37", DecT36(37)},
		[]string{"100", "1296", DecT36(1296)},
		[]string{"500", "6480", DecT36(6480)},
		[]string{"1000", "46656", DecT36(46656)},
		[]string{"zzzz", "1679615", DecT36(1679615)},
		[]string{"zzzzzz", "2176782335", DecT36(2176782335)},
		[]string{"zzzzzzz", "78364164095", DecT36(78364164095)},
	}
	dd := unit.StrSingLine(tcs, "%s != {[DecT36(%s)] => %s}")
	if s, isStrig := dd.(string); isStrig {
		t.Fatal(s)
	} else if success, isBool := dd.(bool); isBool && !success {
		t.Fail()
	}
}

func TestDecT62(t *testing.T) {
	tcs := [][]string{
		[]string{"0", "0", DecT62(0)},
		[]string{"9", "9", DecT62(9)},
		[]string{"a", "10", DecT62(10)},
		[]string{"b", "11", DecT62(11)},
		[]string{"A", "36", DecT62(36)},
		[]string{"B", "37", DecT62(37)},
		[]string{"Z", "61", DecT62(61)},
		[]string{"10", "62", DecT62(62)},
		[]string{"ZZ", "3843", DecT62(3843)},
		[]string{"100", "3844", DecT62(3844)},
		[]string{"ZZZ", "238327", DecT62(238327)},
		[]string{"ZZZZ", "14776335", DecT62(14776335)},
		[]string{"ZZZZZ", "916132831", DecT62(916132831)},
		[]string{"ZZZZZZ", "56800235583", DecT62(56800235583)},           // 6
		[]string{"ZZZZZZZZ", "218340105584895", DecT62(218340105584895)}, // 8
		//[]string{"ZZZZZZZZZZ", "839299365868340264", DecT62(839299365868340264)}, // 10
	}
	dd := unit.StrSingLine(tcs, "%s != {[DecT62(%s)] => %s}")
	if s, isStrig := dd.(string); isStrig {
		t.Fatal(s)
	} else if success, isBool := dd.(bool); isBool && !success {
		t.Fail()
	}
}
