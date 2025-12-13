package ui

import (
	"fmt"
	"image/color"
)

var Theme map[string]interface{}

var DefaultBackground = colorFromHex("#121212")

var DefaultPrimary = colorFromHex("#42A5F5")
var DefaultSecondary = colorFromHex("#26C6DA")

var DefaultOnSurface = colorFromHex("#FFFFFF")
var DefaultOnPrimary = colorFromHex("#000000")
var DefaultOnSecondary = colorFromHex("#000000")

func colorFromHex(hexStr string) color.RGBA {
	var r, g, b, a uint8
	a = 255
	if len(hexStr) == 7 {
		_, err := fmt.Sscanf(hexStr, "#%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return color.RGBA{0, 0, 0, 255}
		}
	} else if len(hexStr) == 9 {
		_, err := fmt.Sscanf(hexStr, "#%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil {
			return color.RGBA{0, 0, 0, 255}
		}
	} else {
		return color.RGBA{0, 0, 0, 255}
	}
	return color.RGBA{r, g, b, a}
}

func init() {
	Theme = make(map[string]interface{})
	Theme["background"] = colorFromHex("#121212")
	Theme["background.accent1"] = colorFromHex("#0d668aff")
	Theme["background.accent2"] = colorFromHex("#2D2D2D")
	Theme["background.selection"] = colorFromHex("#264F78")
	Theme["foreground"] = colorFromHex("#FFFFFF")
	Theme["fontFamily"] = "robotomono"
	Theme["fontSize"] = 18
}

func ThemeBackgroundColor(elevation int) color.RGBA {
	if v, ok := Theme["background"]; ok {
		if bgColor, ok := v.(color.RGBA); ok {
			gray := 0x12
			if elevation > 0 {
				gray = gray + 16 + elevation*6
			}
			if gray > 255 {
				gray = 255
			}
			r := int(bgColor.R) + gray
			g := int(bgColor.G) + gray
			b := int(bgColor.B) + gray

			if r > 255 {
				r = 255
			}

			if g > 255 {
				g = 255
			}

			if b > 255 {
				b = 255
			}

			return color.RGBA{uint8(r), uint8(g), uint8(b), bgColor.A}
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

func ThemeBackgroundColorSelection() color.RGBA {
	if v, ok := Theme["background.selection"]; ok {
		if bgSelectionColor, ok := v.(color.RGBA); ok {
			return bgSelectionColor
		}
	}
	return color.RGBA{38, 79, 120, 255} // Default selection background color
}

func ThemePrimaryColor() color.RGBA {
	if v, ok := Theme["primary"]; ok {
		if fgColor, ok := v.(color.RGBA); ok {
			return fgColor
		}
	}
	return DefaultPrimary
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
	return 18.0 // Default font size
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
