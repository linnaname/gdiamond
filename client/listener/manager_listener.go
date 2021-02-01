package listener

//ManagerListener watch config change
type ManagerListener interface {
	//
	ReceiveConfigInfo(configInfo string)
}

//ManagerListenerFunc have a idea but not implement yet
type ManagerListenerFunc func()

func (f ManagerListenerFunc) delegate() {
	// delegate to the anonymous function
	f()
}
