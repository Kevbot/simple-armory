package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func sendGetRequest(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &v)
}

func constructAccessTokenRequestURL(serverRegion string) string {
	url := fmt.Sprintf(
		"https://%s.battle.net/oauth/token",
		serverRegion,
	)
	queryParams := fmt.Sprintf(
		"?grant_type=client_credentials&client_id=%s&client_secret=%s",
		os.Getenv("WOW_CLIENT_ID"),
		os.Getenv("WOW_CLIENT_SECRET"),
	)
	return url + queryParams
}

func getAccessTokenFromBlizzardAPI(serverRegion string) (string, error) {
	if _, ok := findRegion(serverRegion); !ok {
		return "", errors.New("could not find server region")
	}
	accessToken := struct {
		Value string `json:"access_token"`
	}{}
	requestURL := constructAccessTokenRequestURL(serverRegion)
	if err := sendGetRequest(requestURL, &accessToken); err != nil {
		return "", err
	}
	return accessToken.Value, nil
}

func constructRaidProfileRequestURL(params apiRaidSearchData, region region, accessToken string) string {
	url := fmt.Sprintf(
		"https://%s.api.blizzard.com/profile/wow/character/%s/%s/encounters/raids",
		region.abbreviation,
		params.realmSlug,
		params.characterSlug,
	)
	queryParams := fmt.Sprintf(
		"?namespace=%s&locale=%s&access_token=%s",
		region.namespace,
		region.locale,
		accessToken,
	)
	return url + queryParams
}

func getRaidProfileFromBlizzardAPI(params apiRaidSearchData, accessToken string) (string, error) {
	region, ok := findRegion(params.serverRegion)
	if !ok {
		return "", errors.New("could not find server region")
	}
	requestURL := constructRaidProfileRequestURL(params, region, accessToken)
	var respData raidProfileData
	if err := sendGetRequest(requestURL, &respData); err != nil {
		return "", err
	}
	key, ok := getExpansionKey(params.expansionName)
	if !ok {
		return "", errors.New("could not find expansion")
	}
	characterHasRaidProgress := key < len(respData.Expansions)
	if !characterHasRaidProgress {
		return "", errors.New("character has no raid progress in this expansion")
	}
	return respData.ExpansionString(key), nil
}

func constructCharProfileRequestURL(params apiCharSearchData, region region, accessToken string) string {
	url := fmt.Sprintf(
		"https://%s.api.blizzard.com/profile/wow/character/%s/%s",
		region.abbreviation,
		params.realmSlug,
		params.characterSlug,
	)
	queryParams := fmt.Sprintf(
		"?namespace=%s&locale=%s&access_token=%s",
		region.namespace,
		region.locale,
		accessToken,
	)
	return url + queryParams
}

func getCharProfileFromBlizzardAPI(params apiCharSearchData, accessToken string) (string, error) {
	region, ok := findRegion(params.serverRegion)
	if !ok {
		return "", errors.New("could not find server region")
	}
	requestURL := constructCharProfileRequestURL(params, region, accessToken)
	var respData charProfileData
	if err := sendGetRequest(requestURL, &respData); err != nil {
		return "", err
	}
	return respData.String(), nil
}
