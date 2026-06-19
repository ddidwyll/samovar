package mqtt

import (
	"ergo.services/ergo/gen"

	natiu "github.com/soypat/natiu-mqtt"

	"context"
	"errors"
	"time"
)

type listener struct {
	gen.MetaProcess
	client *natiu.Client
}

func newListener(nc *natiu.Client) gen.MetaBehavior {
	return &listener{client: nc}
}

func (l *listener) Init(process gen.MetaProcess) error {
	l.MetaProcess = process
	return nil
}

func (l *listener) Start() error {
	for {
		// err := l.withTimeout(func(ctx context.Context) error {
		// 	return l.client.Ping(ctx)
		// })
		// if err == nil && !l.client.IsConnected() {
		// 	err = errors.New("not connected")
		// }
		// if err != nil {
		// 	l.Log().Error("mqtt.listener error: %s", err)
		// 	return err
		// }
		// if !l.client.IsConnected() {
		// 	return errors.New("not connected")
		// } else {
		// l.Log().Info("connected")
		// }
		l.Send(l.Parent(), "ping")
		// err := l.client.StartPing()
  	// ok := l.client.AwaitingPingresp()
  	// if err == nil && ok {
  	if true {
  	  l.client.HandleNext()
  		time.Sleep(100 * time.Millisecond)
  	} else {
    	return errors.New("mqtt connection failed")
  	}
	// for {
 //  	err := l.client.HandleNext()
	// 	l.Send(l.Parent(), "ping")
	// 	if err != nil {
	// 		l.Log().Error("mqtt.listener error: %s", err)
	// 		return err
	// 	}
	// 	time.Sleep(time.Second)
	// }
	// 	// err := l.withTimeout(func(ctx context.Context) error {
	// 	// 	return l.client.Ping(ctx)
	// 	// })
	// 	// if err != nil {
	// 	// 	l.Log().Error("mqtt.listener error: %s", err)
	// 	// 	return err
	// 	// } else {
 //  // 		l.Send(l.Parent(), "ping")
 //  // 		time.Sleep(time.Second)
	// 	// }
	// // }
	}
}

func (l *listener) withTimeout(f func(context.Context) error) error {
	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	err := f(ctx)
	cancel()
	return err
}

func (l *listener) Terminate(reason error) {
	l.Log().Debug("mqtt.listener terminated (%s)", reason)
}

func (l *listener) HandleCall(_ gen.PID, _ gen.Ref, request any) (any, error) {
	l.Log().Debug("mqtt.listener receive call request: %v", request)
	return nil, nil

}
func (l *listener) HandleMessage(_ gen.PID, msg any) error {
	l.Log().Debug("mqtt.listener receive message: %s", msg)
	return nil
}

func (l *listener) HandleInspect(_ gen.PID, items ...string) map[string]string {
	l.Log().Debug("mqtt.listener receive inspect request: %v", items)
	return nil
}
