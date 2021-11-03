package entity

type (
	QueueState uint
)

const (
	// Wait 等待
	Wait QueueState = 0
	// Ready 就绪态
	Ready QueueState = 1
	// Running 执行态
	Running QueueState = 2
	// Sleeping 休眠|暂停 态
	Sleeping QueueState = 3
	// Idle 空闲态
	Idle QueueState = 4
	// Stop 停止状态
	Stop QueueState = 5
	// ReStarting 重启态
	ReStarting QueueState = 6
)

func (state QueueState) Describe() string {
	switch state {
	case Wait:
		return "waiting for run"
	case Ready:
		return "ready for run"
	case Running:
		return "queue running"
	case Sleeping:
		return "queue sleeping"
	case Idle:
		return "queue idle"
	case Stop:
		return "queue stopped"
	case ReStarting:
		return "queue restarting"
	}
	return ""
}

func (state QueueState) Check() bool {
	switch state {
	case Wait, Ready, Running, Sleeping, Idle, Stop, ReStarting:
		return true
	default:
		return false
	}
}

func (state QueueState) Int() uint {
	return uint(state)
}

func (state QueueState) Is(s uint) bool {
	switch s {
	case Wait.Int(), Ready.Int(), Running.Int(), Sleeping.Int(), Idle.Int(), Stop.Int(), ReStarting.Int():
		return true
	}
	return false
}

func (state QueueState) String() string {
	switch state {
	case Wait:
		return "wait"
	case Ready:
		return "ready"
	case Running:
		return "running"
	case Sleeping:
		return "sleeping"
	case Idle:
		return "idle"
	case Stop:
		return "stopped"
	case ReStarting:
		return "restart"
	}
	return ""
}
