package pubsub

type operation int

const (
	_ = iota
	sub
	unsub
	pub
)

type cmd struct {
	op  operation
	ch  chan string
	rch <-chan string
	msg string
}

type PubSub struct {
	ch   chan cmd
	subs []chan string
}

func New() *PubSub {
	p := &PubSub{}
	p.ch = make(chan cmd)
	go p.start()
	return p
}

func (p *PubSub) start() {
	for {
		select {
		case cmd := <-p.ch:
			switch cmd.op {
			case sub:
				p.subs = append(p.subs, cmd.ch)
			case unsub:
				for i, test := range p.subs {
					if test == cmd.rch {
						p.subs = append(p.subs[:i], p.subs[i+1:]...)
						break
					}
				}
			case pub:
				for _, ch := range p.subs {
					ch <- cmd.msg
				}
			}
		}
	}
}

func (p *PubSub) Sub() <-chan string {
	ch := make(chan string, 128)

	p.ch <- cmd{
		op: sub,
		ch: ch,
	}

	return ch
}

func (p *PubSub) UnSub(ch <-chan string) {
	p.ch <- cmd{
		op:  unsub,
		rch: ch,
	}
}

func (p *PubSub) Pub(msg string) {
	p.ch <- cmd{
		op:  pub,
		msg: msg,
	}
}
