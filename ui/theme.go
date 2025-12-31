package ui

import (
	"fmt"
	"image/color"
)

var Theme map[string]interface{}

/*var DefaultSurfaceBackground = ColorFromHex("#121212")
var DefaultPrimaryBackground = ColorFromHex("#0D47A1")
var DefaultSecondaryBackground = ColorFromHex("#26C6DA")

var DefaultSurface = ColorFromHex("#FFFFFF")
var DefaultPrimary = ColorFromHex("#FFFFFF")
var DefaultSecondary = ColorFromHex("#FFFFFF")*/

var IsDarkTheme = true

func ApplyDarkTheme() {
	IsDarkTheme = true
	Theme["background.surface"] = ColorFromHex("#121212")
	Theme["background.primary"] = ColorFromHex("#0D47A1")
	Theme["background.secondary"] = ColorFromHex("#26C6DA")
	Theme["foreground.surface"] = ColorFromHex("#FFFFFF")
	Theme["foreground.primary"] = ColorFromHex("#FFFFFF")
	Theme["foreground.secondary"] = ColorFromHex("#FFFFFF")
	UpdateMainForm()
}

func ApplyLightTheme() {
	IsDarkTheme = false
	Theme["background.surface"] = ColorFromHex("#FFFFFF")
	Theme["background.primary"] = ColorFromHex("#4d9dec")
	Theme["background.secondary"] = ColorFromHex("#00BCD4")
	Theme["foreground.surface"] = ColorFromHex("#000000")
	Theme["foreground.primary"] = ColorFromHex("#FFFFFF")
	Theme["foreground.secondary"] = ColorFromHex("#FFFFFF")
	UpdateMainForm()
}

func ApplyBaseFontSize(fontSize float64) {
	Theme["fontSize"] = fontSize
	UpdateMainFormLayout()
	UpdateMainForm()
}

func ColorFromHex(hexStr string) color.RGBA {
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

func ColorToHex(col color.Color) string {
	r, g, b, a := col.RGBA()
	return fmt.Sprintf("#%02X%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

func init() {
	Theme = make(map[string]interface{})

	/*Theme["background.surface"] = DefaultSurfaceBackground
	Theme["background.primary"] = DefaultPrimaryBackground
	Theme["background.secondary"] = DefaultSecondaryBackground

	Theme["foreground.surface"] = DefaultSurface
	Theme["foreground.primary"] = DefaultPrimary
	Theme["foreground.secondary"] = DefaultSecondary*/

	ApplyDarkTheme()

	Theme["fontFamily"] = "robotomono"
	Theme["fontSize"] = 18
}

func ThemeBackgroundColorDarkTheme(elevation int, role string) color.RGBA {
	if role == "" {
		role = "surface"
	}
	if v, ok := Theme["background."+role]; ok {
		if bgColor, ok := v.(color.RGBA); ok {
			gray := 0x12
			gray = gray + 16 + elevation*6
			if gray < 0 {
				gray = 0
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
	return ColorFromHex("#000000")
}

func ThemeBackgroundColorLightTheme(elevation int, role string) color.RGBA {
	if role == "" {
		role = "surface"
	}
	if v, ok := Theme["background."+role]; ok {
		if bgColor, ok := v.(color.RGBA); ok {
			gray := 0x12

			gray = gray + 16 + elevation*6
			if gray < 0 {
				gray = 0
			}
			if gray > 255 {
				gray = 255
			}

			r := int(bgColor.R) - gray
			g := int(bgColor.G) - gray
			b := int(bgColor.B) - gray

			if r < 0 {
				r = 0
			}

			if g < 0 {
				g = 0
			}

			if b < 0 {
				b = 0
			}

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
	return ColorFromHex("#FFFFFF")
}

func ThemeBackgroundColor(elevation int, role string) color.RGBA {
	if IsDarkTheme {
		return ThemeBackgroundColorDarkTheme(elevation, role)
	} else {
		return ThemeBackgroundColorLightTheme(elevation, role)
	}
}

func ThemeForegroundColor(role string) color.RGBA {
	if role == "" {
		role = "surface"
	}
	if v, ok := Theme["foreground."+role]; ok {
		if fgColor, ok := v.(color.RGBA); ok {
			return fgColor
		}
	}
	return ColorFromHex("#777777")
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
