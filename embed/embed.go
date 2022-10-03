package embed

import (
	"golang.org/x/image/colornames"
	"image/color"
	"main/info"
	"time"
)

func colorToInt(color color.RGBA) int {
	return 256*256*int(color.R) + 256*int(color.G) + int(color.B)
}

type Type uint8

const (
	TypeInfo Type = iota
	TypeSuccess
	TypeWarning
	TypeError
)

func colorFromType(t Type) color.RGBA {
	var c color.RGBA
	switch t {
	case TypeInfo:
		c = colornames.Cyan
	case TypeSuccess:
		c = colornames.Greenyellow
	case TypeError:
		c = colornames.Red
	case TypeWarning:
		c = colornames.Orangered
	}
	return c
}

func NewSimpleEmbed(title string, description string, t Type) *Builder {
	return NewBuilder().Title(title).Description(description).Footer().Text(info.BotName + " " + info.BotVersion).Timestamp(time.Now()).Color(colorFromType(t))
}
