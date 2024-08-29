package netcat

import (
	"fmt"
	"math/rand"
)

// Définition des constantes pour les couleurs ANSI
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

// Fonction Color qui retourne le code ANSI pour une couleur donnée
func Color(name string) string {
	switch name {
	case "red":
		return Red
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "blue":
		return Blue
	case "magenta":
		return Magenta
	case "cyan":
		return Cyan
	default:
		return Reset
	}
}

func Colorize(name, colorName string) string {
	return fmt.Sprintf("%s%s%s", Color(colorName), name, Reset)
}

// Fonction pour obtenir une couleur aléatoire
func randomColor() string {
	colors := []string{"red", "green", "yellow", "blue", "magenta", "cyan"}
	return colors[rand.Intn(len(colors))]
}
