// videopoker - command line video poker game with color Unicode suits.
//
// version 1.0
package main

/* Some of the code in this file may look very C-like because it
   was translated to Go from the C language version of videopoker. */

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	)

const VERSION = "videopoker 1.0"

/* values for -uN unicode flag */

const (
	UNICODE_TTY = iota
	UNICODE_OFF   
	UNICODE_SUITS
	UNICODE_CARDS
	)

var unicode int = UNICODE_SUITS
// unicode is nonzero if unicode output is enabled.

var hands int = 0	// number of hands played

/* the card type, for holding infomation about the deck of cards */

type card struct
{
	index int	/* cards value, minus 1 */
	sym string	/* textual appearance */
	uc string	/* Unicode value for the card */
	suit int	/* card's suit (see just below) */
	gone int	/* true if it's been dealt */
}

/* The number of cards in the hand */

const CARDS = 5

/* The number of cards in the deck */

const CARDSINDECK = 52

const (
	CLUBS	= iota
	DIAMONDS
	HEARTS
	SPADES
	)

/* one-character suit designations */

const NUMSUITS = 4

var suitname [NUMSUITS]string = [NUMSUITS]string {
	"C",
	"D",
	"H",
	"S",
	}

/* Needed for recognizing royal flush, or tens or better (TEN),
   or jacks or better (JACK) */

const TEN = 9	/* the index of the "10" card */
const JACK = 10	/* the index of the "jack" card */

/* The standard deck of 52 cards */

var deck [CARDSINDECK]card = [CARDSINDECK]card {
/*	index, card, Unicode, suit, gone */
	{  1, " 2", "\U0001F0D2", CLUBS, 0 },
	{  2, " 3", "\U0001F0D3", CLUBS, 0 },
	{  3, " 4", "\U0001F0D4", CLUBS, 0 },
	{  4, " 5", "\U0001F0D5", CLUBS, 0 },
	{  5, " 6", "\U0001F0D6", CLUBS, 0 },
	{  6, " 7", "\U0001F0D7", CLUBS, 0 },
	{  7, " 8", "\U0001F0D8", CLUBS, 0 },
	{  8, " 9", "\U0001F0D9", CLUBS, 0 },
	{  9, "10", "\U0001F0Da", CLUBS, 0 },
	{ 10, " J", "\U0001F0Db", CLUBS, 0 },
	{ 11, " Q", "\U0001F0Dd", CLUBS, 0 },
	{ 12, " K", "\U0001F0De", CLUBS, 0 },
	{ 13, " A", "\U0001F0D1", CLUBS, 0 },

	{  1, " 2", "\U0001F0C2", DIAMONDS, 0 },
	{  2, " 3", "\U0001F0C3", DIAMONDS, 0 },
	{  3, " 4", "\U0001F0C4", DIAMONDS, 0 },
	{  4, " 5", "\U0001F0C5", DIAMONDS, 0 },
	{  5, " 6", "\U0001F0C6", DIAMONDS, 0 },
	{  6, " 7", "\U0001F0C7", DIAMONDS, 0 },
	{  7, " 8", "\U0001F0C8", DIAMONDS, 0 },
	{  8, " 9", "\U0001F0C9", DIAMONDS, 0 },
	{  9, "10", "\U0001F0Ca", DIAMONDS, 0 },
	{ 10, " J", "\U0001F0Cb", DIAMONDS, 0 },
	{ 11, " Q", "\U0001F0Cd", DIAMONDS, 0 },
	{ 12, " K", "\U0001F0Ce", DIAMONDS, 0 },
	{ 13, " A", "\U0001F0C1", DIAMONDS, 0 },

	{  1, " 2", "\U0001F0B2", HEARTS, 0 },
	{  2, " 3", "\U0001F0B3", HEARTS, 0 },
	{  3, " 4", "\U0001F0B4", HEARTS, 0 },
	{  4, " 5", "\U0001F0B5", HEARTS, 0 },
	{  5, " 6", "\U0001F0B6", HEARTS, 0 },
	{  6, " 7", "\U0001F0B7", HEARTS, 0 },
	{  7, " 8", "\U0001F0B8", HEARTS, 0 },
	{  8, " 9", "\U0001F0B9", HEARTS, 0 },
	{  9, "10", "\U0001F0Ba", HEARTS, 0 },
	{ 10, " J", "\U0001F0Bb", HEARTS, 0 },
	{ 11, " Q", "\U0001F0Bd", HEARTS, 0 },
	{ 12, " K", "\U0001F0Be", HEARTS, 0 },
	{ 13, " A", "\U0001F0B1", HEARTS, 0 },

	{  1, " 2", "\U0001F0A2", SPADES, 0 },
	{  2, " 3", "\U0001F0A3", SPADES, 0 },
	{  3, " 4", "\U0001F0A4", SPADES, 0 },
	{  4, " 5", "\U0001F0A5", SPADES, 0 },
	{  5, " 6", "\U0001F0A6", SPADES, 0 },
	{  6, " 7", "\U0001F0A7", SPADES, 0 },
	{  7, " 8", "\U0001F0A8", SPADES, 0 },
	{  8, " 9", "\U0001F0A9", SPADES, 0 },
	{  9, "10", "\U0001F0Aa", SPADES, 0 },
	{ 10, " J", "\U0001F0Ab", SPADES, 0 },
	{ 11, " Q", "\U0001F0Ad", SPADES, 0 },
	{ 12, " K", "\U0001F0Ae", SPADES, 0 },
	{ 13, " A", "\U0001F0A1", SPADES, 0 },
}

