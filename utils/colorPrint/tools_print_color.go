package colorPrint

// Séquences d'échappement ANSI pour les couleurs
const (
	Reset   = "\033[0m"
	Blue    = "\033[34m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Grey    = "\033[90m"
    Black   = "\033[30m"
    Purple  = "\033[35m"
)

// Sequence d'échapement ANSI pour les couleur en Bold
const (
    BlueB    = "\033[1;34m"
    GreenB   = "\033[1;32m"
    YellowB  = "\033[1;33m"
    RedB     = "\033[1;31m"
    GreyB    = "\033[1;90m"
    BlackB   = "\033[1;30m"
    PurpleB  = "\033[1;35m"
)

//Sequence d'échapement ANSI pour les couleur en Italic
const (
    BlueI    = "\033[3;34m"
    GreenI   = "\033[3;32m"
    YellowI  = "\033[3;33m"
    RedI     = "\033[3;31m"
    GreyI    = "\033[3;90m"
    BlackI   = "\033[3;30m"
    PurpleI  = "\033[3;35m"
)

//Sequence d'échapement ANSI pour les couleur en Underline
const (
    BlueU    = "\033[4;34m"
    GreenU   = "\033[4;32m"
    YellowU  = "\033[4;33m"
    RedU     = "\033[4;31m"
    GreyU    = "\033[4;90m"
    BlackU   = "\033[4;30m"
    PurpleU  = "\033[4;35m"
)

//Sequence d'échapement ANSI pour le Background des couleurs
const (
    BlueBg    = "\033[44m"
    GreenBg   = "\033[42m"
    YellowBg  = "\033[43m"
    RedBg     = "\033[41m"
    GreyBg    = "\033[100m"
    BlackBg   = "\033[40m"
    PurpleBg  = "\033[45m"
    WhiteBg   = "\033[47m"
)

//Sequence d'échapement ANSI pour le Background des couleurs en Bold
const (
    BlueBgB    = "\033[1;44m"
    GreenBgB   = "\033[1;42m"
    YellowBgB  = "\033[1;43m"
    RedBgB     = "\033[1;41m"
    GreyBgB    = "\033[1;100m"
    BlackBgB   = "\033[1;40m"
    PurpleBgB  = "\033[1;45m"
)

func BlueP(content string) string {
	return Blue + content + Reset
}

func GreenP(content string) string {
    return Green + content + Reset
}

func YellowP(content string) string {
    return Yellow + content + Reset
}

func RedP(content string) string {
    return Red + content + Reset
}

func GreyP(content string) string {
    return Grey + content + Reset
}

func BlackP(content string) string {
    return Black + content + Reset
}

func PurpleP(content string) string {
    return Purple + content + Reset
}