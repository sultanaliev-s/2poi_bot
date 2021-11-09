package main

import (
	"2poi_bot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

var conf *config.Config

var qwertyLtoc map[rune]rune
var qwertyCtol map[rune]rune

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf = config.New()
	populateQwertys()

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":"+conf.PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	run()
}

func run() {
	offset := 0
	for {
		updates, err := getUpdates(conf.BOT_URL, offset)
		if err != nil {
			log.Println("Something went wrong:", err)
		}

		for _, update := range updates {
			err = respond(conf.BOT_URL, update)
			if err != nil {
				log.Println("Something went wrong:", err)
			}
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	requestURL := fmt.Sprintf("%s/%s?%s=%d&?%s=%d&?%s=%s",
		botUrl, "getupdates", "offset", offset, "timeout=", 5,
		"allowed_updates=", "[\"message\"]")

	resp, err := http.Get(requestURL)
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
	if !shouldRespond(update) {
		return nil
	}

	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = translate(&update.Message.ReplyToMessage.Text) +
		"\n" + getQuote(&update.Message.ReplyToMessage.Sender)

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	requestURl := botUrl + "/sendMessage"
	_, err = http.Post(requestURl, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	return nil
}

func shouldRespond(update Update) bool {
	if update.Message.ReplyToMessage == nil {
		return false
	}
	return strings.Contains(update.Message.Text, "@TwoPoiBot")
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

func getQuote(sender *User) string {
	if len(sender.Username) == 0 {
		return "©" + sender.Name
	}
	return "©" + sender.Username
}

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
