package leaderboard

import (
	"log"

	dsedb "dont-slack-evil/db"
)

// SendLeaderboardNotification sends the leaderboard notification
func SendLeaderboardNotification() (int, error) {
	log.Println("Hi Im here")
	log.Println("mdr")
	notificationsSent := 0
	log.Println("LOL")
	log.Println(dsedb.GetTeams)
	teams, teamsErr := dsedb.GetTeams()
	if teamsErr != nil {
		log.Println(teamsErr)
		return 0, teamsErr
	}
	log.Println(teams)
	return notificationsSent, nil
}