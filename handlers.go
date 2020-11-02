package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	cmdPrefix     = "!"
	charSearchStr = "char"
	raidSearchStr = "raid"
	helpStr       = "help"
)

var (
	charSearchCmd = fmt.Sprintf("%s%s", cmdPrefix, charSearchStr)
	raidSearchCmd = fmt.Sprintf("%s%s", cmdPrefix, raidSearchStr)
	helpCmd       = fmt.Sprintf("%s%s", cmdPrefix, helpStr)
)

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
		s.ChannelMessageSend(m.ChannelID, "Searching character data...")
		msg, err := generateCharSearchMsg(userMsg)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	isRaidSearchMsg := strings.HasPrefix(userMsg, raidSearchCmd)
	if isRaidSearchMsg {
		s.ChannelMessageSend(m.ChannelID, "Searching raid data...")
		msg, err := generateRaidSearchMsg(userMsg)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
}

func generateHelpMsg() string {
	helpMsg := fmt.Sprintf("Example: %s bfa asmongold us kel'thuzad\n", raidSearchCmd)
	return helpMsg
}

type apiCharSearchData struct {
	characterSlug string
	serverRegion  string
	realmSlug     string
}

type apiRaidSearchData struct {
	expansionName string
	apiCharSearchData
}

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
