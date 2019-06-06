package peterstd

import (
	"encoding"
	"gitlab.wallstcn.com/spider/peterstd/json"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
)

// check if err is not nil, panic
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Datetime

const (
	DateLayout     = "2006-01-02"
	DatetimeLayout = "2006-01-02 15:04:05"

	Key       = 0x01010101
	id2       = 113
	openBigId = true
)

// Return date string of today
func Today() string {
	return time.Now().Format(DateLayout)
}

// Date string trans to time.Time
// eg. DateStr2Time("2018-07-23")
func DateStr2Time(date string) (time.Time, error) {
	return time.ParseInLocation(DateLayout, date, time.UTC)
}

// Datetime string trans to time.Time
// eg. DatetimeStr2Time("2018-07-23 13:00:00")
func DatetimeStr2Time(date string) (time.Time, error) {
	return time.ParseInLocation(DatetimeLayout, date, time.UTC)
}

// Date string trans to unix timestamp
// eg. DateStr2Unix("2018-07-23")
func DateStr2Unix(date string) (int64, error) {
	t, err := DateStr2Time(date)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// Datetime string trans to unix timestamp
// eg. DatetimeStr2Unix("2018-07-23 13:00:00")
func DatetimeStr2Unix(date string) (int64, error) {
	t, err := DateStr2Time(date)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// Float64 to string
func FormatFloat(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

// ParseFloat string to float64
func ParseFloat(s string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		Error(err)
	}
	return res
}

// Int64 to string
func FormatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

// String to int64
func ParseInt64(s string) int64 {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		WithError(err).WithField("string", s).Error("Error in parse string to int64")
	}
	return res
}

// Json Translate
func JSONTranslate(s interface{}, d interface{}) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, d)
}

func ScanStruct(s, d interface{}, fields ...string) {
	var (
		tagname  = "json"
		sv       = reflect.Indirect(reflect.ValueOf(s))
		dv       = reflect.Indirect(reflect.ValueOf(d))
		st       = sv.Type()
		dt       = dv.Type()
		slen     = st.NumField()
		dlen     = dt.NumField()
		smap     = make(map[string]reflect.Value)
		fieldMap = make(map[string]bool)
		flen     = len(fields)
	)

	for i := 0; i < slen; i++ {
		stag := st.Field(i).Tag.Get(tagname)
		stag = strings.SplitN(stag, ",", 2)[0]
		smap[stag] = sv.Field(i)
	}

	for _, field := range fields {
		fieldMap[field] = true
	}

	for i := 0; i < dlen; i++ {
		dfield := dt.Field(i)
		dtag := dfield.Tag.Get(tagname)
		dtag = strings.SplitN(dtag, ",", 2)[0]
		if _, ok := fieldMap[dtag]; flen > 0 && !ok {
			continue
		}
		if svalue, ok := smap[dtag]; ok && svalue.Type() == dfield.Type {
			dv.Field(i).Set(svalue)
		}
	}
}

func StructFieldsMap(s interface{}, fields ...string) map[string]interface{} {
	var (
		res     = make(map[string]interface{})
		sv      = reflect.Indirect(reflect.ValueOf(s))
		st      = sv.Type()
		slen    = st.NumField()
		tagname = "json"
	)

	for i := 0; i < slen; i++ {
		if st.Field(i).Anonymous {
			submap := StructFieldsMap(sv.Field(i).Interface(), fields...)
			for k, v := range submap {
				res[k] = v
			}
		} else {
			tag := st.Field(i).Tag.Get(tagname)
			tag = strings.SplitN(tag, ",", 2)[0]
			res[tag] = sv.Field(i).Interface()
		}
	}

	if len(fields) == 0 {
		return res
	}

	fieldsMap := make(map[string]bool)
	for _, field := range fields {
		fieldsMap[field] = true
	}

	// Remove not need fields
	for k := range res {
		if ok := fieldsMap[k]; !ok {
			delete(res, k)
		}
	}

	// Add null fields
	for _, field := range fields {
		if _, ok := res[field]; !ok {
			res[field] = nil
		}
	}
	return res
}

func StrictFieldsMap(s interface{}, fields ...string) (map[string]interface{}, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	fieldsMap := make(map[string]bool)
	for _, field := range fields {
		fieldsMap[field] = true
		if _, ok := res[field]; !ok {
			res[field] = nil
		}
	}

	for k := range res {
		if _, ok := fieldsMap[k]; !ok {
			delete(res, k)
		}
	}
	return res, nil
}

// Get function name of fn
func GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// GenerateRequestID
func GenerateRequestID() string {
	return random.String(32)
}

// Check if any of errors is not nil, then panic
func Check(errors ...error) {
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}

// ScanRedisHash scan redis hash type to golang struct
func ScanRedisHash(data map[string]string, result interface{}) {
	var (
		s        = reflect.ValueOf(result).Elem()
		t        = reflect.TypeOf(result).Elem()
		numField = t.NumField()
		timeKind = reflect.TypeOf(time.Time{}).Kind()
	)

	for i := 0; i < numField; i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("json")
		name := strings.SplitN(tag, ",", 2)[0]
		if !ok {
			continue
		}
		v, ok := data[name]
		if !ok {
			continue
		}

		val := s.Field(i)
		switch val.Type().Kind() {
		case reflect.Int:
			tmp, _ := strconv.ParseInt(v, 10, 32)
			val.SetInt(int64(tmp))
		case reflect.Int64:
			tmp, _ := strconv.ParseInt(v, 10, 64)
			val.SetInt(tmp)
		case reflect.Float64:
			tmp, _ := strconv.ParseFloat(v, 64)
			val.SetFloat(tmp)
		case reflect.String:
			val.SetString(v)
		case timeKind:
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				val.Set(reflect.ValueOf(t))
			}
		case reflect.Struct:
			if m, ok := val.Interface().(encoding.BinaryUnmarshaler); ok {
				m.UnmarshalBinary([]byte(v))
			}
		}
	}
}

// RFC3339 to unix timestamp
func RFC3339ToUnix(datetime string) int64 {
	if datetime == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		Errorln(err)
	}
	return t.Unix()
}

// Error
func DeferError(err error, msg string) {
	if err != nil {
		WithError(err).Error(msg)
	}
}

func Round(f float64, digits int) float64 {
	shift := math.Pow(10, float64(digits))
	v := f * shift
	v1, v2 := math.Modf(v)
	if v1 < 0 {
		if v2 <= -0.5 {
			v1 -= 1
		}
	} else {
		if v2 >= 0.5 {
			v1 += 1
		}
	}
	return v1 / shift
}

func ContainsString(strs []string, target string) bool {
	for _, str := range strs {
		if str == target {
			return true
		}
	}
	return false
}

func ContainsInt(nums []int, target int) bool {
	for _, num := range nums {
		if num == target {
			return true
		}
	}
	return false
}

func ContainsInt64(nums []int64, target int64) bool {
	for _, num := range nums {
		if num == target {
			return true
		}
	}
	return false
}

// 大数变小数，长ID变短ID
func DecodePlateId(id int64) int64 {
	if !openBigId {
		return id
	}
	if id > 10000 {
		id = int64(math.Sqrt(float64((id^Key)/id2 + 1)))
	}
	return id
}

// 小数变大数，短ID变长ID
func EncodePlateId(id int64) int64 {
	if !openBigId {
		return id
	}
	if id < 10000 {
		id = ((id*id - 1) * id2) ^ Key
	}
	return id
}
