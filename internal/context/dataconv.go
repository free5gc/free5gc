package context

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"unicode"
)

// Ipv6ToInt - Convert Ipv6 string to *bigInt
func Ipv6ToInt(ipv6 string) *big.Int {
	ipv6 = ipv6 + "/32"
	ip, _, err := net.ParseCIDR(ipv6)
	if err != nil {
		fmt.Println("Error", ip, err)
	}
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ip)
	return IPv6Int
}

// Ipv4ToInt - Convert Ipv4 string to int64
func Ipv4ToInt(ipv4 string) int64 {
	ipv4 = ipv4 + "/24"
	ip, _, err := net.ParseCIDR(ipv4)
	if err != nil {
		fmt.Println("Error", ip, err)
	}
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(ip)
	return IPv4Int.Int64()
}

// Ipv4IntToIpv4String - Convert Ipv4 int64 to string
func Ipv4IntToIpv4String(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// Ipv6IntToIpv6String - Convert Ipv6 *big.Int to string
func Ipv6IntToIpv6String(ip *big.Int) string {
	ipv6Bytes := ip.Bytes()
	ipv6String := hex.EncodeToString(ipv6Bytes)

	for i := 1; i < 8; i++ {
		ipv6String = ipv6String[:i-1+4*i] + ":" + ipv6String[i-1+4*i:]
	}
	return ipv6String
}

// EncodeGroupId - Encode GroupId to number string(output pattern: [10][3][3][25])
func EncodeGroupId(groupId string) string {
	externalGroupIdentitySplit := strings.Split(groupId, "-")

	var encodedGroupId string

	var encodedGroupServiceIdentifier int
	var groupServiceIdentifierPadding string
	for _, c := range externalGroupIdentitySplit[0] {
		if unicode.IsNumber(c) {
			encodedGroupServiceIdentifier = encodedGroupServiceIdentifier*16 + int(c)
		} else {
			n := int(unicode.ToLower(c) - 'a')
			encodedGroupServiceIdentifier = encodedGroupServiceIdentifier*16 + n
		}
	}
	for i := 0; i < (10 - len(externalGroupIdentitySplit[0])); i++ {
		groupServiceIdentifierPadding += "0"
	}
	encodedGroupId = encodedGroupId + groupServiceIdentifierPadding + strconv.Itoa(encodedGroupServiceIdentifier)

	encodedGroupId = encodedGroupId + externalGroupIdentitySplit[1]

	if len(externalGroupIdentitySplit[2]) == 2 {
		encodedGroupId = encodedGroupId + "0" + externalGroupIdentitySplit[2]
	} else {
		encodedGroupId = encodedGroupId + externalGroupIdentitySplit[2]
	}

	var encodedLocalGroupId int
	var localGroupIdPadding string
	for _, c := range externalGroupIdentitySplit[3] {
		if unicode.IsNumber(c) {
			encodedLocalGroupId = encodedGroupServiceIdentifier*16 + int(c)
		} else {
			n := int(unicode.ToLower(c) - 'a')
			encodedLocalGroupId = encodedGroupServiceIdentifier*16 + n
		}
	}

	for i := 0; i < (25 - len(externalGroupIdentitySplit[3])); i++ {
		localGroupIdPadding += "0"
	}
	encodedGroupId = encodedGroupId + localGroupIdPadding + strconv.Itoa(encodedLocalGroupId)

	return encodedGroupId
}
