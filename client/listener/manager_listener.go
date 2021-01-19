package listener

type ManagerListener interface {
	ReceiveConfigInfo(configInfo string)
}

type ManagerListenerFunc func()

func (f ManagerListenerFunc) delegate() {
	// delegate to the anonymous function
	f()
}
