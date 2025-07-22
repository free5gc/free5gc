package util

import "hash/crc32"

func HashCRC32(text string) (uint32, error) {
	h := crc32.NewIEEE()
	_, err := h.Write([]byte(text))
	if err != nil {
		return 0, err
	}

	return h.Sum32(), nil
}
