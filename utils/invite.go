package utils

const (
	PRIME1 = 3         // 与字符集长度 62 互质
	PRIME2 = 3         // 与邀请码长度 6 互质
	SALT   = 123456789 // 随意取一个数值
)

var AlphanumericSet = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	//'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// GetInvCodeByUIDUniqueNew 获取指定长度的邀请码。
func GetInvCodeByUIDUniqueNew(uid uint64, l int) string {
	// 放大 + 加盐。
	uid = uid*PRIME1 + SALT

	var code []rune
	slIdx := make([]byte, l)

	// 扩散。
	for i := 0; i < l; i++ {
		slIdx[i] = byte(uid % uint64(len(AlphanumericSet)))                   // 获取 62 进制的每一位值
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(len(AlphanumericSet)) // 其他位与个位加和再取余（让个位的变化影响到所有位）
		uid = uid / uint64(len(AlphanumericSet))                              // 相当于右移一位（62进制）
	}

	// 混淆。
	for i := 0; i < l; i++ {
		idx := (byte(i) * PRIME2) % byte(l)
		code = append(code, AlphanumericSet[slIdx[idx]])
	}
	return string(code)
}