/* The hand. It holds five cards. */

var hand [5]card

/* sorted hand, for internal use when recognizing winners */

var shand [5]card

/* keys used to select kept cards */

const NUMKEYS = 5

var keys [NUMKEYS]byte = [NUMKEYS]byte {
	' ', 'j', 'k', 'l', ';',
}

/* initial number of chips held */

const INITCHIPS = 1000
var score int = INITCHIPS

/* minimum and maximum swing of score during this game */

var score_low int = INITCHIPS
var score_high int = INITCHIPS

/* The games starts with a bet of 10, the minimum allowed */

const INITMINBET = 10

var minbet int = INITMINBET
var bet int = INITMINBET

/* number of chips or groups of 10 chips bet */

var betmultiplier int = 1

/* Options */

/* -b (Bold): print in boldface */

var boldface bool = false

/* -mh (Mark Held): Mark cards that are held */

var markheld bool = false

/* -q (Quiet): Don't print banner or final report */

var quiet bool = false

/*
	Some ANSI Terminal escape codes:
	ESC[38;5; then one of (0m = black, 1m = red, 2m = green, 3m = yellow,
			       4m = blue, 5m = magenta, 6m = cyan, 7m = white)
	ESC[1m = bold, ESC[0m = reset all attributes
*/

/* The below are part of the escape codes listed above.
   Do not change the values. Do not use iota. */
const BLACK = 0
const RED = 1
const GREEN = 2
const YELLOW = 3
const BLUE = 4
const MAGENTA = 5
const CYAN = 6
const WHITE = 7

func color(color int) {
	if unicode == UNICODE_TTY { return }
	fmt.Printf("\033[38;5;%dm",color)
}

func ANSIbold() {
	if unicode == UNICODE_TTY { return }
	fmt.Printf("\033[1m")
}

func ANSIreset() {
	if unicode == UNICODE_TTY { return }
	fmt.Printf("\033[0m")
	if boldface { ANSIbold() }
}

/*
	Display the hand
	This is where the Unicode output setting (-u<N> option) takes effect,
	so there are three different ways it can display the cards.
*/

