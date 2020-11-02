package main

const (
	classicKey = iota
	theBurningCrusadeKey
	wrathOfTheLichKingKey
	cataclysmKey
	mistsOfPandariaKey
	warlordsOfDraenorKey
	legionKey
	battleForAzerothKey
	// Todo: add key for Shadowlands when it is added to the Blizzard API
)

var expansionKeys = map[string]int{
	"classic":   classicKey,
	"vanilla":   classicKey,
	"tbc":       theBurningCrusadeKey,
	"bc":        theBurningCrusadeKey,
	"wrath":     wrathOfTheLichKingKey,
	"wotlk":     wrathOfTheLichKingKey,
	"lich king": wrathOfTheLichKingKey,
	"cataclysm": cataclysmKey,
	"cata":      cataclysmKey,
	"mop":       mistsOfPandariaKey,
	"mists":     mistsOfPandariaKey,
	"pandaria":  mistsOfPandariaKey,
	"wod":       warlordsOfDraenorKey,
	"warlords":  warlordsOfDraenorKey,
	"draenor":   warlordsOfDraenorKey,
	"legion":    legionKey,
	"bfa":       battleForAzerothKey,
	"battle":    battleForAzerothKey,
}

func getExpansionKey(expansionAlias string) (int, bool) {
	key, ok := expansionKeys[expansionAlias]
	return key, ok
}
