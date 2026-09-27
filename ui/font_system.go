package ui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/image/font/sfnt"
)

// Chinese text is drawn with a font of the system: CJK fonts are too large
// to embed. The font is looked up the first time text has CJK characters,
// so applications without such text never touch the system fonts.

// SystemFallbackFonts enables the system CJK font. Set it to false before
// any CJK text is shown to use only the fonts registered by the application.
var SystemFallbackFonts = true

// FontFamilyCJK is the family the system CJK font is registered as, when found.
const FontFamilyCJK = "system-cjk"

// systemFontCandidate is a font file to try; families are the preferred
// fonts of a collection (.ttc), in order.
type systemFontCandidate struct {
	pattern  string // a path, may have filepath.Glob wildcards
	families []string
}

var systemFallbackOnce sync.Once

// ensureFallbackFonts loads the system CJK font the first time text has
// CJK characters. Call it before taking a face: loading resets the faces.
func ensureFallbackFonts(text string) {
	if hasCJK(text) {
		systemFallbackOnce.Do(loadSystemCJKFont)
	}
}

// hasCJK reports whether text has characters from the CJK blocks: radicals,
// punctuation, kana, ideographs, Hangul, full-width forms.
func hasCJK(text string) bool {
	for _, r := range text {
		if r >= 0x2E80 {
			return true
		}
	}
	return false
}

func loadSystemCJKFont() {
	if !SystemFallbackFonts {
		return
	}
	for _, candidate := range systemCJKFontCandidates() {
		paths, _ := filepath.Glob(candidate.pattern)
		for _, path := range paths {
			if f, err := openSystemFont(path, candidate.families); err == nil {
				setFont(FontFamilyCJK, f)
				AddFallbackFont(FontFamilyCJK)
				return
			}
		}
	}
	// Any other Unix: ask fontconfig
	if path := fontconfigMatch(":lang=zh-cn"); path != "" {
		if f, err := openSystemFont(path, nil); err == nil {
			setFont(FontFamilyCJK, f)
			AddFallbackFont(FontFamilyCJK)
		}
	}
}

func systemCJKFontCandidates() []systemFontCandidate {
	switch runtime.GOOS {
	case "windows":
		dir := filepath.Join(os.Getenv("WINDIR"), "Fonts")
		if os.Getenv("WINDIR") == "" {
			dir = `C:\Windows\Fonts`
		}
		return []systemFontCandidate{
			{filepath.Join(dir, "msyh.ttc"), []string{"Microsoft YaHei UI", "Microsoft YaHei"}},
			{filepath.Join(dir, "msyh.ttf"), []string{"Microsoft YaHei"}},
			{filepath.Join(dir, "Deng.ttf"), []string{"DengXian"}},
			{filepath.Join(dir, "simhei.ttf"), []string{"SimHei"}},
			{filepath.Join(dir, "simsun.ttc"), []string{"SimSun"}},
		}
	case "darwin":
		return []systemFontCandidate{
			{"/System/Library/Fonts/PingFang.ttc", []string{"PingFang SC"}},
			// Since macOS 13 PingFang is a downloadable asset
			{"/System/Library/AssetsV2/com_apple_MobileAsset_Font*/*/AssetData/PingFang.ttc", []string{"PingFang SC"}},
			{"/System/Library/Fonts/Hiragino Sans GB.ttc", []string{"Hiragino Sans GB"}},
			{"/System/Library/Fonts/STHeiti Light.ttc", []string{"Heiti SC", "STHeiti"}},
			{"/System/Library/Fonts/STHeiti Medium.ttc", []string{"Heiti SC", "STHeiti"}},
			{"/System/Library/Fonts/Supplemental/Arial Unicode.ttf", nil},
			{"/Library/Fonts/Arial Unicode.ttf", nil},
		}
	}
	notoCJK := []string{"Noto Sans CJK SC"}
	sourceHan := []string{"Source Han Sans SC", "Source Han Sans CN"}
	return []systemFontCandidate{
		{"/usr/share/fonts/*/NotoSansCJK-Regular.ttc", notoCJK},
		{"/usr/share/fonts/*/*/NotoSansCJK-Regular.ttc", notoCJK},
		{"/usr/share/fonts/*/NotoSansCJKsc-Regular.otf", notoCJK},
		{"/usr/share/fonts/*/*/NotoSansCJKsc-Regular.otf", notoCJK},
		{"/usr/share/fonts/*/SourceHanSans*-Regular.tt?", sourceHan},
		{"/usr/share/fonts/*/*/SourceHanSans*-Regular.?t?", sourceHan},
		{"/usr/share/fonts/*/*/wqy-microhei.ttc", []string{"WenQuanYi Micro Hei"}},
		{"/usr/share/fonts/*/wqy-microhei.ttc", []string{"WenQuanYi Micro Hei"}},
		{"/usr/share/fonts/*/*/wqy-zenhei.ttc", []string{"WenQuanYi Zen Hei"}},
		{"/usr/share/fonts/*/*/DroidSansFallbackFull.ttf", nil},
	}
}

// fontconfigMatch returns the file of the font fontconfig picks for the
// pattern, "" if fontconfig isn't there.
func fontconfigMatch(pattern string) string {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return ""
	}
	out, err := exec.Command("fc-match", "-f", "%{file}", pattern).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// openSystemFont opens the font file, reading it on demand rather than
// loading it: CJK fonts are tens of megabytes. From a collection it picks
// the first of the preferred families, the Regular one if there are several.
// The font must have Chinese characters.
func openSystemFont(path string, families []string) (*sfnt.Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	collection, err := sfnt.ParseCollectionReaderAt(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	var buf sfnt.Buffer
	fonts := make([]*sfnt.Font, 0, collection.NumFonts())
	for i := 0; i < collection.NumFonts(); i++ {
		if f, err := collection.Font(i); err == nil && hasGlyph(f, &buf, '中') {
			fonts = append(fonts, f)
		}
	}

	// The file stays open while the font is in use: the application's lifetime
	for _, family := range families {
		var match *sfnt.Font
		for _, f := range fonts {
			if !strings.EqualFold(fontName(f, &buf, sfnt.NameIDTypographicFamily, sfnt.NameIDFamily), family) {
				continue
			}
			if match == nil || strings.EqualFold(fontName(f, &buf, sfnt.NameIDTypographicSubfamily, sfnt.NameIDSubfamily), "Regular") {
				match = f
			}
		}
		if match != nil {
			return match, nil
		}
	}
	if len(fonts) > 0 {
		return fonts[0], nil
	}
	file.Close()
	return nil, errors.New("no Chinese characters in " + path)
}

func hasGlyph(f *sfnt.Font, buf *sfnt.Buffer, r rune) bool {
	index, err := f.GlyphIndex(buf, r)
	return err == nil && index != 0
}

// fontName returns the first of the names the font has.
func fontName(f *sfnt.Font, buf *sfnt.Buffer, ids ...sfnt.NameID) string {
	for _, id := range ids {
		if name, err := f.Name(buf, id); err == nil && name != "" {
			return name
		}
	}
	return ""
}
