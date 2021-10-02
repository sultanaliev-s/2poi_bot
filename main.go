package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var qwertyLtoc map[rune]rune
var qwertyCtol map[rune]rune

func populateQwertys() {
	qwertyLtoc = make(map[rune]rune)
	qwertyCtol = make(map[rune]rune)

	eng := []rune(
		"qwertyuiop[]asdfghjkl;'zxcvbnm,./QWERTYUIOP{}ASDFGHJKL:\"ZXCVBNM<>?&")
	rus := []rune(
		"йцукенгшщзхъфывапролджэячсмитьбю.ЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,?")

	for i, v := range eng {
		qwertyLtoc[v] = rune(rus[i])
	}
	for i, v := range rus {
		qwertyCtol[v] = rune(eng[i])
	}
}

// func init() {
// 	// loads values from .env into the system
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No .env file found")
// 	}
// }

func main() {
	populateQwertys()

	port := os.Getenv("PORT")

	if len(port) != 0 {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}

	botToken := os.Getenv("BOT_TOKEN")
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	if len(botToken) != 0 {
		log.Print("token found")
	}
	if len(botToken) == 0 {
		log.Print("token not found")
	}

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Something went wrong:", err)
		}

		for _, update := range updates {
			err = respond(botUrl, update)
			if err != nil {
				log.Println("Something went wrong:", err)
			}
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getupdates" +
		"?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func respond(botUrl string, update Update) error {
	// if !shouldRespond(update) {
	// 	return nil
	// }

	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = "TUP" //translate(&update.Message.ReplyToMessage.Text)

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	_, err = http.Post(
		botUrl+"/sendMessage",
		"application/json",
		bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	return nil
}

func shouldRespond(update Update) bool {
	if update.Message.ReplyToMessage == nil {
		return false
	}
	if !strings.Contains(update.Message.Text, "@TwoPoiBot") {
		return false
	}
	return true
}

func translate(str *string) (res string) {
	text := []rune(*str)

	if len(text) > 0 && isCyrillic(text[0]) {
		return translateSpec(&text, &qwertyCtol)
	}

	return translateSpec(&text, &qwertyLtoc)
}

func isCyrillic(char rune) bool {
	return unicode.Is(unicode.Cyrillic, char)
}

func translateSpec(text *[]rune, dict *map[rune]rune) (res string) {
	for _, v := range *text {
		ch, ok := (*dict)[v]
		if ok {
			res += string(ch)
		} else {
			res += string(v)
		}
	}

	return res
}
