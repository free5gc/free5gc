package cdrFile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type CDRFile struct {
	Hdr     CdrFileHeader
	CdrList []CDR
}

type CDR struct {
	Hdr     CdrHeader
	CdrByte []byte
}

// TS 32.297 6.1.1 CDR file header format
type CdrFileHeader struct {
	FileLength                            uint32
	HeaderLength                          uint32
	HighReleaseIdentifier                 uint8 // octet 9 bit 6..8
	HighVersionIdentifier                 uint8 // octet 9 bit 1..5
	LowReleaseIdentifier                  uint8 // octet 10 bit 6..8
	LowVersionIdentifier                  uint8 // octet 10 bit 1..5
	FileOpeningTimestamp                  CdrHdrTimeStamp
	TimestampWhenLastCdrWasAppendedToFIle CdrHdrTimeStamp
	NumberOfCdrsInFile                    uint32
	FileSequenceNumber                    uint32
	FileClosureTriggerReason              FileClosureTriggerReasonType
	IpAddressOfNodeThatGeneratedFile      [20]byte // ip address in ipv6 format
	LostCdrIndicator                      uint8
	LengthOfCdrRouteingFilter             uint16
	CDRRouteingFilter                     []byte // vendor specific
	LengthOfPrivateExtension              uint16
	PrivateExtension                      []byte // vendor specific
	HighReleaseIdentifierExtension        uint8
	LowReleaseIdentifierExtension         uint8
}

type CdrHeader struct {
	CdrLength                  uint16
	ReleaseIdentifier          ReleaseIdentifierType // octet 3 bit 6..8
	VersionIdentifier          uint8                 // otcet 3 bit 1..5
	DataRecordFormat           DataRecordFormatType  // octet 4 bit 6..8
	TsNumber                   TsNumberIdentifier    // octet 4 bit 1..5
	ReleaseIdentifierExtension uint8
}

type CdrHdrTimeStamp struct {
	MonthLocal  uint8
	DateLocal   uint8
	HourLocal   uint8
	MinuteLocal uint8

	// bit set to "1" expresses "+" or bit set to "0" expresses "-" time deviation)
	SignOfTheLocalTimeDifferentialFromUtc uint8
	HourDeviation                         uint8
	MinuteDeviation                       uint8
}

type FileClosureTriggerReasonType uint8

const (
	NormalClosure                     FileClosureTriggerReasonType = 0
	FileSizeLimitReached              FileClosureTriggerReasonType = 1
	FileOpentimeLimitedReached        FileClosureTriggerReasonType = 2
	MaximumNumberOfCdrsInFileReached  FileClosureTriggerReasonType = 3
	FileClosedByManualIntervention    FileClosureTriggerReasonType = 4
	CdrReleaseVersionOrEncodingChange FileClosureTriggerReasonType = 5
	AbnormalFileClosure               FileClosureTriggerReasonType = 128
	FileSystemError                   FileClosureTriggerReasonType = 129
	FileSystemStorageExhausted        FileClosureTriggerReasonType = 130
	FileIntegrityError                FileClosureTriggerReasonType = 131
)

type ReleaseIdentifierType uint8

const (
	Rel99 ReleaseIdentifierType = iota
	Rel4
	Rel5
	Rel6
	Rel7
	Rel8
	Rel9
	BeyondRel9
)

type DataRecordFormatType uint8

const (
	BasicEncodingRules DataRecordFormatType = iota + 1
	UnalignedPackedEncodingRules
	AlignedPackedEncodingRules1
	XMLEncodingRules
)

type TsNumberIdentifier uint8

