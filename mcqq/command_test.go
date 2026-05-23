package mcqq

import (
	"strings"
	"testing"
)

func TestCutField(t *testing.T) {
	field, rest := cutField(" Server   say hello ")
	if field != "Server" {
		t.Fatalf("expected field Server, got %q", field)
	}
	if rest != "say hello" {
		t.Fatalf("expected rest say hello, got %q", rest)
	}

	field, rest = cutField("Server")
	if field != "Server" || rest != "" {
		t.Fatalf("unexpected single field result: %q %q", field, rest)
	}
}

func TestPlainComponent(t *testing.T) {
	component := plainComponent("hello")
	if component.Text == nil || *component.Text != "hello" {
		t.Fatalf("unexpected component text: %#v", component.Text)
	}
	if component.Color == nil || *component.Color != White {
		t.Fatalf("unexpected component color: %#v", component.Color)
	}
}

func TestMCCommandHelp(t *testing.T) {
	help := mcCommandHelp()
	for _, expected := range []string{"/mc status", "/mc rcon", "/mc title", "/mc actionbar", "/mc private"} {
		if !strings.Contains(help, expected) {
			t.Fatalf("help should contain %q: %s", expected, help)
		}
	}
}
