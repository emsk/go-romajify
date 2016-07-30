package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

// Commands is list of subcommands for Romajify
var Commands = []cli.Command{
	hepburnCommand,
	nihonCommand,
	kunreiCommand,
}

var hepburnCommand = cli.Command{
	Name:   "hepburn",
	Usage:  "Convert kana to Hepburn romaji",
	Action: hepburnAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "upcase",
			Usage: "Convert to uppercase",
		},
		cli.BoolFlag{
			Name:  "traditional",
			Usage: "Convert to traditional Hepburn romaji",
		},
	},
}

var nihonCommand = cli.Command{
	Name:   "nihon",
	Usage:  "Convert kana to Nihon-shiki romaji",
	Action: nihonAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "upcase",
			Usage: "Convert to uppercase",
		},
	},
}

var kunreiCommand = cli.Command{
	Name:   "kunrei",
	Usage:  "Convert kana to Kunrei-shiki romaji",
	Action: kunreiAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "upcase",
			Usage: "Convert to uppercase",
		},
	},
}

var monographs = map[string]string{
	"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
	"ア": "a", "イ": "i", "ウ": "u", "エ": "e", "オ": "o",

	"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko",
	"カ": "ka", "キ": "ki", "ク": "ku", "ケ": "ke", "コ": "ko",

	"が": "ga", "ぎ": "gi", "ぐ": "gu", "げ": "ge", "ご": "go",
	"ガ": "ga", "ギ": "gi", "グ": "gu", "ゲ": "ge", "ゴ": "go",

	"さ": "sa", "し": "shi", "す": "su", "せ": "se", "そ": "so",
	"サ": "sa", "シ": "shi", "ス": "su", "セ": "se", "ソ": "so",

	"ざ": "za", "じ": "ji", "ず": "zu", "ぜ": "ze", "ぞ": "zo",
	"ザ": "za", "ジ": "ji", "ズ": "zu", "ゼ": "ze", "ゾ": "zo",

	"た": "ta", "ち": "chi", "つ": "tsu", "て": "te", "と": "to",
	"タ": "ta", "チ": "chi", "ツ": "tsu", "テ": "te", "ト": "to",

	"だ": "da", "ぢ": "ji", "づ": "zu", "で": "de", "ど": "do",
	"ダ": "da", "ヂ": "ji", "ヅ": "zu", "デ": "de", "ド": "do",

	"な": "na", "に": "ni", "ぬ": "nu", "ね": "ne", "の": "no",
	"ナ": "na", "ニ": "ni", "ヌ": "nu", "ネ": "ne", "ノ": "no",

	"は": "ha", "ひ": "hi", "ふ": "fu", "へ": "he", "ほ": "ho",
	"ハ": "ha", "ヒ": "hi", "フ": "fu", "ヘ": "he", "ホ": "ho",

	"ば": "ba", "び": "bi", "ぶ": "bu", "べ": "be", "ぼ": "bo",
	"バ": "ba", "ビ": "bi", "ブ": "bu", "ベ": "be", "ボ": "bo",

	"ぱ": "pa", "ぴ": "pi", "ぷ": "pu", "ぺ": "pe", "ぽ": "po",
	"パ": "pa", "ピ": "pi", "プ": "pu", "ペ": "pe", "ポ": "po",

	"ま": "ma", "み": "mi", "む": "mu", "め": "me", "も": "mo",
	"マ": "ma", "ミ": "mi", "ム": "mu", "メ": "me", "モ": "mo",

	"や": "ya", "ゆ": "yu", "よ": "yo",
	"ヤ": "ya", "ユ": "yu", "ヨ": "yo",

	"ら": "ra", "り": "ri", "る": "ru", "れ": "re", "ろ": "ro",
	"ラ": "ra", "リ": "ri", "ル": "ru", "レ": "re", "ロ": "ro",

	"わ": "wa", "ゐ": "i", "ゑ": "e", "を": "o", "ん": "n",
	"ワ": "wa", "ヰ": "i", "ヱ": "e", "ヲ": "o", "ン": "n",

	"ぁ": "a", "ぃ": "i", "ぅ": "u", "ぇ": "e", "ぉ": "o",
	"ァ": "a", "ィ": "i", "ゥ": "u", "ェ": "e", "ォ": "o",

	"ゃ": "ya", "ゅ": "yu", "ょ": "yo",
	"ャ": "ya", "ュ": "yu", "ョ": "yo",

	"ゔ": "bu", "ヴ": "bu", "ー": "", "＿": "_",
}