func showhand() {
//
	var i int
	/* Unicode characters for the suits */
	const spade string = "\u2660"
	const heart string = "\u2665"
	const diamond string = "\u2666"
	const club string = "\u2663"

	/* Method 1: Unicode Card Faces (-u3 option),
	   which requires only one line  */

	if unicode == UNICODE_CARDS {	/* print the Unicode card faces */
	//
		for i = 0; i < CARDS; i++ {
		//
			switch hand[i].suit {
				case DIAMONDS:
					fallthrough
				case HEARTS:
					color(RED)
//					fmt.Printf("\033[38;5;1m%s\033[0m", hand[i].uc)
					fmt.Printf("%s ", hand[i].uc)
					ANSIreset()
				case CLUBS:
					fallthrough
				case SPADES:
					fmt.Printf("%s ", hand[i].uc)
			}
		}
		/* print a space to separate output from user input */
		fmt.Printf(" ")
		return
	}

	/* Method 2: Two Line Output for -u0/-u1/-u2 options, requires 2 lines */

	if boldface { ANSIbold() }

	/* First Line of output: show card values */

	for i = 0; i < CARDS; i++ {
	//
		switch hand[i].suit {
			case DIAMONDS:
				fallthrough
			case HEARTS:
				/* print in red */
				color(RED)
//				fmt.Printf("\033[38;5;1m%s\033[0m ", hand[i].sym)
				fmt.Printf("%s", hand[i].sym)
				ANSIreset()
				fmt.Printf(" ")
			case CLUBS:
				fallthrough
			case SPADES:
				/* print in default text color */
				fmt.Printf("%s ", hand[i].sym)
		}
	}

	fmt.Printf("\n")

	for i = 0; i < CARDS; i++ {
	//
		if unicode == UNICODE_SUITS {
			/* Unicode method */
			switch hand[i].suit {
				case DIAMONDS:
					/* print in red */
					fmt.Printf(" ")
					color(RED)
//					fmt.Printf("\033[38;5;1m %s\033[0m ", diamond)
					fmt.Printf("%s", diamond)
					ANSIreset()
					fmt.Printf(" ")
				case HEARTS:
					fmt.Printf(" ")
					color(RED)
//					fmt.Printf("\033[38;5;1m %s\033[0m ", heart)
					fmt.Printf("%s", heart)
					ANSIreset()
					fmt.Printf(" ")
				case CLUBS:
					/* print in default text color */
					fmt.Printf(" %s ", club)
				case SPADES:
					fmt.Printf(" %s ", spade)
			}
		} else {
			/* ASCII method */
			switch hand[i].suit {
				case DIAMONDS:
					fallthrough
				case HEARTS:
					/* print H or D in red */
					fmt.Printf(" ")
					color(RED)
//					fmt.Printf(" \033[38;5;1m%s\033[0m ", suitname[hand[i].suit])
					fmt.Printf("%s", suitname[hand[i].suit])
					ANSIreset()
					fmt.Printf(" ")
				case CLUBS:
					fallthrough
				case SPADES:
					/* print S or C in default text color */
					fmt.Printf(" %s ", suitname[hand[i].suit])
			}
		}
	}

	/* print a space to separate output from user input */
	fmt.Printf(" ")
/*
	if check_for_dupes() == 0 { fmt.Printf("\n!!! DUPLICATE CARD !!!\n\n") }
*/
}

/* The various video poker games that are supported */

const (
	AllAmerican = iota
	TensOrBetter
	BonusPoker
	DoubleBonus
	DoubleBonusBonus
	JacksOrBetter	// default
	JacksOrBetter95
	JacksOrBetter86
	JacksOrBetter85
	JacksOrBetter75
	JacksOrBetter65
	NUMGAMES
)

/*
	The game in play. Default is Jacks or Better,
	which is coded into initialization of static data
*/

var game int = JacksOrBetter

var gamenames [NUMGAMES]string = [NUMGAMES]string {
	"All American",
	"Tens or Better",
	"Bonus Poker",
	"Double Bonus",
	"Double Bonus Bonus",
	"Jacks or Better",
	"9/5 Jacks or Better",
	"8/6 Jacks or Better",
	"8/5 Jacks or Better",
	"7/5 Jacks or Better",
	"6/5 Jacks or Better",
}

// TODO: the above consts, gamenames[], and option strings can
// be put in an array of structures:
// type gameinfo struct {} gameinfo;
// var game []gameinfo = { ... };  (etc.)
// Then change badgame(), option handling in main(), and setgame()

/* Error message for -g option. Also a way to display the list of games */

