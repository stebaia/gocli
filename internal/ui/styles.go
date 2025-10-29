package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	// Colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Error     = lipgloss.Color("#EF4444") // Red
	Warning   = lipgloss.Color("#F59E0B") // Orange
	Info      = lipgloss.Color("#3B82F6") // Blue
	Muted     = lipgloss.Color("#6B7280") // Gray

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Legacy color functions for simple usage
	ColorPrimary   = color.New(color.FgMagenta, color.Bold)
	ColorSuccess   = color.New(color.FgGreen, color.Bold)
	ColorError     = color.New(color.FgRed, color.Bold)
	ColorInfo      = color.New(color.FgCyan)
	ColorWarning   = color.New(color.FgYellow, color.Bold)
	ColorMuted     = color.New(color.FgHiBlack)
)

// Icons
const (
	IconSuccess = "✓"
	IconError   = "✗"
	IconInfo    = "ℹ"
	IconWarning = "⚠"
	IconRocket  = "🚀"
	IconPackage = "📦"
	IconFile    = "📄"
	IconFolder  = "📁"
	IconCode    = "💻"
	IconFire    = "🔥"
	IconCheck   = "✅"
	IconCross   = "❌"
	IconPin     = "📌"
)
