package app

type ThemeOption struct {
	Key  string
	Name string
	Desc string
}

var themeOptions = []ThemeOption{
	{"rose-pine", "Rosé Pine", "Elegant rose tones"},
	{"catppuccin", "Catppuccin Mocha", "Soft pastels, gentle on the eyes"},
	{"nord", "Nord", "Cool Arctic blues, minimal"},
	{"tokyo-night", "Tokyo Night", "Vibrant cyberpunk vibes"},
	{"gruvbox", "Gruvbox", "Warm retro colors"},
	{"dracula", "Dracula", "Classic high contrast"},
}

func themeIndexByKey(key string) int {
	for i, opt := range themeOptions {
		if opt.Key == key {
			return i
		}
	}
	return 0
}
