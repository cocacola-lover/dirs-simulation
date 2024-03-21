package nlogger

func (l *Logger) CountRouteMessageReceives() int {
	l.rmrLock.Lock()
	defer l.rmrLock.Unlock()

	ans := 0
	for _, emap := range l.routeMessageReceives {
		for _, earr := range emap {
			ans += len(earr)
		}
	}

	return ans
}

func (l *Logger) CountRouteMessageTimeouts() int {
	l.rmtLock.Lock()
	defer l.rmtLock.Unlock()

	ans := 0
	for _, emap := range l.routeMessageTimeouts {
		for _, earr := range emap {
			ans += len(earr)
		}
	}

	return ans
}

func (l *Logger) CountRouteMessageConfirms() int {
	l.rmcLock.Lock()
	defer l.rmcLock.Unlock()

	ans := 0
	for _, emap := range l.routeMessageConfirms {
		for _, earr := range emap {
			ans += len(earr)
		}
	}

	return ans
}

func (l *Logger) CountDownloadMessages() int {
	l.dLock.Lock()
	defer l.dLock.Unlock()

	ans := 0
	for _, earr := range l.downloadMessages {
		ans += len(earr)
	}

	return ans
}
