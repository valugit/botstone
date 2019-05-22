package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Card struct {
	Flavor  string `json:"flavor"`
	ImgGold string `json:"imgGold"`
}

type Infos struct {
	Collection []Card
}

func api(name string) (string, string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://omgvamp-hearthstone-v1.p.rapidapi.com/cards/%s", name), nil)
	if err != nil {
		fmt.Println("error on req")
	}
	req.Header.Add("X-RapidAPI-Host", "omgvamp-hearthstone-v1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", "3466702fb8msh7ef4d1d64f9e41ep183fd8jsnbf0761c6e744")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error on res")
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	infos := make([]Card, 0)
	_ = json.Unmarshal(body, &infos)

	if len(infos) > 1 {
		return infos[1].Flavor, infos[1].ImgGold
	}
	return infos[0].Flavor, infos[0].ImgGold
}

func main() {
	botStone, err := discordgo.New("Bot <YOUR KEY HERE>")

	if err != nil {
		fmt.Println("failed to create discord session", err)
	}

	if err != nil {
		fmt.Println("failed to access account", err)
	}

	botStone.AddHandler(handleCmd)
	err = botStone.Open()

	if err != nil {
		fmt.Println("unable to establish connection", err)
	}

	defer botStone.Close()

	<-make(chan struct{})
}

func handleCmd(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	user := msg.Author
	if user.ID == bot.State.User.ID || user.Bot {
		return
	}

	content := msg.Content

	if strings.HasPrefix(content, "!") {
		flavor, imgGold := api(strings.TrimPrefix(content, "!"))
		fmt.Println(flavor, imgGold)
		bot.ChannelMessageSend(msg.ChannelID, flavor)
		bot.ChannelMessageSend(msg.ChannelID, imgGold)
	}
}
