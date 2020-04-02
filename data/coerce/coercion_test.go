package coerce

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCoerceToString(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToString(valInt)
	assert.Equal(t, "2", cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToString(valBool)
	assert.Equal(t, "true", cval, "not equal")

	var valStr interface{} = "12"
	cval, _ = ToString(valStr)
	assert.Equal(t, "12", cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToString(valFloat)
	assert.Equal(t, "1.23", cval, "not equal")

	var valNil interface{} // = nil
	cval, _ = ToString(valNil)
	assert.Equal(t, "", cval, "not equal")

	var pointerString interface{} = StringPointer("flogo")
	cval, _ = ToString(pointerString)
	assert.Equal(t, "flogo", cval, "not equal")
}

func StringPointer(s string) *string {
	return &s
}

func TestCoerceToInt(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt(valInt)
	assert.Equal(t, 2, cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToInt(valBool)
	assert.Equal(t, 1, cval, "not equal")

	var valStr interface{} = "12"
	cval, _ = ToInt(valStr)
	assert.Equal(t, 12, cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToInt(valFloat)
	assert.Equal(t, 1, cval, "not equal")

	var valNil interface{} //= nil
	cval, _ = ToInt(valNil)
	assert.Equal(t, 0, cval, "not equal")
}

func TestCoerceToInt32(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt32(valInt)
	assert.Equal(t, int32(2), cval)

	var valBool interface{} = true
	cval, _ = ToInt32(valBool)
	assert.Equal(t, int32(1), cval)

	var valStr interface{} = "12"
	cval, _ = ToInt32(valStr)
	assert.Equal(t, int32(12), cval)

	var valFloat interface{} = 1.23
	cval, _ = ToInt32(valFloat)
	assert.Equal(t, int32(1), cval)

	var valNil interface{} //= nil
	cval, _ = ToInt32(valNil)
	assert.Equal(t, int32(0), cval)
}

func TestCoerceToInt64(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToInt64(valInt)
	assert.Equal(t, int64(2), cval)

	var valBool interface{} = true
	cval, _ = ToInt64(valBool)
	assert.Equal(t, int64(1), cval)

	var valStr interface{} = "12"
	cval, _ = ToInt64(valStr)
	assert.Equal(t, int64(12), cval)

	var valFloat interface{} = 1.23
	cval, _ = ToInt64(valFloat)
	assert.Equal(t, int64(1), cval)

	var valNil interface{} //= nil
	cval, _ = ToInt64(valNil)
	assert.Equal(t, int64(0), cval)
}

func TestCoerceToBoolean(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToBool(valInt)
	assert.Equal(t, true, cval, "not equal")

	var valBool interface{} = true
	cval, _ = ToBool(valBool)
	assert.Equal(t, true, cval, "not equal")

	var valStr interface{} = "false"
	cval, _ = ToBool(valStr)
	assert.Equal(t, false, cval, "not equal")

	var valFloat interface{} = 1.23
	cval, _ = ToBool(valFloat)
	assert.Equal(t, true, cval, "not equal")

	var valNil interface{} //= nil
	cval, _ = ToBool(valNil)
	assert.Equal(t, false, cval, "not equal")
}

func TestCoerceToObject(t *testing.T) {

	s := struct {
		Id   int
		Name string
	}{Id: 1001,
		Name: "flogo"}

	obj, err := ToObject(s)
	assert.Nil(t, err)
	assert.Equal(t, "flogo", obj["Name"])
}

func TestCoerceArray(t *testing.T) {

	var valInt interface{} = 2
	cval, _ := ToArray(valInt)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	valArr := []string{"a", "b"}
	cval, _ = ToArray(valArr)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	valInfArr := []interface{}{"a", "b", 2, 5}
	cval, _ = ToArray(valInfArr)
	assert.Equal(t, reflect.Slice, reflect.ValueOf(cval).Kind())

	var valNilArr interface{}
	cval, _ = ToArray(valNilArr)
	assert.Nil(t, cval)

}

func TestParams(t *testing.T) {
	var err error
	simpleParams := map[string]string{"params": "a"}
	cval, _ := ToParams(simpleParams)
	assert.NotNil(t, cval)

	simpleStringParams := `{"params": "a"}`
	cval, _ = ToParams(simpleStringParams)
	assert.NotNil(t, cval["params"])

	simpleStringParams = `{"a"}`
	_, err = ToParams(simpleStringParams)
	assert.NotNil(t, err)

	simpleDefault := 3
	_, err = ToParams(simpleDefault)
	assert.NotNil(t, err)
}

func TestConnection(t *testing.T) {

	var err error

	connString := "sampleConn"
	_, err = ToConnection(connString)
	assert.NotNil(t, err)

}

func TestToDateTime(t *testing.T) {
	var examples = []string{
		"May 8, 2009 5:57:51 PM",
		"oct 7, 1970",
		"oct 7, '70",
		"oct. 7, 1970",
		"oct. 7, 70",
		"Mon Jan  2 15:04:05 2006",
		"Mon Jan  2 15:04:05 MST 2006",
		"Mon Jan 02 15:04:05 -0700 2006",
		"Monday, 02-Jan-06 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Thu, 4 Jan 2018 17:53:36 +0000",
		"Mon Aug 10 15:44:11 UTC+0100 2015",
		"Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
		"September 17, 2012 10:09am",
		"September 17, 2012 at 10:09am PST-08",
		"September 17, 2012, 10:10:09",
		"October 7, 1970",
		"October 7th, 1970",
		"12 Feb 2006, 19:17",
		"12 Feb 2006 19:17",
		"7 oct 70",
		"7 oct 1970",
		"03 February 2013",
		"1 July 2013",
		"2013-Feb-03",
		//   mm/dd/yy
		"3/31/2014",
		"03/31/2014",
		"08/21/71",
		"8/1/71",
		"4/8/2014 22:05",
		"04/08/2014 22:05",
		"4/8/14 22:05",
		"04/2/2014 03:00:51",
		"8/8/1965 12:00:00 AM",
		"8/8/1965 01:00:01 PM",
		"8/8/1965 01:00 PM",
		"8/8/1965 1:00 PM",
		"8/8/1965 12:00 AM",
		"4/02/2014 03:00:51",
		"03/19/2012 10:11:59",
		"03/19/2012 10:11:59.3186369",
		// yyyy/mm/dd
		"2014/3/31",
		"2014/03/31",
		"2014/4/8 22:05",
		"2014/04/08 22:05",
		"2014/04/2 03:00:51",
		"2014/4/02 03:00:51",
		"2012/03/19 10:11:59",
		"2012/03/19 10:11:59.3186369",
		// Chinese
		"2014年04月08日",
		//   yyyy-mm-ddThh
		"2006-01-02T15:04:05+0000",
		"2009-08-12T22:15:09-07:00",
		"2009-08-12T22:15:09",
		"2009-08-12T22:15:09Z",
		//   yyyy-mm-dd hh:mm:ss
		"2014-04-26 17:24:37.3186369",
		"2012-08-03 18:31:59.257000000",
		"2014-04-26 17:24:37.123",
		"2013-04-01 22:43",
		"2013-04-01 22:43:22",
		"2014-12-16 06:20:00 UTC",
		"2014-12-16 06:20:00 GMT",
		"2014-04-26 05:24:37 PM",
		"2014-04-26 13:13:43 +0800",
		"2014-04-26 13:13:43 +0800 +08",
		"2014-04-26 13:13:44 +09:00",
		"2012-08-03 18:31:59.257000000 +0000 UTC",
		"2015-09-30 18:48:56.35272715 +0000 UTC",
		"2015-02-18 00:12:00 +0000 GMT",
		"2015-02-18 00:12:00 +0000 UTC",
		"2015-02-08 03:02:00 +0300 MSK m=+0.000000001",
		"2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001",
		"2017-07-19 03:21:51+00:00",
		"2014-04-26",
		"2014-04",
		"2014",
		"2014-05-11 08:20:13,787",
		// mm.dd.yy
		"3.31.2014",
		"03.31.2014",
		"08.21.71",
		"2014.03",
		"2014.03.30",
		//  yyyymmdd and similar
		"20140601",
		"20140722105203",
		// unix seconds, ms, micro, nano
		"1332151919",
		"1384216367189",
		"1384216367111222",
		"1384216367111222333",
	}

	for _, dt := range examples {
		dateTime, err := ToDateTime(dt)
		assert.Nil(t, err)
		assert.NotNil(t, dateTime)
	}

}
