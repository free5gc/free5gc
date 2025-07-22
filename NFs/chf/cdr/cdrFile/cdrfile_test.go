package cdrFile

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCdrFile(t *testing.T) {
	t.Parallel()

	cdrFile1 := CDRFile{
		Hdr: CdrFileHeader{
			FileLength:                            71,
			HeaderLength:                          63,
			HighReleaseIdentifier:                 7,
			HighVersionIdentifier:                 4,
			LowReleaseIdentifier:                  7,
			LowVersionIdentifier:                  5,
			FileOpeningTimestamp:                  CdrHdrTimeStamp{4, 28, 17, 18, 1, 8, 0},
			TimestampWhenLastCdrWasAppendedToFIle: CdrHdrTimeStamp{1, 2, 3, 4, 1, 6, 30},
			NumberOfCdrsInFile:                    1,
			FileSequenceNumber:                    11,
			FileClosureTriggerReason:              4,
			IpAddressOfNodeThatGeneratedFile: [20]byte{
				0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb, 0xa, 0xb,
			},
			LostCdrIndicator:               4,
			LengthOfCdrRouteingFilter:      4,
			CDRRouteingFilter:              []byte("abcd"),
			LengthOfPrivateExtension:       5,
			PrivateExtension:               []byte("fghjk"), // vendor specific
			HighReleaseIdentifierExtension: 2,
			LowReleaseIdentifierExtension:  3,
		},
		CdrList: []CDR{{
			Hdr: CdrHeader{
				CdrLength:                  3,
				ReleaseIdentifier:          BeyondRel9,                   // octet 3 bit 6..8
				VersionIdentifier:          3,                            // otcet 3 bit 1..5
				DataRecordFormat:           UnalignedPackedEncodingRules, // octet 4 bit 6..8
				TsNumber:                   TS32253,                      // octet 4 bit 1..5
				ReleaseIdentifierExtension: 4,
			},
			CdrByte: []byte("abc"),
		}},
	}

	cdrFile2 := CDRFile{
		Hdr: CdrFileHeader{
			FileLength:                            92,
			HeaderLength:                          66,
			HighReleaseIdentifier:                 7,
			HighVersionIdentifier:                 5,
			LowReleaseIdentifier:                  7,
			LowVersionIdentifier:                  6,
			FileOpeningTimestamp:                  CdrHdrTimeStamp{1, 2, 11, 56, 1, 7, 30},
			TimestampWhenLastCdrWasAppendedToFIle: CdrHdrTimeStamp{4, 3, 2, 1, 0, 4, 0},
			NumberOfCdrsInFile:                    3,
			FileSequenceNumber:                    65,
			FileClosureTriggerReason:              2,
			IpAddressOfNodeThatGeneratedFile: [20]byte{
				0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd, 0xc, 0xd,
			},
			LostCdrIndicator:               4,
			LengthOfCdrRouteingFilter:      5,
			CDRRouteingFilter:              []byte("gfdss"),
			LengthOfPrivateExtension:       7,
			PrivateExtension:               []byte("abcdefg"), // vendor specific
			HighReleaseIdentifierExtension: 1,
			LowReleaseIdentifierExtension:  2,
		},
		CdrList: []CDR{
			{
				Hdr: CdrHeader{
					CdrLength:                  3,
					ReleaseIdentifier:          BeyondRel9,
					VersionIdentifier:          3,
					DataRecordFormat:           UnalignedPackedEncodingRules,
					TsNumber:                   TS32253,
					ReleaseIdentifierExtension: 4,
				},
				CdrByte: []byte("abc"),
			},
			{
				Hdr: CdrHeader{
					CdrLength:                  6,
					ReleaseIdentifier:          BeyondRel9,
					VersionIdentifier:          2,
					DataRecordFormat:           AlignedPackedEncodingRules1,
					TsNumber:                   TS32205,
					ReleaseIdentifierExtension: 2,
				},
				CdrByte: []byte("ghjklm"),
			},
			{
				Hdr: CdrHeader{
					CdrLength:                  2,
					ReleaseIdentifier:          BeyondRel9,
					VersionIdentifier:          3,
					DataRecordFormat:           AlignedPackedEncodingRules1,
					TsNumber:                   TS32225,
					ReleaseIdentifierExtension: 1,
				},
				CdrByte: []byte("cv"),
			},
		},
	}

	testCases := []struct {
		name string
		in   CDRFile
	}{
		{"cdrfile1", cdrFile1},
		{"cdrfile2", cdrFile2},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileName := "encoding" + strconv.Itoa(i) + ".txt"
			tc.in.Encoding(fileName)
			newCdrFile := CDRFile{}
			newCdrFile.Decoding(fileName)
			e := os.Remove(fileName)
			if e != nil {
				fmt.Println(e)
			}

			require.Equal(t, tc.in, newCdrFile)
			// require.True(t, reflect.DeepEqual(tc.in, newCdrFile))
		})
	}
}
