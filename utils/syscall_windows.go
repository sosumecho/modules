package utils

const (
	SIGUSR1 = Integer(0x1e)
	SIGUSR2 = Integer(0x1f)

	LOCK_EX = int(0x2)
	LOCK_NB = int(0x4)
)

type Integer int

func (integer Integer) String() string {
	switch integer {
	case SIGUSR1:
		return "USR1"
	case SIGUSR2:
		return "USR2"
	default:
		return "UNDEFINED"
	}
}

func (integer Integer) Signal() {
	return
}

// Flock 文件加锁
func Flock(fd int, how int) error {
	return nil
}
