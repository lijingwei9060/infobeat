package command

// Commander 接口，用于执行调用
type Commander interface {
	RunCommand(...string) ([]byte, error)
}
