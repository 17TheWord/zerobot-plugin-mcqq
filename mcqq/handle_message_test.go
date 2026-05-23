package mcqq

import "testing"

func TestTranslateTextFallback(t *testing.T) {
	if got := translateText(TranslateModel{Text: "translated", Key: "key"}); got != "translated" {
		t.Fatalf("expected translated text, got %q", got)
	}
	if got := translateText(TranslateModel{Key: "fallback.key"}); got != "fallback.key" {
		t.Fatalf("expected key fallback, got %q", got)
	}
}

func TestDeathTextFallback(t *testing.T) {
	if got := deathText(DeathModel{Text: "Steve was slain", Key: "death.attack.mob"}); got != "Steve was slain" {
		t.Fatalf("expected death text, got %q", got)
	}
	if got := deathText(DeathModel{Key: "death.attack.mob"}); got != "death.attack.mob" {
		t.Fatalf("expected death key fallback, got %q", got)
	}
}

func TestAchievementTextFallback(t *testing.T) {
	achievement := AchievementModel{
		Key:  "minecraft:story/mine_stone",
		Text: "legacy achievement text",
		Translate: TranslateModel{
			Text: "translated achievement text",
		},
	}
	if got := achievementText(achievement); got != "translated achievement text" {
		t.Fatalf("expected translated achievement text, got %q", got)
	}

	achievement.Translate = TranslateModel{}
	if got := achievementText(achievement); got != "legacy achievement text" {
		t.Fatalf("expected legacy achievement text, got %q", got)
	}

	achievement.Text = ""
	if got := achievementText(achievement); got != "minecraft:story/mine_stone" {
		t.Fatalf("expected achievement key fallback, got %q", got)
	}
}

func TestFlexibleTranslateTextUnmarshal(t *testing.T) {
	var text FlexibleTranslateText
	if err := text.UnmarshalJSON([]byte(`"plain title"`)); err != nil {
		t.Fatalf("unmarshal plain text: %v", err)
	}
	if got := text.String(); got != "plain title" {
		t.Fatalf("expected plain title, got %q", got)
	}

	text = FlexibleTranslateText{}
	if err := text.UnmarshalJSON([]byte(`{"key":"advancement.key","text":"Translated title"}`)); err != nil {
		t.Fatalf("unmarshal translate text: %v", err)
	}
	if got := text.String(); got != "Translated title" {
		t.Fatalf("expected translated title, got %q", got)
	}
}