func badgame() {
//
	fmt.Printf("Video Poker: -g option is missing valid game name.\n")
	fmt.Printf("Available games are:\n")
	fmt.Printf("aa   - All American\n")
	fmt.Printf("10s  - Tens or Better\n")
/*
	fmt.Printf("bon  - Bonus Poker\n")
	fmt.Printf("db   - Double Bonus\n")
	fmt.Printf("dbb  - Double Bonus Bonus\n")
*/
	fmt.Printf("jb95 - 9/5 Jacks or Better\n")
	fmt.Printf("jb86 - 8/6 Jacks or Better\n")
	fmt.Printf("jb85 - 8/5 Jacks or Better\n")
	fmt.Printf("jb75 - 7/5 Jacks or Better\n")
	fmt.Printf("jb65 - 6/5 Jacks or Better\n")
	os.Exit(0)
}

/* replacements for C library getchar() and ungetc() functions */

var stdin *bufio.Reader
var char int

func getchar() int {
//
	var inputbyte byte

	inputbyte, _ = stdin.ReadByte()
	char = int(inputbyte)
	return char
}

func ungetc() {
//
	stdin.UnreadByte()
}

/* replacement for C library random() function */

var randomgen *rand.Rand

func srandom() {
//
	randomgen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func random() int {
//
	return randomgen.Int()
}

/* replacement for C library exit() function */

func exit(error int) {
//
	os.Exit(error)
}

/* Initialize, Handle arguments, enter loop */

func main() {
//
	var ai int
	var i, cnt int
	var arg string

	// open stdin
	stdin = bufio.NewReader(os.Stdin)

	/* initialize random number generator */

	srandom()

	/* process arguments */

	ai = 0

	for cnt = len(os.Args); cnt > 1; {
	//
		/* -b (Bold) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'b' &&
		   len(os.Args[1+ai]) == 2 {
		//
			if cnt < 2 { badgame() }

			boldface = true

			/* advance to next argument */
			ai += 1
			cnt -= 1
			continue
		}

		/* -b1 (Bet 1 Chip) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'b' &&
		   os.Args[1+ai][2] == '1' &&
		   len(os.Args[1+ai]) == 3 {

			if cnt < 2 { badgame() }

			/* set minimum bet */
			minbet = 1; bet = 1

			/* advance to next argument */
			ai += 1
			cnt -= 1
			continue
		}

		/* -g <name>  Choose Game */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'g' &&
		   len(os.Args[1+ai]) == 2 {
		//
			if cnt < 3 { badgame() }
			arg = os.Args[2+ai]

			if arg == "jb95" {
				game = JacksOrBetter95
			} else if arg == "jb86" {
				game = JacksOrBetter86
			} else if arg == "jb85" {
				game = JacksOrBetter85
			} else if arg == "jb75" {
				game = JacksOrBetter75
			} else if arg == "jb65" {
				game = JacksOrBetter65
			} else if arg == "aa" {
				game = AllAmerican
			} else if arg == "10s" {
				game = TensOrBetter
			} else if arg == "tens" {
				game = TensOrBetter
/*
			} else if arg == "bon" {
				game = BonusPoker
			} else if arg == "db" {
				game = DoubleBonus
			} else if arg == "dbb" {
				game = DoubleBonusBonus
*/
			} else { badgame() }

			setgame(game)

			/* advance to next argument */
			ai += 2
			cnt -= 2
			continue
		}

		/* -is (Initial Score) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'i' &&
		   os.Args[1+ai][2] == 's' &&
		   len(os.Args[1+ai]) == 3 {

			if cnt < 3 { badgame() }

			score, _ = strconv.Atoi(os.Args[2+ai])

			if score <= 0 || score > 100000  {
				fmt.Printf("Video Poker: bad number given with the -is option.\n")
				os.Exit(1)
			}

			if score%10 != 0 { minbet = 1; bet = 1 }

			/* advance to next argument */
			ai += 2
			cnt -= 2
			continue
		}

		/* -k <5-char-string>  Redefine input keys (default is " jkl;" */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'k' &&
		   len(os.Args[1+ai]) == 2 {
		//
			if cnt < 3 { badgame() }
			arg = os.Args[2+ai]

			if(len(arg) != NUMKEYS) {
				// complain and exit
				fmt.Printf("Video Poker: the string given with the -k option is the wrong length.\n")
				os.Exit(1)
			}

			/* copy the string into keys[] */
			for i = 0; i < NUMKEYS; i++ {

				if arg[i] == 'q' || arg[i] == 'e' {
					fmt.Printf("Video Poker: for the -k option, the string may not contain q or e.\n")
					os.Exit(1)
				}
				keys[i] = byte(arg[i])
			}

			/* advance to next argument */
			ai += 2
			cnt -= 2
			continue
		}

		/* -mh (Mark Held) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'm' &&
		   os.Args[1+ai][2] == 'h' &&
		   len(os.Args[1+ai]) == 3 {
		//
			if cnt < 2 { badgame() }

			/* turn on Mark Held flag */
			markheld = true

			/* advance to next argument */
			ai += 1
			cnt -= 1
			continue
		}

		/* -q (Quiet) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'q' &&
		   len(os.Args[1+ai]) == 2 {
		//
			if cnt < 2 { badgame() }

			/* turn on Quiet flag */
			quiet = true

			/* advance to next argument */
			ai += 1
			cnt -= 1
			continue
		}

		/* -u<n>  Unicode Output */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'u' &&
		   len(os.Args[1+ai]) > 2 {

			switch os.Args[1+ai][2] {

				case '0': unicode = UNICODE_TTY
				case '1': unicode = UNICODE_OFF
				case '2': unicode = UNICODE_SUITS
				case '3': unicode = UNICODE_CARDS
				default:
					fmt.Printf("Video Poker: digit %d given with -u option is out of range.\n",os.Args[1+ai][2]-'0')
					os.Exit(1)
			}
			/* advance to next argument */
			ai += 1
			cnt--
			continue
		}

		/* -v (Version) */

		if os.Args[1+ai][0] == '-' &&
		   os.Args[1+ai][1] == 'v' &&
		   len(os.Args[1+ai]) == 2 {
		//
			if cnt < 2 { badgame() }

			fmt.Printf("%s\n",VERSION)
			os.Exit(0)
		}

		/* unrecognized option */
		fmt.Printf("Video Poker: the %s option was not recognized\n",os.Args[1+ai])
		os.Exit(1)
	}

	/* Before starting play, print the name of the game in green */

	if ! quiet {
	//
		fmt.Printf("\n ")
		color(GREEN)
		ANSIbold()
		fmt.Printf("%s",gamenames[game])
		ANSIreset()
		fmt.Printf("\n\n")
	}

	for { play() }
}

