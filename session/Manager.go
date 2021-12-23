package session

type Manager struct {
	Sessions map[string]map[string]string
}

func (man *Manager) Contains(id string) bool {
	return man.Sessions[id] != nil
}

//var globalSessions *Manager
