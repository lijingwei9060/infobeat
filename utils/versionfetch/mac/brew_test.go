package mac

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lijingwei9060/infobeat/utils/command"
)

func TestFetchBrewVersion(t *testing.T) {
	commander, err := command.New(command.Config{TimeOut: 5 * time.Second, Backoff: 5})
	assert.NotNil(t, commander)
	assert.Nil(t, err)

	out, err := FetchBrewVersion(commander)
	assert.NotNil(t, out)
	assert.Nil(t, err)
	for _, str := range out {
		t.Log(str)
	}
}