const (
	TS32005 TsNumberIdentifier = 0
	TS32015 TsNumberIdentifier = 1
	TS32205 TsNumberIdentifier = 2
	TS32215 TsNumberIdentifier = 3
	TS32225 TsNumberIdentifier = 4
	TS32235 TsNumberIdentifier = 5
	TS32250 TsNumberIdentifier = 6
	TS32251 TsNumberIdentifier = 7
	TS32260 TsNumberIdentifier = 9
	TS32270 TsNumberIdentifier = 10
	TS32271 TsNumberIdentifier = 11
	TS32272 TsNumberIdentifier = 12
	TS32273 TsNumberIdentifier = 13
	TS32275 TsNumberIdentifier = 14
	TS32274 TsNumberIdentifier = 15
	TS32277 TsNumberIdentifier = 16
	TS32296 TsNumberIdentifier = 17
	TS32278 TsNumberIdentifier = 18
	TS32253 TsNumberIdentifier = 19
	TS32255 TsNumberIdentifier = 20
	TS32254 TsNumberIdentifier = 21
	TS32256 TsNumberIdentifier = 22
	TS28201 TsNumberIdentifier = 23
	TS28202 TsNumberIdentifier = 24
)

func (cdrf CdrFileHeader) Encoding() []byte {
	buf := new(bytes.Buffer)

	// File length
	if err := binary.Write(buf, binary.BigEndian, cdrf.FileLength); err != nil {
		fmt.Println("CdrFileHeader File length failed:", err)
	}

	// Header length
	if err := binary.Write(buf, binary.BigEndian, cdrf.HeaderLength); err != nil {
		fmt.Println("CdrFileHeader Header length failed:", err)
	}

	// High release / version identifier

	// if cdrf.HighReleaseIdentifier == 7 {
	// var highIdentifier uint8 =
	//  (cdrf.HighReleaseIdentifier+cdrf.HighReleaseIdentifierExtension+1)*100 + cdrf.HighVersionIdentifier
	// 	if err := binary.Write(buf, binary.BigEndian, highIdentifier); err != nil {
	// 		fmt.Println("binary.Write failed:", err)
	// 	}
	// } else {
	// 	var highIdentifier uint8 = cdrf.HighReleaseIdentifier*100 + cdrf.HighVersionIdentifier
	// 	if err := binary.Write(buf, binary.BigEndian, highIdentifier); err != nil {
	// 		fmt.Println("binary.Write failed:", err)
	// 	}
	// }

	highIdentifier := (cdrf.HighReleaseIdentifier << 5) | cdrf.HighVersionIdentifier

	if err := binary.Write(buf, binary.BigEndian, highIdentifier); err != nil {
		fmt.Println("CdrFileHeader highIdentifier failed:", err)
	}

	// Low release / version identifier

	// if cdrf.LowReleaseIdentifier == 7 {
	// 	var lowIdentifier uint8 =
	//   (cdrf.LowReleaseIdentifier+cdrf.LowReleaseIdentifierExtension+1)*100 + cdrf.LowVersionIdentifier
	// 	if err := binary.Write(buf, binary.BigEndian, lowIdentifier); err != nil {
	// 		fmt.Println("binary.Write failed:", err)
	// 	}
	// } else {
	// 	var lowIdentifier uint8 = cdrf.LowReleaseIdentifier*100 + cdrf.LowVersionIdentifier
	// 	if err := binary.Write(buf, binary.BigEndian, lowIdentifier); err != nil {
	// 		fmt.Println("binary.Write failed:", err)
	// 	}
	// }

	lowIdentifier := (cdrf.LowReleaseIdentifier << 5) | cdrf.LowVersionIdentifier

	if err := binary.Write(buf, binary.BigEndian, lowIdentifier); err != nil {
		fmt.Println("CdrFileHeader lowIdentifier failed:", err)
	}

	// File opening timestamp
	// _, offset := cdrf.FileOpeningTimestamp.Zone()
	// sign := 0
	// offsetHour := -offset/3600
	// offsetMin := -offset/60%60

	// if offset >= 0 {
	// 	sign = 1
	// 	offsetHour = offset/3600
	// 	offsetMin = offset/60%60
	// }

	// var ts uint32 = uint32(cdrf.FileOpeningTimestamp.Month())<<28 |
	// 	uint32(cdrf.FileOpeningTimestamp.Day())<<23 |
	// 	uint32(cdrf.FileOpeningTimestamp.Hour())<<18 |
	// 	uint32(cdrf.FileOpeningTimestamp.Minute())<<12 |
	// 	uint32(sign)<<11 |
	// 	uint32(offsetHour)<<6 |
	// 	uint32(offsetMin)

	ts := uint32(cdrf.FileOpeningTimestamp.MonthLocal)<<28 |
		uint32(cdrf.FileOpeningTimestamp.DateLocal)<<23 |
		uint32(cdrf.FileOpeningTimestamp.HourLocal)<<18 |
		uint32(cdrf.FileOpeningTimestamp.MinuteLocal)<<12 |
		uint32(cdrf.FileOpeningTimestamp.SignOfTheLocalTimeDifferentialFromUtc)<<11 |
		uint32(cdrf.FileOpeningTimestamp.HourDeviation)<<6 |
		uint32(cdrf.FileOpeningTimestamp.MinuteDeviation)

	if err := binary.Write(buf, binary.BigEndian, ts); err != nil {
		fmt.Println("CdrFileHeader File opening timestamp failed:", err)
	}

	// Timestamp when last CDR was appended to file
	// _, offset = cdrf.TimestampWhenLastCdrWasAppendedToFIle.Zone()
	// sign = 0
	// offsetHour = -offset/3600
	// offsetMin = -offset/60%60

	// if offset >= 0 {
	// 	sign = 1
	// 	offsetHour = offset/3600
	// 	offsetMin = offset/60%60
	// }

	// ts  = uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.Month())<<28 |
	// 	uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.Day())<<23 |
	// 	uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.Hour())<<18 |
	// 	uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.Minute())<<12 |
	// 	uint32(sign)<<11 |
	// 	uint32(offsetHour)<<6 |
	// 	uint32(offsetMin)
	ts = uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.MonthLocal)<<28 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.DateLocal)<<23 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.HourLocal)<<18 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.MinuteLocal)<<12 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.SignOfTheLocalTimeDifferentialFromUtc)<<11 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.HourDeviation)<<6 |
		uint32(cdrf.TimestampWhenLastCdrWasAppendedToFIle.MinuteDeviation)

	if err := binary.Write(buf, binary.BigEndian, ts); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Number of CDRs in file
	if err := binary.Write(buf, binary.BigEndian, cdrf.NumberOfCdrsInFile); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// File sequence number
	if err := binary.Write(buf, binary.BigEndian, cdrf.FileSequenceNumber); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// File closure trigger reason
	if err := binary.Write(buf, binary.BigEndian, cdrf.FileClosureTriggerReason); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Node IP address
	if err := binary.Write(buf, binary.BigEndian, cdrf.IpAddressOfNodeThatGeneratedFile); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Lost CDR indicator
	if err := binary.Write(buf, binary.BigEndian, cdrf.LostCdrIndicator); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Length of CDR routeing filter
	if err := binary.Write(buf, binary.BigEndian, cdrf.LengthOfCdrRouteingFilter); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// CDR routeing filter
	if err := binary.Write(buf, binary.BigEndian, cdrf.CDRRouteingFilter); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Length of private extension
	if err := binary.Write(buf, binary.BigEndian, cdrf.LengthOfPrivateExtension); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// Private extension
	if err := binary.Write(buf, binary.BigEndian, cdrf.PrivateExtension); err != nil {
		fmt.Println("CdrFileHeader failed:", err)
	}

	// "High Release Identifier" extension
	if cdrf.HighReleaseIdentifier == 7 {
		if err := binary.Write(buf, binary.BigEndian, cdrf.HighReleaseIdentifierExtension); err != nil {
			fmt.Println("CdrFileHeader failed:", err)
		}
	}

	// "Low Release Identifier" extension
	if cdrf.LowReleaseIdentifier == 7 {
		if err := binary.Write(buf, binary.BigEndian, cdrf.LowReleaseIdentifierExtension); err != nil {
			fmt.Println("CdrFileHeader failed:", err)
		}
	}

	// fmt.Printf("Encoded: % b\n", buf.Bytes())
	// fmt.Printf("%#v\n", sign)
	// fmt.Printf("%#v\n", offsetHour)
	// fmt.Printf("%#v\n", offsetMin)
	// fmt.Printf("%#v\n", cdrf.FileOpeningTimestamp)
	// fmt.Printf("%#v", cdrf.FileOpeningTimestamp.UTC().Sub(*cdrf.FileOpeningTimestamp).Hours())

	// Check
	if cdrf.HeaderLength != uint32(len(buf.Bytes())) && cdrf.HeaderLength != 0xffffffff {
		fmt.Println("[Encoding Warning]HeaderLength field of CdfFile Header not equals to the length of CdfFile Header.")
		fmt.Println("\tExpected", uint32(len(buf.Bytes())), "Get", cdrf.HeaderLength)
	}

	return buf.Bytes()
}

