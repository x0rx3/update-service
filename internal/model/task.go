package model

type Task struct {
	server         *Server
	meta           map[string]any
	processLogChan chan *ProcessLog
}

func NewTask(server *Server, processLog chan *ProcessLog) *Task {
	return &Task{
		server:         server,
		meta:           make(map[string]any),
		processLogChan: processLog,
	}
}

func (inst *Task) Meta(key string) (any, bool) {
	if val, ok := inst.meta[key]; ok {
		return val, ok
	}
	return nil, false
}

func (inst *Task) AddMeta(key string, val any) {
	inst.meta[key] = val
}

func (inst *Task) DeleteMeta(key string) {
	delete(inst.meta, key)
}

func (inst *Task) Server() *Server {
	return inst.server
}

func (inst *Task) SendProcessLog(log *ProcessLog) {
	if inst.processLogChan != nil {
		inst.processLogChan <- log
	}
}