/* Functions that recognize winning hands */

/*
	Flush:
	returns true if the sorted hand is a flush
*/

func flush() bool {
//
	if shand[0].suit == shand[1].suit &&
	   shand[1].suit == shand[2].suit &&
	   shand[2].suit == shand[3].suit &&
	   shand[3].suit == shand[4].suit { return true }

	return false
}

/*
	Straight:
	returns true if the sorted hand is a straight
*/

func straight() bool {
//
	if shand[1].index == shand[0].index + 1 &&
	   shand[2].index == shand[1].index + 1 &&
	   shand[3].index == shand[2].index + 1 &&
	   shand[4].index == shand[3].index + 1 { return true }

	return false
}

/*
	Four of a kind:
	the middle 3 all match, and the first or last matches those
*/

func four() bool {
//
	if (shand[1].index == shand[2].index &&
	    shand[2].index == shand[3].index ) &&
	   ( shand[0].index == shand[2].index ||
	     shand[4].index == shand[2].index) { return true }

	return false
}

/*
	Full house:
	3 of a kind and a pair
*/

func full() bool {
//
	if shand[0].index == shand[1].index &&
	  (shand[2].index == shand[3].index &&
	   shand[3].index == shand[4].index) { return true }

	if shand[3].index == shand[4].index &&
	  (shand[0].index == shand[1].index &&
	   shand[1].index == shand[2].index) { return true }

	return false
}