func (header CdrHeader) Encoding() []byte {
	buf := new(bytes.Buffer)

	// CDR length
	if err := binary.Write(buf, binary.BigEndian, header.CdrLength); err != nil {
		fmt.Println("CdrHeader failed:", err)
	}

	// Release/Version Identifier
	identifier := uint8(header.ReleaseIdentifier)<<5 |
		uint8(header.VersionIdentifier)

	if err := binary.Write(buf, binary.BigEndian, identifier); err != nil {
		fmt.Println("CdrHeader failed:", err)
	}

	// Data Record Format / TS number
	oct4 := uint8(header.DataRecordFormat)<<5 | uint8(header.TsNumber)

	if err := binary.Write(buf, binary.BigEndian, oct4); err != nil {
		fmt.Println("CdrHeader failed:", err)
	}

	// Release Identifier extension
	if header.ReleaseIdentifier == 7 {
		if err := binary.Write(buf, binary.BigEndian, header.ReleaseIdentifierExtension); err != nil {
			fmt.Println("CdrHeader failed:", err)
		}
	}

	// fmt.Printf("Encoded: % b\n", buf.Bytes())
	return buf.Bytes()
}

func (cdfFile CDRFile) Encoding(fileName string) {
	buf := new(bytes.Buffer)

	// Cdr File Header
	bufCdrFileHeader := cdfFile.Hdr.Encoding()
	if err := binary.Write(buf, binary.BigEndian, bufCdrFileHeader); err != nil {
		fmt.Println("CDRFile failed:", err)
	}

	for i, cdr := range cdfFile.CdrList {
		bufCdrHeader := cdr.Hdr.Encoding()
		if err := binary.Write(buf, binary.BigEndian, bufCdrHeader); err != nil {
			fmt.Println("CDRFile failed:", err)
		}

		if err := binary.Write(buf, binary.BigEndian, cdr.CdrByte); err != nil {
			fmt.Println("CDRFile failed:", err)
		}

		if len(cdr.CdrByte) != int(cdr.Hdr.CdrLength) {
			fmt.Println("[Encoding Warning]CdrLength field of cdr", i, "header not equals to the length of encoding cdr", i)
			fmt.Println("\tExpected", len(cdr.CdrByte), "Get", int(cdr.Hdr.CdrLength))
		}
	}

	// Check
	if cdfFile.Hdr.FileLength != uint32(len(buf.Bytes())) && cdfFile.Hdr.FileLength != 0xffffffff {
		fmt.Println("[Encoding Warning]FileLength field of CdfFile Header not equals to the length of encoding file.")
		fmt.Println("\tExpected", uint32(len(buf.Bytes())), "Get", cdfFile.Hdr.FileLength)
	}

	// fmt.Printf("Encoded: %b\n", buf.Bytes())
	err := os.WriteFile(fileName, buf.Bytes(), 0o666)
	if err != nil {
		panic(err)
	}
}

