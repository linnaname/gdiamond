package subscriber

import (
	"gdiamond/client/listener"
	"sync"
)

var once sync.Once
var ins *Subscriber

//GetSubscriberInstance  keep it singleton
func GetSubscriberInstance() *Subscriber {
	once.Do(func() {
		ins, _ = newSubscriber(listener.NewDefaultSubscriberListener())

	})
	return ins
}
