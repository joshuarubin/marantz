package pubsub

import "testing"

func TestNew(t *testing.T) {
	var p *PubSub

	if p != nil {
		t.Error("p not initially nil")
	}

	p = New()

	if p == nil {
		t.Error("p is nil")
	}

	if p.control == nil {
		t.Error("p.control is nil")
	}

	if p.Size() != 0 {
		t.Error("p.Size is not 0")
	}
}

func TestStop(t *testing.T) {
	p := New()

	<-p.Stop()

	if p.control != nil {
		t.Error("p did not stop")
	}
}

func TestSub(t *testing.T) {
	p := New()

	ch, done := p.Sub()
	<-done

	if ch == nil {
		t.Error("ch is nil")
	}

	if p.Size() != 1 {
		t.Error("p.Size is not 1")
	}

	<-p.UnSub(ch)

	if p.Size() != 0 {
		t.Error("p.Size is not 0")
	}
}

func pubTester(t *testing.T, recv <-chan interface{}, done chan bool) {
	val := <-recv

	str, ok := val.(string)
	if !ok {
		t.Error("val not a string")
	}

	if str != "test" {
		t.Error("str did not match 'test'")
	}

	done <- true
}

func TestPub(t *testing.T) {
	p := New()
	done := make(chan bool)

	const numSubs = 1
	for i := 0; i < numSubs; i++ {
		recv, sync := p.Sub()
		<-sync
		go pubTester(t, recv, done)
	}

	p.Pub("test")

	for i := 0; i < numSubs; i++ {
		<-done
	}
}
