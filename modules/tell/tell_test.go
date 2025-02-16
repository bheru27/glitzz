package tell

import (
	"github.com/bheru27/glitzz/config"
	"github.com/bheru27/glitzz/core"
	"github.com/bheru27/glitzz/tests"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTell(t *testing.T) {
	// setup
	tmpDirName, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Could not create a temporary directory: %s", err)
	}
	defer os.Remove(tmpDirName)

	conf := config.Default()
	conf.Tell.TellFile = filepath.Join(tmpDirName, "tell_data.json")

	sender := &tests.SenderMock{}

	m, err := New(sender, conf)
	if err != nil {
		t.Fatalf("Could not create a module: %s", err)
	}

	// test
	output, err := m.RunCommand(core.Command{Text: ".tell nick message", Nick: "author"})
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}

	e := irc.Event{Nick: "othernick", Arguments: []string{"any message"}, Code: "PRIVMSG"}
	m.HandleEvent(&e)
	if len(sender.Replies) != 0 {
		t.Fatalf("Invalid output length: %d", len(sender.Replies))
	}

	e = irc.Event{Nick: "nIcK", Arguments: []string{"any message"}, Code: "PRIVMSG"}
	m.HandleEvent(&e)
	if len(sender.Replies) != 1 {
		t.Fatalf("Invalid output length: %d", len(sender.Replies))
	}

}