var digraphs = map[string]string{
	"きゃ": "kya", "きゅ": "kyu", "きょ": "kyo",
	"キャ": "kya", "キュ": "kyu", "キョ": "kyo",

	"ぎゃ": "gya", "ぎゅ": "gyu", "ぎょ": "gyo",
	"ギャ": "gya", "ギュ": "gyu", "ギョ": "gyo",

	"しゃ": "sha", "しゅ": "shu", "しょ": "sho",
	"シャ": "sha", "シュ": "shu", "ショ": "sho",

	"じゃ": "ja", "じゅ": "ju", "じょ": "jo",
	"ジャ": "ja", "ジュ": "ju", "ジョ": "jo",

	"ちゃ": "cha", "ちゅ": "chu", "ちょ": "cho",
	"チャ": "cha", "チュ": "chu", "チョ": "cho",

	"ぢゃ": "ja", "ぢゅ": "ju", "ぢょ": "jo",
	"ヂャ": "ja", "ヂュ": "ju", "ヂョ": "jo",

	"にゃ": "nya", "にゅ": "nyu", "にょ": "nyo",
	"ニャ": "nya", "ニュ": "nyu", "ニョ": "nyo",

	"ひゃ": "hya", "ひゅ": "hyu", "ひょ": "hyo",
	"ヒャ": "hya", "ヒュ": "hyu", "ヒョ": "hyo",

	"びゃ": "bya", "びゅ": "byu", "びょ": "byo",
	"ビャ": "bya", "ビュ": "byu", "ビョ": "byo",

	"ぴゃ": "pya", "ぴゅ": "pyu", "ぴょ": "pyo",
	"ピャ": "pya", "ピュ": "pyu", "ピョ": "pyo",

	"みゃ": "mya", "みゅ": "myu", "みょ": "myo",
	"ミャ": "mya", "ミュ": "myu", "ミョ": "myo",

	"りゃ": "rya", "りゅ": "ryu", "りょ": "ryo",
	"リャ": "rya", "リュ": "ryu", "リョ": "ryo",
}

var nihonMonographs = map[string]string{
	"し": "si", "ち": "ti", "つ": "tu", "ふ": "hu", "じ": "zi", "ぢ": "di", "づ": "du",
	"シ": "si", "チ": "ti", "ツ": "tu", "フ": "hu", "ジ": "zi", "ヂ": "di", "ヅ": "du",

	"ゐ": "wi", "ゑ": "we", "を": "wo",
	"ヰ": "wi", "ヱ": "we", "ヲ": "wo",
}

var nihonDigraphs = map[string]string{
	"しゃ": "sya", "しゅ": "syu", "しょ": "syo",
	"シャ": "sya", "シュ": "syu", "ショ": "syo",

	"じゃ": "zya", "じゅ": "zyu", "じょ": "zyo",
	"ジャ": "zya", "ジュ": "zyu", "ジョ": "zyo",

	"ちゃ": "tya", "ちゅ": "tyu", "ちょ": "tyo",
	"チャ": "tya", "チュ": "tyu", "チョ": "tyo",

	"ぢゃ": "dya", "ぢゅ": "dyu", "ぢょ": "dyo",
	"ヂャ": "dya", "ヂュ": "dyu", "ヂョ": "dyo",
}

var kunreiMonographs = map[string]string{
	"し": "si", "ち": "ti", "つ": "tu", "ふ": "hu", "じ": "zi", "ぢ": "zi",
	"シ": "si", "チ": "ti", "ツ": "tu", "フ": "hu", "ジ": "zi", "ヂ": "zi",
}

