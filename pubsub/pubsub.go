package pubsub

type operation int

const (
	_ = iota
	sub
	unsub
	pub
	size
	stop
)

type cmdT struct {
	op   operation
	ch   chan interface{}
	rch  <-chan interface{}
	msg  interface{}
	done chan<- bool
}

type PubSub struct {
	control chan cmdT
	subs    []chan interface{}
}

func New() *PubSub {
	return (&PubSub{}).Start()
}

func (this *PubSub) Start() *PubSub {
	if this.control != nil {
		return this
	}

	this.control = make(chan cmdT, 128)

	go func() {
		var stopCmd *cmdT

	loop:
		for {
			select {
			case cmd := <-this.control:
				switch cmd.op {
				case sub:
					this.subs = append(this.subs, cmd.ch)
					cmd.done <- true
				case unsub:
					for i, test := range this.subs {
						if test == cmd.rch {
							this.subs = append(this.subs[:i], this.subs[i+1:]...)
							break
						}
					}
					cmd.done <- true
				case pub:
					for _, ch := range this.subs {
						ch <- cmd.msg
					}
					cmd.done <- true
				case size:
					cmd.ch <- len(this.subs)
				case stop:
					stopCmd = &cmd
					break loop
				}
			}
		}

		this.control = nil
		stopCmd.done <- true
	}()

	return this
}

func (this *PubSub) Sub() (<-chan interface{}, chan bool) {
	ch := make(chan interface{}, 128)
	done := make(chan bool, 1)

	this.control <- cmdT{
		op:   sub,
		ch:   ch,
		done: done,
	}

	return ch, done
}

func (this *PubSub) UnSub(ch <-chan interface{}) <-chan bool {
	done := make(chan bool, 1)

	this.control <- cmdT{
		op:   unsub,
		rch:  ch,
		done: done,
	}

	return done
}

func (this *PubSub) Pub(msg interface{}) <-chan bool {
	done := make(chan bool, 1)

	this.control <- cmdT{
		op:   pub,
		msg:  msg,
		done: done,
	}

	return done
}

func (this *PubSub) Size() int {
	resp := make(chan interface{})

	this.control <- cmdT{
		op: size,
		ch: resp,
	}

	return (<-resp).(int)
}

func (this *PubSub) Stop() <-chan bool {
	done := make(chan bool, 1)

	this.control <- cmdT{
		op:   stop,
		done: done,
	}

	return done
}
