package entity

type (
	QueueState uint
)

const (
	Wait     QueueState = 0
	Init     QueueState = 1
	Running  QueueState = 2
	Sleeping QueueState = 3
	Idle     QueueState = 4
	Stop     QueueState = 5
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
	}
	return ""
}

func (state QueueState) Check() bool {
	switch state {
	case Wait, Init, Running, Sleeping, Idle, Stop:
		return true
	default:
		return false
	}
}

func (state QueueState) Int() uint {
	return uint(state)
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
	}
	return ""
}
