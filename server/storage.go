package server

type storage interface {
	execCmd(cmd cmd) (cmdResult, error)
}
