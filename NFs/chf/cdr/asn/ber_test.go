package asn

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type intStruct struct {
	A int `ber:"tagNum:0"`
}
type twoIntStruct struct {
	A int `ber:"tagNum:0"`
	B int `ber:"tagNum:1"`
}
type nestedStruct struct {
	A intStruct `ber:"tagNum:0,set"`
	B intStruct `ber:"tagNum:1,seq"`
}
type choiceTest struct {
	Present int
	A       *int       `ber:"tagNum:0"`
	B       *BitString `ber:"tagNum:1"`
	C       *intStruct `ber:"tagNum:2,seq"`
	D       *int       `ber:"tagNum:32"`
	E       *int       `ber:"tagNum:128"`
}
type choiceInStruct struct {
	A int        `ber:"tagNum:0"`
	B choiceTest `ber:"tagNum:1,choice"`
}
type intSlice struct {
	List []int
}
type sliceInStruct struct {
	A []int `ber:"tagNum:0,seq"`
}

var i int

// TODO: slice test
func TestMarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		in    interface{}
		out   string
		param string
	}{
		{"intTest1", 10, "02010a", ""},
		{"intTest2", 127, "02017f", ""},
		{"intTest3", 128, "02020080", ""},
		{"intTest6", 0, "020100", ""},
		{"boolTest1", true, "0101ff", ""},
		{"boolTest2", false, "010100", ""},
		{"BitStringTest1", BitString{[]byte{0x80}, 1}, "03020780", ""},
		{"BitStringTest2", BitString{[]byte{0x81, 0xf0}, 12}, "03030481f0", ""},
		{"OctetStringTest1", OctetString([]byte{1, 2, 3}), "0403010203", ""},
		{"StringTest1", "test", "0c0474657374", "utf8"},
		{
			"StringTest2",
			"" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // This is 127 times 'x'
			"0c7f" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"78787878787878787878787878787878787878787878787878787878787878",
			"utf8",
		},
		{
			"StringTest3",
			"" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // This is 128 times 'x'
			"0c8180" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878",
			"utf8",
		},
		{"enumTest1", Enumerated(127), "0a017f", ""},
		{"enumTest2", Enumerated(128), "0a020080", ""},
		{
			"structTest1",
			intStruct{64},
			"3003" + "800140",
			"seq",
		},
		{
			"structTest2",
			twoIntStruct{64, 65},
			"3006" + "800140" + "810141",
			"seq",
		},
		{
			"structTest3",
			nestedStruct{intStruct{64}, intStruct{65}},
			"300a" +
				"a003" +
				"800140" +
				"a103" +
				"800141",
			"seq",
		},
		{
			"choiceTest1",
			choiceTest{
				1, &i, nil, nil, nil, nil,
			},
			"800100",
			"choice",
		},
		{
			"choiceTest2",
			choiceTest{
				2,
				nil,
				&BitString{[]byte{0x80}, 1},
				nil,
				nil,
				nil,
			},
			"81020780",
			"choice",
		},
		{
			"choiceTest3",
			choiceTest{
				3,
				nil,
				nil,
				&intStruct{64},
				nil,
				nil,
			},
			"a203" + "800140",
			"choice",
		},
		{
			"choiceTest4",
			choiceTest{
				4, nil, nil, nil, &i, nil,
			},
			"9f200100",
			"choice",
		},
		{
			"choiceTest5",
			choiceTest{
				5, nil, nil, nil, nil, &i,
			},
			"9f81000100",
			"choice",
		},
		{
			"choiceTest6",
			choiceInStruct{
				1,
				choiceTest{
					1, &i, nil, nil, nil, nil,
				},
			},
			"3008" + "800101" + "a103" + "800100",
			"seq",
		},
		{
			"sliceTest1",
			[]int{1, 2, 3},
			"3009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest2",
			[]intStruct{{1}, {2}, {3}},
			"300f" + "3003800101" + "3003800102" + "3003800103",
			"seq",
		},
		{
			"sliceTest3",
			intSlice{[]int{1, 2, 3}},
			"3009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest4",
			sliceInStruct{[]int{1, 2, 3}},
			"300b" + "a009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest5",
			[]int{},
			"3000",
			"seq",
		},
		{
			"sliceTest6",
			[]choiceTest{
				{
					1, &i, nil, nil, nil, nil,
				},
				{
					3, nil, nil, &intStruct{64}, nil, nil,
				},
			},
			"3008" + "800100" + "a203" + "800140",
			"seq",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := BerMarshalWithParams(tc.in, tc.param)
			require.NoError(t, err)
			require.Equal(t, tc.out, hex.EncodeToString(out))
		})
	}
}

