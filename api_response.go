package main

import "fmt"

type charProfileData struct {
	Name    string `json:"name"`
	Faction struct {
		Name string `json:"name"`
	} `json:"faction"`
	Race struct {
		Name string `json:"name"`
	} `json:"race"`
	Class struct {
		Name string `json:"name"`
	} `json:"character_class"`
	Specialization struct {
		Name string `json:"name"`
	} `json:"active_spec"`
	Guild struct {
		Name string `json:"name"`
	} `json:"guild"`
	Level     int `json:"level"`
	ItemLevel int `json:"average_item_level"`
}

func (data charProfileData) String() string {
	usableGuildName := "none"
	if data.Guild.Name != "" {
		usableGuildName = data.Guild.Name
	}
	msg := fmt.Sprintln(
		"**Character:**",
		data.Name,
		data.Level,
		data.Faction.Name,
		data.Race.Name,
		data.Specialization.Name,
		data.Class.Name,
		"\n**Guild:**",
		usableGuildName,
		"\n**Item level:**",
		data.ItemLevel,
	)
	return msg
}

type raidProfileData struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Character struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Realm struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
			Slug string `json:"slug"`
		} `json:"realm"`
	} `json:"character"`
	Expansions []struct {
		Expansion struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"expansion"`
		Instances []struct {
			Instance struct {
				Key struct {
					Href string `json:"href"`
				} `json:"key"`
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"instance"`
			Modes []struct {
				Difficulty struct {
					Type string `json:"type"`
					Name string `json:"name"`
				} `json:"difficulty"`
				Status struct {
					Type string `json:"type"`
					Name string `json:"name"`
				} `json:"status"`
				Progress struct {
					CompletedCount int `json:"completed_count"`
					TotalCount     int `json:"total_count"`
					Encounters     []struct {
						Encounter struct {
							Key struct {
								Href string `json:"href"`
							} `json:"key"`
							Name string `json:"name"`
							ID   int    `json:"id"`
						} `json:"encounter"`
						CompletedCount    int   `json:"completed_count"`
						LastKillTimestamp int64 `json:"last_kill_timestamp"`
					} `json:"encounters"`
				} `json:"progress"`
			} `json:"modes"`
		} `json:"instances"`
	} `json:"expansions"`
}

func (data raidProfileData) ExpansionString(expansionKey int) string {
	instances := data.Expansions[expansionKey].Instances
	var msg string
	for _, instance := range instances {
		msg += fmt.Sprintf("**%s**\n", instance.Instance.Name)
		for _, mode := range instance.Modes {
			msg += fmt.Sprintf("• %s %d/%d\n", mode.Difficulty.Name, mode.Progress.CompletedCount, mode.Progress.TotalCount)
		}
		msg += "\n"
	}
	return msg
}