func (cdfFile *CDRFile) Decoding(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// fileLength := binary.BigEndian.Uint32(data[0:4])

	// File opening timestamp
	ts := binary.BigEndian.Uint32(data[10:14])
	// month := ts >> 28
	// date := int((ts >> 23) & 0b11111)
	// hour := int((ts >> 18) & 0b11111)
	// minute := int((ts >> 12) & 0b111111)
	// sign := (ts >> 11) & 0b1
	// hourDeviation := (ts >> 6) & 0b11111
	// minuteDeviation := ts & 0b111111

	// var offset int
	// if sign == 1 {
	// 	offset = int(hourDeviation*3600 + minuteDeviation*60)
	// } else if sign == 0 {
	// 	offset = int(hourDeviation*3600 + minuteDeviation*60) * -1
	// }
	// loc := time.FixedZone("", offset)
	// // The year is temporarily set to the current year
	// fileOpeningTimestamp := time.Date(time.Now().Year(), time.Month(month), date, hour, minute, 0, 0, loc)
	fileOpeningTimestamp := CdrHdrTimeStamp{
		MonthLocal:                            uint8(ts >> 28),
		DateLocal:                             uint8((ts >> 23) & 0b11111),
		HourLocal:                             uint8((ts >> 18) & 0b11111),
		MinuteLocal:                           uint8((ts >> 12) & 0b111111),
		SignOfTheLocalTimeDifferentialFromUtc: uint8((ts >> 11) & 0b1),
		HourDeviation:                         uint8((ts >> 6) & 0b11111),
		MinuteDeviation:                       uint8(ts & 0b111111),
	}

	// Last CDR append timestamp
	ts = binary.BigEndian.Uint32(data[14:18])
	// month = ts >> 28
	// date = int((ts >> 23) & 0b11111)
	// hour = int((ts >> 18) & 0b11111)
	// minute = int((ts >> 12) & 0b111111)
	// sign = (ts >> 11) & 0b1
	// hourDeviation = (ts >> 6) & 0b11111
	// minuteDeviation = ts & 0b111111

	// if sign == 1 {
	// 	offset = int(hourDeviation*3600 + minuteDeviation*60)
	// } else if sign == 0 {
	// 	offset = int(hourDeviation*3600 + minuteDeviation*60) * -1
	// }
	// loc = time.FixedZone("", offset)
	// // The year is temporarily set to the current year
	// lastCDRAppendTimestamp := time.Date(time.Now().Year(), time.Month(month), date, hour, minute, 0, 0, loc)
	// fmt.Println(fileLength)

	lastCDRAppendTimestamp := CdrHdrTimeStamp{
		MonthLocal:                            uint8(ts >> 28),
		DateLocal:                             uint8((ts >> 23) & 0b11111),
		HourLocal:                             uint8((ts >> 18) & 0b11111),
		MinuteLocal:                           uint8((ts >> 12) & 0b111111),
		SignOfTheLocalTimeDifferentialFromUtc: uint8((ts >> 11) & 0b1),
		HourDeviation:                         uint8((ts >> 6) & 0b11111),
		MinuteDeviation:                       uint8(ts & 0b111111),
	}

	// Length
	numberOfCdrsInFile := binary.BigEndian.Uint32(data[18:22])
	lengthOfCdrRouteingFilter := binary.BigEndian.Uint16(data[48:50])
	xy := 50 + lengthOfCdrRouteingFilter
	LengthOfPrivateExtension := binary.BigEndian.Uint16(data[xy : xy+2])
	n := xy + 2 + LengthOfPrivateExtension

	// ip
	var IpAddressOfNodeThatGeneratedFile [20]byte
	copy(IpAddressOfNodeThatGeneratedFile[:], data[27:47])

	cdfFile.Hdr = CdrFileHeader{
		FileLength:                            binary.BigEndian.Uint32(data[0:4]),
		HeaderLength:                          binary.BigEndian.Uint32(data[4:8]),
		HighReleaseIdentifier:                 data[8] >> 5,
		HighVersionIdentifier:                 data[8] & 0b11111,
		LowReleaseIdentifier:                  data[9] >> 5,
		LowVersionIdentifier:                  data[9] & 0b11111,
		FileOpeningTimestamp:                  fileOpeningTimestamp,
		TimestampWhenLastCdrWasAppendedToFIle: lastCDRAppendTimestamp,
		NumberOfCdrsInFile:                    numberOfCdrsInFile,
		FileSequenceNumber:                    binary.BigEndian.Uint32(data[22:26]),
		FileClosureTriggerReason:              FileClosureTriggerReasonType(data[26]),
		IpAddressOfNodeThatGeneratedFile:      IpAddressOfNodeThatGeneratedFile,
		LostCdrIndicator:                      data[47],
		LengthOfCdrRouteingFilter:             lengthOfCdrRouteingFilter,
		CDRRouteingFilter:                     data[50:xy],
		LengthOfPrivateExtension:              LengthOfPrivateExtension,
		PrivateExtension:                      data[xy+2 : n],
	}

	tail := uint32(n)

	if cdfFile.Hdr.HighReleaseIdentifier == 7 {
		cdfFile.Hdr.HighReleaseIdentifierExtension = data[n]
		tail++
	}
	if cdfFile.Hdr.LowReleaseIdentifier == 7 {
		cdfFile.Hdr.LowReleaseIdentifierExtension = data[n+1]
		tail++
	}

	// fmt.Println("[Decode]cdrfileheader:\n", cdfFile.Hdr)

	for i := 1; i <= int(numberOfCdrsInFile); i++ {
		cdrLength := binary.BigEndian.Uint16(data[tail : tail+2])
		releaseIdentifier := ReleaseIdentifierType(data[tail+2] >> 5)
		expectedLength := int(tail) + 4 + int(cdrLength)
		if releaseIdentifier == 7 {
			expectedLength += 1
		}
		if len(data) < expectedLength {
			fmt.Println("[Decoding Error]Length of cdrfile is wrong. cdr:", i)
		}

		cdrHeader := CdrHeader{
			CdrLength:         cdrLength,
			ReleaseIdentifier: ReleaseIdentifierType(data[tail+2] >> 5),
			VersionIdentifier: data[tail+2] & 0b11111,
			DataRecordFormat:  DataRecordFormatType(data[tail+3] >> 5),
			TsNumber:          TsNumberIdentifier(data[tail+3] & 0b11111),
		}

		i := uint32(4)
		if cdrHeader.ReleaseIdentifier == 7 {
			cdrHeader.ReleaseIdentifierExtension = data[tail+4]
			i++
		}

		cdr := CDR{
			Hdr:     cdrHeader,
			CdrByte: data[tail+i : tail+i+uint32(cdrLength)],
		}
		cdfFile.CdrList = append(cdfFile.CdrList, cdr)
		tail += i + uint32(cdrLength)
	}
	// fmt.Println("[Decode]cdrfile:\n", cdfFile)
	// fmt.Printf("%#v\n", cdfFile)
}