func newInt(i int) *int                         { return &i }
func newString(s string) *string                { return &s }
func newBool(b bool) *bool                      { return &b }
func newEnum(i Enumerated) *Enumerated          { return &i }
func newOctetString(o OctetString) *OctetString { return &o }

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		out   interface{}
		in    string
		param string
	}{
		{"intTest1", newInt(10), "02010a", ""},
		{"intTest2", newInt(127), "02017f", ""},
		{"intTest3", newInt(128), "02020080", ""},
		{"intTest6", newInt(0), "020100", ""},
		{"boolTest1", newBool(true), "0101ff", ""},
		{"boolTest2", newBool(false), "010100", ""},
		{"BitStringTest1", &BitString{[]byte{0x80}, 1}, "03020780", ""},
		{"BitStringTest2", &BitString{[]byte{0x81, 0xf0}, 12}, "03030481f0", ""},
		{"OctetStringTest1", newOctetString([]byte{1, 2, 3}), "0403010203", ""},
		{"StringTest1", newString("test"), "0c0474657374", "utf8"},
		{
			"StringTest2",
			newString(
				"" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"), // This is 127 times 'x'
			"0c7f" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"78787878787878787878787878787878787878787878787878787878787878",
			"utf8",
		},
		{
			"StringTest3",
			newString("" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"), // This is 128 times 'x'
			"0c8180" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878" +
				"7878787878787878787878787878787878787878787878787878787878787878",
			"utf8",
		},
		{"enumTest1", newEnum(127), "0a017f", ""},
		{"enumTest2", newEnum(128), "0a020080", ""},
		{
			"structTest1",
			&intStruct{64},
			"3003" + "800140",
			"seq",
		},
		{
			"structTest2",
			&twoIntStruct{64, 65},
			"3006" + "800140" + "810141",
			"seq",
		},
		{
			"structTest3",
			&nestedStruct{intStruct{64}, intStruct{65}},
			"300a" +
				"a003" +
				"800140" +
				"a103" +
				"800141",
			"seq",
		},
		{
			"choiceTest1",
			&choiceTest{
				1, &i, nil, nil, nil, nil,
			},
			"800100",
			"choice",
		},
		{
			"choiceTest2",
			&choiceTest{
				2,
				nil,
				&BitString{[]byte{0x80}, 1},
				nil,
				nil,
				nil,
			},
			"81020780",
			"choice",
		},
		{
			"choiceTest3",
			&choiceTest{
				3,
				nil,
				nil,
				&intStruct{64},
				nil,
				nil,
			},
			"a203" + "800140",
			"choice",
		},
		{
			"choiceTest4",
			&choiceTest{
				4, nil, nil, nil, &i, nil,
			},
			"9f200100",
			"choice",
		},
		{
			"choiceTest5",
			&choiceTest{
				5, nil, nil, nil, nil, &i,
			},
			"9f81000100",
			"choice",
		},
		{
			"choiceTest6",
			&choiceInStruct{
				1,
				choiceTest{
					1, &i, nil, nil, nil, nil,
				},
			},
			"3008" + "800101" + "a103" + "800100",
			"seq",
		},
		{
			"sliceTest1",
			&[]int{1, 2, 3},
			"3009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest2",
			&[]intStruct{{1}, {2}, {3}},
			"300f" + "3003800101" + "3003800102" + "3003800103",
			"seq",
		},
		{
			"sliceTest3",
			&intSlice{[]int{1, 2, 3}},
			"3009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest4",
			&sliceInStruct{[]int{1, 2, 3}},
			"300b" + "a009" + "020101" + "020102" + "020103",
			"seq",
		},
		{
			"sliceTest5",
			&[]int{},
			"3000",
			"seq",
		},
		{
			"sliceTest6",
			&[]choiceTest{
				{
					1, &i, nil, nil, nil, nil,
				},
				{
					3, nil, nil, &intStruct{64}, nil, nil,
				},
			},
			"3008" + "800100" + "a203" + "800140",
			"seq",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, err := hex.DecodeString(tc.in)
			require.NoError(t, err)
			out := reflect.New(reflect.TypeOf(tc.out).Elem())
			val := out.Interface()
			err = UnmarshalWithParams(in, val, tc.param)
			require.NoError(t, err)
			// require.Equal(t, tc.out, val)
			require.True(t, reflect.DeepEqual(tc.out, val))
		})
	}
}

func TestParseInt64(t *testing.T) {
	testCases := [][]byte{}
	origInts := []int64{0, 1, 127, 128, 32767}

	for _, origInt := range origInts {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, int64(origInt))
		require.NoError(t, err)
		testCases = append(testCases, buf.Bytes())

		buf = new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, int32(origInt))
		require.NoError(t, err)
		testCases = append(testCases, buf.Bytes())

		buf = new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, int16(origInt))
		require.NoError(t, err)
		testCases = append(testCases, buf.Bytes())
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%x", tc), func(t *testing.T) {
			r, err := parseInt64(tc)
			require.NoError(t, err)
			require.Equal(t, origInts[i/3], r)
		})
	}
}
