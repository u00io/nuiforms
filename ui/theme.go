package ui

import "image/color"

var Theme map[string]interface{}

var DefaultBackground = color.RGBA{0, 0, 0, 255}
var DefaultForeground = color.RGBA{255, 255, 255, 255}

func init() {
	Theme = make(map[string]interface{})
	Theme["background"] = color.RGBA{0, 0, 0, 255}
	Theme["foreground"] = color.RGBA{255, 255, 255, 255}
	Theme["fontFamily"] = "robotomono"
	Theme["fontSize"] = 16
	Theme["borderColor"] = color.RGBA{255, 255, 255, 255}

	Theme["button.background"] = color.RGBA{50, 50, 50, 255}
	Theme["button.background.hover"] = color.RGBA{70, 70, 70, 255}
	Theme["button.background.pressed"] = color.RGBA{90, 90, 90, 255}
}

func GetThemeColor(name string, defaultColor color.RGBA) color.RGBA {
	if v, ok := Theme[name]; ok {
		if colorValue, ok := v.(color.RGBA); ok {
			return colorValue
		}
	}
	return defaultColor
}

func GetThemeString(name string, defaultValue string) string {
	if v, ok := Theme[name]; ok {
		if strValue, ok := v.(string); ok {
			return strValue
		}
	}
	return defaultValue
}
