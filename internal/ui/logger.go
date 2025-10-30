package ui

import (
	"fmt"
	"strings"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Title(msg string) {
	fmt.Println()
	fmt.Println(TitleStyle.Render("✨ " + msg))
}

func (l *Logger) Subtitle(msg string) {
	fmt.Println(SubtitleStyle.Render(msg))
}

func (l *Logger) Success(msg string) {
	fmt.Println(SuccessStyle.Render(IconSuccess + " " + msg))
}

func (l *Logger) Error(msg string) {
	fmt.Println(ErrorStyle.Render(IconError + " " + msg))
}

func (l *Logger) Info(msg string) {
	fmt.Println(InfoStyle.Render(IconInfo + " " + msg))
}

func (l *Logger) Warning(msg string) {
	fmt.Println(ErrorStyle.Render(IconWarning + " " + msg))
}

func (l *Logger) Step(step int, total int, msg string) {
	stepStr := fmt.Sprintf("[%d/%d]", step, total)
	fmt.Println(InfoStyle.Render(stepStr) + " " + msg)
}

func (l *Logger) Box(title string, items []string) {
	content := title + "\n\n"
	for _, item := range items {
		content += "  • " + item + "\n"
	}
	fmt.Println(BoxStyle.Render(content))
}

func (l *Logger) Separator() {
	fmt.Println(strings.Repeat("─", 60))
}

func (l *Logger) NewLine() {
	fmt.Println()
}

// PrintBanner prints a beautiful banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║  ███████╗██╗     ██╗███╗   ██╗███████╗    ██████╗██╗     ██╗  ║
║  ██╔════╝██║     ██║████╗  ██║██╔════╝   ██╔════╝██║     ██║  ║
║  █████╗  ██║     ██║██╔██╗ ██║█████╗     ██║     ██║     ██║  ║
║  ██╔══╝  ██║     ██║██║╚██╗██║██╔══╝     ██║     ██║     ██║  ║
║  ██║     ███████╗██║██║ ╚████║███████╗   ╚██████╗███████╗██║  ║
║  ╚═╝     ╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝    ╚═════╝╚══════╝╚═╝  ║
║                                                               ║
║        Flutter Project Generator with Pine Architecture       ║
║                   written by yooocoding                       ║
╚═══════════════════════════════════════════════════════════════╝
`
	fmt.Println(TitleStyle.Render(banner))
}
