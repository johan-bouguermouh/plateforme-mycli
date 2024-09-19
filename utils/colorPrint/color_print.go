package colorPrint

type Options struct {
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	Background    string
}

func ColorPrint(color string, content string, option *Options) string {
	var colorCode string = determineColor(color)

	if option != nil {
		if option.Background != "" {
			background := determineBackground(option.Background)
			colorCode = background + colorCode
		}
		if option.Bold {
			colorCode = "\033[1m" + colorCode
		}
		if option.Italic {
			colorCode = "\033[3m" + colorCode
		}
		if option.Strikethrough {
			colorCode = "\033[9m" + colorCode
		}
		if option.Underline {
			colorCode = "\033[4m" + colorCode
		}
	}
	return colorCode + content + Reset
}

func determineColor(value string) string {
	switch value {
	case "Blue":
		return Blue
	case "Green":
		return Green
	case "Yellow":
		return Yellow
	case "Red":
		return Red
	case "Grey":
		return Grey
	case "Black":
		return Black
	case "Purple":
		return Purple
	default:
		println(YellowP("WARNIG | Color not found, argument ignored : "), value)
		return ""
	}
}

func determineBackground(value string) string {
	switch value {
	case "Blue":
		return BlueBg
	case "Green":
		return GreenBg
	case "Yellow":
		return YellowBg
	case "Red":
		return RedBg
	case "Grey":
		return GreyBg
	case "Black":
		return BlackBg
	case "Purple":
		return PurpleBg
	case "White":
		return WhiteBg
	default:
		println(YellowP("WARNIG | Background color not found, argument ignored : "), value)
		return ""
	}
}