package sproto_test

import (
	"bytes"
	"testing"

	"github.com/xjdrew/gosproto"
)

type Test int

var inst Test = Test(5)

var barCalled bool

func (t *Test) Foobar(req *FoobarRequest, resp *FoobarResponse) {
	if t != &inst {
		panic(t)
	}
	resp.What = req.What
}

func (t *Test) Foo(resp *FooResponse) {
	if t != &inst {
		panic(t)
	}
	resp.Ok = sproto.Bool(true)
}

func (t *Test) Bar() {
	if t != &inst {
		panic(t)
	}
	barCalled = true
}

func TestFoobarService(t *testing.T) {
	name := "test.foobar"
	input := "hello"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := sproto.NewService(rw, protocols)
	req := FoobarRequest{
		What: &input,
	}
	call, err := client.Go(name, &req, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	server, _ := sproto.NewService(rw, protocols)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if err := client.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}
	<-call.Done
	resp := call.Resp.(*FoobarResponse)
	if resp.What == nil || *resp.What != input {
		t.Fatalf("unexpected response:%v", resp.What)
	}
}

func TestFooService(t *testing.T) {
	name := "test.foo"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := sproto.NewService(rw, protocols)
	call, err := client.Go(name, nil, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	server, _ := sproto.NewService(rw, protocols)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if err := client.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}
	<-call.Done
	resp := call.Resp.(*FooResponse)
	if resp.Ok == nil || !*resp.Ok {
		t.Fatalf("unexpected response:%v", resp.Ok)
	}
}

func TestBarService(t *testing.T) {
	name := "test.bar"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := sproto.NewService(rw, protocols)
	err := client.Invoke(name, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	barCalled = false
	server, _ := sproto.NewService(rw, protocols)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if !barCalled {
		t.Fatal("unexpected dispatch")
	}
}
