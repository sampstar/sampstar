package network

func bitsToBytes(bits int) int {
	return (bits + 7) >> 3
}

func bytesToBits(bytes int) int {
	return bytes << 3
}
