package main

import (
	"regexp"
	"strings"
)

var rulesTextClassifiers = map[string]*regexp.Regexp{}
var typeLineClassifiers = map[string]*regexp.Regexp{}
var isCreature = regexp.MustCompile("(?i)^(summon|(legendary )?creature)")

func classify(card *mtgv4Card) {
	if card.CustomClassification == nil {
		card.CustomClassification = []string{}
	}
	var tags = map[string]interface{}{}
	for tag, rx := range rulesTextClassifiers {
		if rx.MatchString(card.Text) {
			tags[tag] = nil
		}
	}
	for tag, rx := range typeLineClassifiers {
		if rx.MatchString(card.Type) {
			tags[tag] = nil
		}
	}
	for tag := range tags {
		if tag == "commander" {
			if card.Legalities.Commander == "" || !isCreature.MatchString(card.Type) {
				continue
			}
		}
		card.CustomClassification = append(card.CustomClassification, tag)
	}

}

func init() {
	typeLineClassifiers["commander"] = regexp.MustCompile(
		"(?i)^(Summon Legend|Legendary Creature)",
	)
	rulesTextClassifiers["commander"] = regexp.MustCompile(
		"(?i)(can be your commander)",
	)
	rulesTextClassifiers["ramp"] = regexp.MustCompile(
		strings.Join(
			[]string{
				"(?i)((^| )(add {[^}]+}",
				"play[^.]+additional[^.]+lands?",
				"put[^.]+land[^.]*onto the battlefield",
				"search[^.]+your library for[^.]+(forest|swamp|island|plains|mountain|land)s?)( |\\.|$))",
			},
			"|",
		),
	)
	rulesTextClassifiers["draw"] = regexp.MustCompile(
		"(?i)((^| )(cards? into your hand|draw (a|x|that many|one|\\d+) cards?)( |\\.|$))",
	)
	rulesTextClassifiers["board-wipes"] = regexp.MustCompile(
		"(?i)((^| )(sacrifice|destroy|exile) all)( |\\.|$)",
	)
	rulesTextClassifiers["tribe-support"] = regexp.MustCompile(
		strings.Join(
			[]string{
				// needs work...
				"(?i)((^| )((addition|is|instances?|hooses?|shares?|all)[^.]+creature types?|creature[^.]+choice)( |\\.|$))",
			},
			"|",
		),
	)
	rulesTextClassifiers["counters"] = regexp.MustCompile(
		strings.Join(
			[]string{
				"(?i)((^| )(proliferates?",
				"((any|another) kind of)? counters? (from|to|on|already)?",
				"Bolster \\d+",
				"Suspend \\d+",
				"{E}",
				// begin types of counters
				"([+-]\\d+/[+-]\\d+",
				"Age",
				"Aim",
				"Arrow",
				"Arrowhead",
				"Awakening",
				"Blaze",
				"Blood",
				"Bounty",
				"Bribery",
				"Brick",
				"Carrion",
				"Charge",
				"CRANK!",
				"Credit",
				"Corpse",
				"Crystal",
				"Cube",
				"Currency",
				"Death",
				"Delay",
				"Depletion",
				"Despair",
				"Devotion",
				"Divinity",
				"Doom",
				"Dream",
				"Echo",
				"Egg",
				"Elixir",
				"Energy",
				"Eon",
				"Experience",
				"Eyeball",
				"Fade",
				"Fate",
				"Feather",
				"Filibuster",
				"Flood",
				"Fungus",
				"Fuse",
				"Gem",
				"Glyph",
				"Golds",
				"Growth",
				"Hatchling",
				"Healing",
				"Hit",
				"Hoofprint",
				"Hour",
				"Hourglass",
				"Hunger",
				"Ice",
				"Incubation",
				"Infection",
				"Intervention",
				"Isolation",
				"Javelin",
				"Ki",
				"Level",
				"Lore",
				"Loyalty",
				"Luck",
				"Magnet",
				"Manifestation",
				"Mannequin",
				"Mask",
				"Matrix",
				"Mine",
				"Mining",
				"Mire",
				"Music",
				"Muster",
				"Net",
				"Omen",
				"Ore",
				"Page",
				"Pain",
				"Paralyzation",
				"Petal",
				"Petrification",
				"Phylactery",
				"Pin",
				"Plague",
				"Poison",
				"Polyp",
				"Pressure",
				"Prey",
				"Pupa",
				"Quest",
				"Rust",
				"Scream",
				"Shell",
				"Shield",
				"Silver",
				"Shred",
				"Sleep",
				"Sleight",
				"Slime",
				"Slumber",
				"Soot",
				"Spore",
				"Storage",
				"Strife",
				"Study",
				"Theft",
				"Tide",
				"Time",
				"Tower",
				"Training",
				"Trap",
				"Treasure",
				"Velocity",
				"Verse",
				"Vitality",
				"Volatile",
				"Wage",
				"Winch",
				"Wind",
				"Wish) counters?)( |\\.|$))",
			},
			"|",
		),
	)
}
