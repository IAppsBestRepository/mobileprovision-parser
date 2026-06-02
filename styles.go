package main

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7c3aed"))
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6366f1")).Bold(true).Width(20)
	valueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e4e4e7"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#525252"))
	sectionStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b5cf6")).Bold(true)
	deviceStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#a1a1aa"))
	borderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3f3f46"))
	entEnabled    = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	warningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	entConfigured = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))
	entDisabled   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
)