/*
	Three of a kind:
	it can appear 3 ways
*/

func three() bool {
//
	if shand[0].index == shand[1].index &&
	   shand[1].index == shand[2].index { return true }

	if shand[1].index == shand[2].index &&
	   shand[2].index == shand[3].index { return true }

	if shand[2].index == shand[3].index &&
	   shand[3].index == shand[4].index { return true }

	return false
}

/*
	Two pair:
	it can appear in 3 ways
*/

func twopair() bool {
//
	if ((shand[0].index == shand[1].index) && (shand[2].index == shand[3].index)) ||
	   ((shand[0].index == shand[1].index) && (shand[3].index == shand[4].index)) ||
	   ((shand[1].index == shand[2].index) && (shand[3].index == shand[4].index)) { return true }

	return false
}

/*
	Two of a kind (pair), jacks or better
	or if the game is Tens or Better, 10s or better.
*/

func two() bool {
//
	var min int = JACK

	if game == TensOrBetter { min = TEN }

	if shand[0].index == shand[1].index && shand[1].index >= min { return true }
	if shand[1].index == shand[2].index && shand[2].index >= min { return true }
	if shand[2].index == shand[3].index && shand[3].index >= min { return true }
	if shand[3].index == shand[4].index && shand[4].index >= min { return true }

	return false
}

/*
	This bunch of consts is used to index into
	paytable[] and handname[], so make sure the two match.
*/

const (
	ROYAL = iota
	STRFL
	FOUR
	FULL
	FLUSH
	STR
	THREE
	TWOPAIR
	PAIR
	NOTHING
	/* the number of the above: */
	NUMHANDTYPES
	)

var paytable [NUMHANDTYPES]int = [NUMHANDTYPES]int {
	800,	/* royal flush: 800 */
	50,	/* straight flush: 50 */
	25,	/* 4 of a kind: 25 */
	9,	/* full house: 9 */
	6,	/* flush: 6 */
	4,	/* straight: 4 */
	3,	/* 3 of a kind: 3 */
	2,	/* two pair: 2 */
	1,	/* jacks or better: 1 */
	0,	/* nothing */
}

var handname [NUMHANDTYPES]string = [NUMHANDTYPES]string {
	"Royal Flush    ",
	"Straight Flush ",
	"Four of a Kind ",
	"Full House     ",
	"Flush          ",
	"Straight       ",
	"Three of a Kind",
	"Two Pair       ",
	"Pair           ",
	"Nothing        ",
}

const INVALID = 100	/* higher than any valid card index */

/* returns type of hand */

func recognize() int {
//
	var i, j, f int
	var min int = INVALID
	var tmp [CARDS]card
	var st, fl bool	/* both are auto-initialized to 0 */

	/* Sort hand into sorted hand (shand) */

	/* make copy of hand */
	for i = 0; i < CARDS; i++ { tmp[i] = hand[i] }

	for i = 0; i < CARDS; i++ {
		/* put lowest card in hand into next place in shand */

		for j = 0; j < CARDS; j++ {
			if tmp[j].index <= min {
				min = tmp[j].index
				f = j
			}
		}

		shand[i] = tmp[f]
		tmp[f].index = INVALID	/* larger than any card */
		min = INVALID
	}

	/* royal and straight flushes, straight, and flush */

	fl = flush()
	st = straight()

	if st && fl && shand[0].index == TEN { return ROYAL }
	if st && fl { return STRFL }
	if four() { return FOUR }
	if full() { return FULL }
	if fl { return FLUSH }
	if st { return STR }
	if three() { return THREE }
	if twopair() { return TWOPAIR }
	if two() { return PAIR }

	/* Nothing */

	return NOTHING
}

/* The loop */

