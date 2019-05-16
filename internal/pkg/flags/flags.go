package flags

type Flag int

const (
	In Flag = iota
	Out
	_
	File
	Timeout
	Memory
	QEMU
	Dirty
	Dumb
	Dict
	Banner
	Exploration
)

func (f Flag) Int() int {
	return int(f)
}
