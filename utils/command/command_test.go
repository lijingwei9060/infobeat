package command

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	command, err := New(Config{TimeOut: 30 * time.Second, Backoff: 3})
	assert.NotNil(t, command)
	assert.Nil(t, err)

	out, err := command.RunCommand("ls")

	assert.NotNil(t, out)
	assert.Nil(t, err)
	t.Log(string(out))
}

func TestCommandWithTimeout(t *testing.T) {
	command, err := New(Config{TimeOut: 3 * time.Second, Backoff: 3})
	assert.NotNil(t, command)
	assert.Nil(t, err)

	out, err := command.RunCommand("sleep", "5")

	assert.NotNil(t, err)
	t.Log(string(out))
	t.Log(err)
}

func TestPing(t *testing.T) {
	command, err := New(Config{TimeOut: 3 * time.Second, Backoff: 3})
	assert.NotNil(t, command)
	assert.Nil(t, err)

	out, err := command.RunCommand("rpm", "-qa", "--queryformat", "'%{name}:%{version}.%{release}\n'")

	assert.NotNil(t, err)
	t.Log(string(out))
	t.Log(err)
}

func TestConcurrency(t *testing.T) {

	command, err := New(Config{TimeOut: 5 * time.Second, Backoff: 3})
	assert.NotNil(t, command)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}

	sleepFun := func() {
		defer wg.Done()
		now := time.Now()
		_, err1 := command.RunCommand("sleep", "3")
		t.Logf("time:%s\n", time.Now().Sub(now))
		assert.Nil(t, err1)
	}

	wg.Add(1)
	go sleepFun()
	wg.Add(1)
	go sleepFun()
	wg.Add(1)
	go sleepFun()
	wg.Add(1)
	go sleepFun()
	wg.Add(1)
	go sleepFun()

	wg.Wait()
}