func play() {
//
	var i int
	var crd int
	var c int
	var hold[CARDS] int
	var digit int

	/* initialize deck */
	for i = 0; i < CARDSINDECK; i++ { deck[i].gone = 0 }

	/* initialize hold[] */
	for i = 0; i < CARDS; i++ { hold[i] = 0 }

	score -= bet

	for i = 0; i < CARDS; i++ {
	//
		/* find a card not already dealt */

		for crd = random()%CARDSINDECK; deck[crd].gone != 0 ; crd = random()%CARDSINDECK { }

		deck[crd].gone = 1
		hand[i] = deck[crd]
	}


	showhand()

	/* get cards to hold, and replace others */

	for {
	//
		c = getchar()
		if c == '\n' { break }

		if c == 'q' || c == 'e' {
		//
			boldface = false
			ANSIreset()

			if ! quiet {
				fmt.Printf("\nYou quit with %d chips after playing %d hands.\n",score+bet,hands)
				fmt.Printf("Range: %d - %d\n", score_low, score_high)
			}
			os.Exit(0)
		}

		if c == 'b' {	/* Change the bet. Only 1, 2, 3, 4, and 5 are allowed. */
		//
			digit = getchar()
			if digit <= '1' || digit >= '6' {
			//
				ungetc()
			} else {
			//
				betmultiplier = digit - '0'
				bet = betmultiplier * minbet
			}
			continue
		}

		for i = 0; i < NUMKEYS; i++ {
		//
			if int(keys[i]) == c {
			//
				/* flip bit to hold/discard it */
				hold[i] ^= 1
			}
		}
	}

	/* Optional Line: mark held cards */

	if markheld {
	//
		for i = 0; i < CARDS; i++ {
		//
			var pm string;
			if hold[i] != 0 { pm = " +" } else { pm = "  " }
			fmt.Printf("%s ", pm)
		}
		fmt.Printf("\n")
	}

	/* replace cards not held */

	for i = 0; i < CARDS; i++ {
	//
		if hold[i] == 0 {
		//
			for crd = random()%CARDSINDECK; deck[crd].gone != 0; crd = random()%CARDSINDECK { }

			deck[crd].gone = 1
			hand[i] = deck[crd]
		}
	}

	/* print final hand */

	showhand()

	/* recognize and score hand */

	i = recognize()

	score += paytable[i] * bet

	fmt.Printf("%s  ",handname[i])
	fmt.Printf("%d\n\n",score)

	hands++

	if score < score_low  { score_low  = score }
	if score > score_high { score_high = score }

	if score < bet {
	//
		for ; score < bet && betmultiplier > 1; {
		//
			betmultiplier--;
			bet = minbet * betmultiplier
		}

		if score < bet {
		//
			boldface = false
			ANSIreset()
			if ! quiet {
			//
				fmt.Printf("You ran out of chips after playing %d hands.\n", hands)
				if(score_high > INITCHIPS) { fmt.Printf("At one point, you had %d chips.\n", score_high) }
			}
			os.Exit(0)
		} else {
		//
			fmt.Printf("You are low on chips. Your bet has been reduced to %d.\n\n", bet)
		}
	}
}

/* do the work for the -g option */

func setgame(game int) {
//
	switch game {
	//
		case JacksOrBetter95:
			paytable[FLUSH] = 5
		case JacksOrBetter86:
			paytable[FULL] = 8
		case JacksOrBetter85:
			paytable[FULL] = 8
			paytable[FLUSH] = 5
		case JacksOrBetter75:
			paytable[FULL] = 7
			paytable[FLUSH] = 5
		case JacksOrBetter65:
			paytable[FULL] = 6
			paytable[FLUSH] = 5
		case AllAmerican:
			paytable[FULL] = 8
			paytable[FLUSH] = 8
			paytable[STR] = 8
			paytable[PAIR] = 1
		case TensOrBetter:
			/* pay table same as JacksOrBetter65 */
			paytable[FULL] = 6
			paytable[FLUSH] = 5
/*
		case BonusPoker:
			fmt.Printf("Bonus Poker is unimplemented in this version.\n")
			os.Exit(0)
		case DoubleBonus:
			fmt.Printf("Double Bonus is unimplemented in this version.\n")
			os.Exit(0)
		case DoubleBonusBonus:
			fmt.Printf("Double Bonus Bonus is unimplemented in this version.\n")
			os.Exit(0)
*/
	}
}
