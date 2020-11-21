package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot-specific command parameters.
const (
	cmdPrefix     = "!"
	charSearchStr = "char"
	raidSearchStr = "raid"
	helpStr       = "help"
)

// Command strings that a user's message will start with.
var (
	charSearchCmd = fmt.Sprintf("%s%s", cmdPrefix, charSearchStr)
	raidSearchCmd = fmt.Sprintf("%s%s", cmdPrefix, raidSearchStr)
	helpCmd       = fmt.Sprintf("%s%s", cmdPrefix, helpStr)
)

// When registered with `s`, this function will examine newly-created
// user messages. If the user message is found to contain a bot command,
// that command's logic will be invoked.
func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	isMsgFromSelf := m.Author.ID == s.State.User.ID
	if isMsgFromSelf {
		return
	}

	userMsg := strings.ToLower(m.Content)

	isHelpMsg := strings.HasPrefix(userMsg, helpCmd)
	if isHelpMsg {
		s.ChannelMessageSend(m.ChannelID, generateHelpMsg())
		return
	}

	isCharSearchMsg := strings.HasPrefix(userMsg, charSearchCmd)
	if isCharSearchMsg {
		s.ChannelMessageSend(m.ChannelID, generateUserFacingMsg(userMsg, generateCharSearchMsg))
		return
	}

	isRaidSearchMsg := strings.HasPrefix(userMsg, raidSearchCmd)
	if isRaidSearchMsg {
		s.ChannelMessageSend(m.ChannelID, generateUserFacingMsg(userMsg, generateRaidSearchMsg))
		return
	}
}

// generateHelpMsg returns an instructional message that can be passed
// back to a user.
func generateHelpMsg() string {
	helpMsg := fmt.Sprintf("Example: %s bfa asmongold us kel'thuzad\n", raidSearchCmd)
	return helpMsg
}

// generateUserFacingMsg is a wrapper function that invokes
// `msgGenerator` on `userInput`, returning the resulting message from
// `msgGenerator` if successful.
func generateUserFacingMsg(userInput string, msgGenerator func(string) (string, error)) string {
	msg, err := msgGenerator(userInput)
	if err != nil {
		return err.Error()
	}
	return msg
}

// apiCharSearchData encapsulates the data required in sending a query
// to Blizzard's character summary API.
type apiCharSearchData struct {
	characterSlug string
	serverRegion  string
	realmSlug     string
}

// apiRaidSearchData encapsulates the data required in sending a query
// to Blizzard's raid encounters API.
type apiRaidSearchData struct {
	expansionName string
	apiCharSearchData
}

// extractIndividualInputArgs attempts to break `query` into a
// `[]string` which represents individual arguments. An error is
// returned if `query` cannot be parsed into an `expectedArgsCount`
// number of arguments. `cmdToIgnore` represents text that should not be
// considered part of the user's unique input.
func extractIndividualInputArgs(query, cmdToIgnore string, expectedArgsCount int) ([]string, error) {
	// Strip the cmdToIgnore text from args since we do not need to process it anymore.
	argsWithoutCmd := strings.Replace(query, cmdToIgnore, "", 1)
	// Interpret user input as a space-separated list.
	// An argument wrapped in quotes will not be split.
	argsSplitterRegex := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	args := argsSplitterRegex.FindAllString(argsWithoutCmd, -1)
	if args == nil {
		return nil, errors.New("could not parse user input")
	}
	inputArgsCountInvalid := len(args) != expectedArgsCount
	if inputArgsCountInvalid {
		return nil, errors.New("could not parse user input")
	}
	return args, nil
}

// prepareRealmSlugForBlizzardAPI ensures `realmSlug` is correctly
// formatted in the way that the Blizzard API expects it to be.
func prepareRealmSlugForBlizzardAPI(realmSlug string) string {
	realmSlug = strings.Trim(realmSlug, "\"")
	if strings.Contains(realmSlug, " ") {
		realmSlug = strings.ReplaceAll(realmSlug, " ", "-")
	}
	if strings.Contains(realmSlug, "'") {
		realmSlug = strings.ReplaceAll(realmSlug, "'", "")
	}
	return realmSlug
}

// parseRaidSearchQuery analyzes `query` and attempts to return data
// that can be used in an HTTP request to the Blizzard encounter summary
// API.
func parseRaidSearchQuery(query string) (apiRaidSearchData, error) {
	// expected form: !raid <expansion name> <character name> <server region> "<server name>"
	expectedArgsCount := 4
	args, err := extractIndividualInputArgs(query, raidSearchCmd, expectedArgsCount)
	if err != nil {
		return apiRaidSearchData{}, err
	}
	realmSlug := prepareRealmSlugForBlizzardAPI(args[expectedArgsCount-1])
	apiRaidSearchData := apiRaidSearchData{
		expansionName: args[0],
		apiCharSearchData: apiCharSearchData{
			characterSlug: args[1],
			serverRegion:  args[2],
			realmSlug:     realmSlug,
		},
	}
	return apiRaidSearchData, nil
}

// generateRaidSearchMsg attempts to query the Blizzard encounter
// summary API using the information obtained from `query`, returning a
// message that is ready to send back to the user who issued the bot
// command.
func generateRaidSearchMsg(query string) (string, error) {
	apiRaidSearchData, err := parseRaidSearchQuery(query)
	if err != nil {
		log.Println(err)
		return "", err
	}
	accessToken, err := getAccessTokenFromBlizzardAPI(apiRaidSearchData.serverRegion)
	if err != nil {
		log.Println(err)
		return "", err
	}
	msg, err := getRaidProfileFromBlizzardAPI(apiRaidSearchData, accessToken)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return msg, nil
}

// parseCharSearchQuery analyzes `query` and attempts to return data
// that can be used in an HTTP request to the Blizzard character summary
// API.
func parseCharSearchQuery(query string) (apiCharSearchData, error) {
	expectedArgsCount := 3
	args, err := extractIndividualInputArgs(query, charSearchCmd, expectedArgsCount)
	if err != nil {
		return apiCharSearchData{}, err
	}
	realmSlug := prepareRealmSlugForBlizzardAPI(args[expectedArgsCount-1])
	apiCharSearchData := apiCharSearchData{
		characterSlug: args[0],
		serverRegion:  args[1],
		realmSlug:     realmSlug,
	}
	return apiCharSearchData, nil
}

// generateCharSearchMsg attempts to query the Blizzard character
// summary summary API using the information obtained from `query`,
// returning a message that is ready to send back to the user who issued
// the bot command.
func generateCharSearchMsg(query string) (string, error) {
	apiCharSearchData, err := parseCharSearchQuery(query)
	if err != nil {
		log.Println(err)
		return "", err
	}
	accessToken, err := getAccessTokenFromBlizzardAPI(apiCharSearchData.serverRegion)
	if err != nil {
		log.Println(err)
		return "", err
	}
	msg, err := getCharProfileFromBlizzardAPI(apiCharSearchData, accessToken)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return msg, nil
}