// func main() {
// 	loc, _ := time.LoadLocation("Asia/Kolkata")
// 	timeNow := time.Now().In(loc)
// 	// timeNow := time.Now()
// 	cdrf := CdrFileHeader{
// 		FileLength:                            5,
// 		HeaderLength:                          6,
// 		HighReleaseIdentifier:                 2,
// 		HighVersionIdentifier:                 3,
// 		LowReleaseIdentifier:                  4,
// 		LowVersionIdentifier:                  5,
// 		FileOpeningTimestamp:                  &timeNow,
// 		TimestampWhenLastCdrWasAppendedToFIle: &timeNow,
// 		NumberOfCdrsInFile:                    1,
// 		FileSequenceNumber:                    11,
// 		FileClosureTriggerReason:              4,
// 		//IpAddressOfNodeThatGeneratedFile      [20]byte(),
// 		LostCdrIndicator:          4,
// 		LengthOfCdrRouteingFilter: 4,
// 		CDRRouteingFilter:                     []byte("abcd"),
// 		LengthOfPrivateExtension: 5,
// 		PrivateExtension:                      []byte("fghjk"), // vendor specific
// 		HighReleaseIdentifierExtension: 2,
// 		LowReleaseIdentifierExtension:  3,
// 	}

// 	cdrHeader := CdrHeader {
// 		CdrLength                  :3,
// 		ReleaseIdentifier          :Rel6, // octet 3 bit 6..8
// 		VersionIdentifier          :3,                // otcet 3 bit 1..5
// 		DataRecordFormat           :UnalignedPackedEncodingRules,  // octet 4 bit 6..8
// 		TsNumber                   : TS32253,   // octet 4 bit 1..5
// 		ReleaseIdentifierExtension :4,
// 	}

// 	cdrFile := CDRFile{
// 		hdr: cdrf,
// 		cdrList: []CDR{{hdr:cdrHeader, cdrByte:[]byte("abc")},},
// 	}

// 	cdrFile.Encoding()
// 	cdrFile = CDRFile{}
// 	cdrFile.Decoding("encoding.txt")
// }
