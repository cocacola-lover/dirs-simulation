package nlogger

func (l *Logger) changeIdForMessageWithoutLock(oldId, newId int) {
	l.startedSearches[newId] = l.startedSearches[oldId]
}

func (l *Logger) ChangeIdForMessages(oldIds, newIds []int) {
	l.dLock.Lock()
	for i, oldId := range oldIds {
		l.changeIdForMessageWithoutLock(oldId, newIds[i])
	}
	l.dLock.Unlock()
}
