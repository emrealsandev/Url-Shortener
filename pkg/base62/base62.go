package base62

const DICTIONARY = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(n uint64) string {
	if n == 0 {
		return "0"
	}
	out := make([]byte, 0, 11)
	for n > 0 {
		out = append(out, DICTIONARY[n%62])
		n /= 62
	}
	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}
