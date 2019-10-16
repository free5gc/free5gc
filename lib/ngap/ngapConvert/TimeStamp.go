//go:binary-only-package

package ngapConvert

import (
	"encoding/binary"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
)

/*
RFC 5905 Section 6 https://tools.ietf.org/html/rfc5905#section-6

       0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                            Seconds                            |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                            Fraction                           |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                             NTP Timestamp Format

   The 64-bit timestamp format is used in packet headers and other
   places with limited word size.  It includes a 32-bit unsigned seconds
   field spanning 136 years and a 32-bit fraction field resolving 232
   picoseconds.  The 32-bit short format is used in delay and dispersion
   header fields where the full resolution and range of the other
   formats are not justified.  It includes a 16-bit unsigned seconds
   field and a 16-bit fraction field.

   In the date and timestamp formats, the prime epoch, or base date of
   era 0, is 0 h 1 January 1900 UTC, when all bits are zero.  It should
   be noted that strictly speaking, UTC did not exist prior to 1 January
   1972, but it is convenient to assume it has existed for all eternity,
   even if all knowledge of historic leap seconds has been lost.  Dates
   are relative to the prime epoch; values greater than zero represent

*/
func TimeStampToInt32(timeStampNgap aper.OctetString) (timeStamp int32) {}

func TimeStampToNgap(timeStamp int32) (timeStampNgap aper.OctetString) {}
