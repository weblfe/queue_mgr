package entity

type (
	QueueState uint
)

const (
	Wait       QueueState = 0
	Init       QueueState = 1
	Running    QueueState = 2
	Sleeping   QueueState = 3
	Idle       QueueState = 4
	Stop       QueueState = 5
	ReStarting QueueState = 6
)

func (state QueueState) Describe() string {
	switch state {
	case Wait:
		return "waiting for run"
	case Init:
		return "initial for run"
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
	case Wait, Init, Running, Sleeping, Idle, Stop, ReStarting:
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
	case Wait.Int(), Init.Int(), Running.Int(), Sleeping.Int(), Idle.Int(), Stop.Int(), ReStarting.Int():
		return true
	}
	return false
}

func (state QueueState) String() string {
	switch state {
	case Wait:
		return "wait"
	case Init:
		return "init"
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
