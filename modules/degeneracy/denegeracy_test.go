package degeneracy

import (
	"github.com/bheru27/glitzz/config"
	"github.com/bheru27/glitzz/tests"
	"github.com/thoj/go-ircevent"
	"testing"
)

func TestShots(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	senderMock := &tests.SenderMock{}
	m, err := New(senderMock, config.Default())
	if err != nil {
		t.Fatalf("error creating module: %s", err)
	}
	e := &irc.Event{Code: "PRIVMSG", Arguments: []string{"#channel", ".shots nick!"}}

	m.HandleEvent(e)
	if len(senderMock.Replies) != 4 {
		t.Errorf("invalid reply: %+v", senderMock.Replies)
	}
}
