package ui

import "image/color"

var Theme map[string]interface{}

var DefaultBackground = color.RGBA{0, 0, 0, 255}
var DefaultForeground = color.RGBA{255, 255, 255, 255}

func init() {
	Theme = make(map[string]interface{})
	Theme["background"] = color.RGBA{10, 10, 10, 255}
	Theme["background.accent1"] = color.RGBA{30, 30, 30, 255}
	Theme["background.accent2"] = color.RGBA{50, 50, 50, 255}
	Theme["foreground"] = color.RGBA{255, 255, 255, 255}
	Theme["fontFamily"] = "robotomono"
	Theme["fontSize"] = 16
}

func ThemeBackgroundColor() color.RGBA {
	if v, ok := Theme["background"]; ok {
		if bgColor, ok := v.(color.RGBA); ok {
			return bgColor
		}
	}
	return DefaultBackground
}

func ThemeBackgroundColorAccent1() color.RGBA {
	if v, ok := Theme["background.accent1"]; ok {
		if bgHoverColor, ok := v.(color.RGBA); ok {
			return bgHoverColor
		}
	}
	return color.RGBA{50, 50, 50, 255} // Default accent1 background color
}

func ThemeBackgroundColorAccent2() color.RGBA {
	if v, ok := Theme["background.accent2"]; ok {
		if bgAccent2Color, ok := v.(color.RGBA); ok {
			return bgAccent2Color
		}
	}
	return color.RGBA{50, 50, 50, 255} // Default accent2 background color
}

func ThemeForegroundColor() color.RGBA {
	if v, ok := Theme["foreground"]; ok {
		if fgColor, ok := v.(color.RGBA); ok {
			return fgColor
		}
	}
	return DefaultForeground
}

func ThemeFontFamily() string {
	if v, ok := Theme["fontFamily"]; ok {
		if fontFamily, ok := v.(string); ok {
			return fontFamily
		}
	}
	return "robotomono" // Default font family
}

func ThemeFontSize() float64 {
	if v, ok := Theme["fontSize"]; ok {
		if fontSize, ok := v.(float64); ok {
			return fontSize
		}
	}
	return 16.0 // Default font size
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

func GetThemeInt(name string, defaultValue int) int {
	if v, ok := Theme[name]; ok {
		if intValue, ok := v.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

func GetThemeFloat(name string, defaultValue float64) float64 {
	if v, ok := Theme[name]; ok {
		if floatValue, ok := v.(float64); ok {
			return floatValue
		}
	}
	return defaultValue
}
