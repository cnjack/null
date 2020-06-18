package null

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)

	nullJSON          = []byte(`null`)
	invalidJSON       = []byte(`:)`)
	testStringBson, _ = bson.Marshal(bson.M{
		"key": "test",
	})
	nullBson, _ = bson.Marshal(bson.M{
		"key": nil,
	})
)

type stringInStruct struct {
	Test String `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assertStr(t, str, "StringFrom() string")

	zero := StringFrom("")
	if !zero.Valid {
		t.Error("StringFrom(0)", "is invalid, but should be valid")
	}
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	sptr := &s
	str := StringFromPtr(sptr)
	assertStr(t, str, "StringFromPtr() string")

	null := StringFromPtr(nil)
	assertNullStr(t, null, "StringFromPtr(nil)")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assertStr(t, str, "string json")

	var ns String
	err = json.Unmarshal(nullStringJSON, &ns)
	if err == nil {
		panic("err should not be nil")
	}

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullStr(t, null, "null json")

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := StringFrom("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	null := StringFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, `null`, "null json marshal")
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")
}

// Tests omitempty... broken until Go 1.4
// func TestMarshalStringInStruct(t *testing.T) {
// 	obj := stringInStruct{Test: StringFrom("")}
// 	data, err := json.Marshal(obj)
// 	maybePanic(err)
// 	assertJSONEquals(t, data, `{}`, "null string in struct")
// }

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := NewString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := StringFrom("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := NewString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestStringSetValid(t *testing.T) {
	change := NewString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func TestStringValueOrZero(t *testing.T) {
	valid := NewString("test", true)
	if valid.ValueOrZero() != "test" {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewString("test", false)
	if invalid.ValueOrZero() != "" {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestStringEqual(t *testing.T) {
	str1 := NewString("foo", false)
	str2 := NewString("foo", false)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", false)
	str2 = NewString("bar", false)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("foo", true)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("foo", false)
	assertStringEqualIsFalse(t, str1, str2)

	str1 = NewString("foo", false)
	str2 = NewString("foo", true)
	assertStringEqualIsFalse(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("bar", true)
	assertStringEqualIsFalse(t, str1, str2)
}

func TestStringEncodeValue(t *testing.T) {
	rb := RegisterNullStruct(bsoncodec.NewRegistryBuilder())
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.DefaultRegistry = rb.Build()
	out, err := bson.Marshal(bson.M{
		"key": NewString("test", true),
	})
	assert.NoError(t, err)
	out2, err := bson.Marshal(bson.M{
		"key": "test",
	})
	assert.NoError(t, err)
	assert.Equal(t, out2, out)
}

func TestStringNullEncodeValue(t *testing.T) {
	rb := RegisterNullStruct(bsoncodec.NewRegistryBuilder())
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.DefaultRegistry = rb.Build()
	out, err := bson.Marshal(bson.M{
		"key": NewString("xxxx", false),
	})
	assert.NoError(t, err)
	out2, err := bson.Marshal(bson.M{
		"key": nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, out2, out)
}

func TestStringNullEncodeValue1(t *testing.T) {
	rb := RegisterNullStruct(bsoncodec.NewRegistryBuilder())
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.DefaultRegistry = rb.Build()
	out, err := bson.Marshal(bson.M{
		"key": NewString("some invalid string", false),
	})
	assert.NoError(t, err)
	out2, err := bson.Marshal(bson.M{
		"key": nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, out2, out)
}

func TestStringNullDecodeValue(t *testing.T) {
	rb := RegisterNullStruct(bsoncodec.NewRegistryBuilder())
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.DefaultRegistry = rb.Build()
	strBson := &struct {
		Key String `bson:"key"`
	}{}
	err := bson.Unmarshal(testStringBson, strBson)
	if assert.NoError(t, err) {
		if assert.True(t, strBson.Key.Valid) {
			assert.Equal(t, "test", strBson.Key.String)
		}
	}
}

func TestStringNullDecodeValue2(t *testing.T) {
	rb := RegisterNullStruct(bsoncodec.NewRegistryBuilder())
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bson.DefaultRegistry = rb.Build()
	strBson := &struct {
		Key String `bson:"key"`
	}{}
	err := bson.Unmarshal(nullBson, strBson)
	if assert.NoError(t, err) {
		assert.False(t, strBson.Key.Valid)
	}
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStr(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func assertStringEqualIsTrue(t *testing.T, a, b String) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return true", a.String, a.Valid, b.String, b.Valid)
	}
}

func assertStringEqualIsFalse(t *testing.T, a, b String) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return false", a.String, a.Valid, b.String, b.Valid)
	}
}