var kunreiDigraphs = map[string]string{
	"しゃ": "sya", "しゅ": "syu", "しょ": "syo",
	"シャ": "sya", "シュ": "syu", "ショ": "syo",

	"じゃ": "zya", "じゅ": "zyu", "じょ": "zyo",
	"ジャ": "zya", "ジュ": "zyu", "ジョ": "zyo",

	"ちゃ": "tya", "ちゅ": "tyu", "ちょ": "tyo",
	"チャ": "tya", "チュ": "tyu", "チョ": "tyo",

	"ぢゃ": "zya", "ぢゅ": "zyu", "ぢょ": "zyo",
	"ヂャ": "zya", "ヂュ": "zyu", "ヂョ": "zyo",
}

func hepburnAction(c *cli.Context) error {
	resultText := c.Args().First()

	resultText = romanize(resultText, digraphs)
	resultText = romanize(resultText, monographs)

	// Double consonants: 促音
	r := regexp.MustCompile(`[っッ]c`)
	resultText = r.ReplaceAllString(resultText, "tc")
	r = regexp.MustCompile(`[っッ](.)`)
	resultText = r.ReplaceAllString(resultText, "$1$1")

	// Syllabic n: 撥音
	if c.Bool("traditional") {
		r = regexp.MustCompile(`n([bmp])`)
		resultText = r.ReplaceAllString(resultText, "m$1")
	}

	// Long vowels: 長音
	r = regexp.MustCompile(`oo(.+)`)
	resultText = r.ReplaceAllString(resultText, "o$1")
	r = regexp.MustCompile(`ou`)
	resultText = r.ReplaceAllString(resultText, "o")
	r = regexp.MustCompile(`uu`)
	resultText = r.ReplaceAllString(resultText, "u")

	if c.Bool("upcase") {
		resultText = strings.ToUpper(resultText)
	}

	fmt.Println(resultText)
	return nil
}

func nihonAction(c *cli.Context) error {
	resultText := c.Args().First()

	resultText = romanize(resultText, mergeMap(digraphs, nihonDigraphs))
	resultText = romanize(resultText, mergeMap(monographs, nihonMonographs))

	// Double consonants: 促音
	r := regexp.MustCompile(`[っッ](.)`)
	resultText = r.ReplaceAllString(resultText, "$1$1")

	// Long vowels: 長音
	r = regexp.MustCompile(`ou|oo`)
	resultText = r.ReplaceAllString(resultText, "o")
	r = regexp.MustCompile(`uu`)
	resultText = r.ReplaceAllString(resultText, "u")

	if c.Bool("upcase") {
		resultText = strings.ToUpper(resultText)
	}

	fmt.Println(resultText)
	return nil
}

func kunreiAction(c *cli.Context) error {
	resultText := c.Args().First()

	resultText = romanize(resultText, mergeMap(digraphs, kunreiDigraphs))
	resultText = romanize(resultText, mergeMap(monographs, kunreiMonographs))

	// Double consonants: 促音
	r := regexp.MustCompile(`[っッ](.)`)
	resultText = r.ReplaceAllString(resultText, "$1$1")

	// Long vowels: 長音
	r = regexp.MustCompile(`ou|oo`)
	resultText = r.ReplaceAllString(resultText, "o")
	r = regexp.MustCompile(`uu`)
	resultText = r.ReplaceAllString(resultText, "u")

	if c.Bool("upcase") {
		resultText = strings.ToUpper(resultText)
	}

	fmt.Println(resultText)
	return nil
}

func romanize(text string, chars map[string]string) string {
	resultText := text

	for kana, romaji := range chars {
		resultText = strings.Replace(resultText, kana, romaji, -1)
	}

	return resultText
}

func mergeMap(map1, map2 map[string]string) map[string]string {
	resultMap := map1

	for k, v := range map2 {
		resultMap[k] = v
	}

	return resultMap
}
