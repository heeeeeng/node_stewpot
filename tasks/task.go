package tasks

type TaskType int

const (
	TaskTypeBroadcast         TaskType = 0
	TaskTypeConnRecv          TaskType = 1
	TaskTypeMsgProcess        TaskType = 2
	TaskTypeMsgProcessCPUReq  TaskType = 3
	TaskTypeMsgProcessCalc    TaskType = 4
	TaskTypeMsgTransmitUpload TaskType = 5
	TaskTypeMsgTransmitDelay  TaskType = 6
)

func (t TaskType) String() string {
	switch t {
	case TaskTypeBroadcast:
		return "TaskTypeBroadcast"
	case TaskTypeConnRecv:
		return "TaskTypeConnRecv"
	case TaskTypeMsgProcess:
		return "TaskTypeMsgProcess"
	case TaskTypeMsgProcessCPUReq:
		return "TaskTypeMsgProcessCPUReq"
	case TaskTypeMsgProcessCalc:
		return "TaskTypeMsgProcessCalc"
	case TaskTypeMsgTransmitUpload:
		return "TaskTypeMsgTransmitUpload"
	case TaskTypeMsgTransmitDelay:
		return "TaskTypeMsgTransmitDelay"

	default:
		return "unknown type"
	}
}
