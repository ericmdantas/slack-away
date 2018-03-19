package away

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nlopes/slack"
)

const (
	botID    = "U9N01EY01"
	botIDMsg = "<@U9N01EY01>"
)

func Start() {
	client := slack.New(os.Getenv("SLACK_BOT_AWAY_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot-away: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	client.SetDebug(true)

	rtm := client.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go mightCreateNewConversationAfterTime(client, rtm, ev)
			replyOrStartNewConversation(client, rtm, ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}

func replyOrStartNewConversation(client *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent) {
	for _, trigger := range wordsThatTriggerReply {
		if caseInsensitiveContains(ev.Text, trigger) {
			replyToUser(rtm, ev)
			break
		} else {
			mightReplyToConversation(client, rtm, ev)
		}
	}
}

func replyToUser(rtm *slack.RTM, ev *slack.MessageEvent) {
	if fmt.Sprintf("<@%s>", ev.User) == botIDMsg {
		return
	}

	userFormatted := formatUser(ev.User)

	rand.Seed(time.Now().UnixNano())
	msg := fmt.Sprintf("%s %s", userFormatted, freeOffenses[rand.Intn(len(freeOffenses))])
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
}

func mightCreateNewConversationAfterTime(client *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent) {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		log.Println("Channel is idle, might reply to conversation")
		mightReplyToConversation(client, rtm, ev)
	}
}

func mightReplyToConversation(client *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent) {
	rand.Seed(time.Now().UnixNano())
	max := 999
	div := 333
	result := rand.Intn(max) % div
	shouldInteract := result == 0

	fmt.Printf("----> interact: %t -- %d possible numbers with %d chance of happening -- result: %d\n", shouldInteract, max, max/div, result)

	if shouldInteract {
		channel, err := client.GetChannelInfo(ev.Channel)

		log.Printf("----------> Interacting with channel: %s\n", channel.Name)

		if err != nil {
			fmt.Printf("Error retrieving channel info: %s\n", err.Error())
			return
		}

		member := channel.Members[rand.Intn(len(channel.Members))]

		if member == botID {
			return
		}

		user, err := client.GetUserInfo(member)

		if err != nil {
			fmt.Printf("Error retrieving user info: %s\n", err.Error())
			return
		}

		log.Printf("------------> Interacting with %s\n", user.Name)

		msg := fmt.Sprintf("<@%s> %s", user.ID, freeOffenses[rand.Intn(len(freeOffenses))])
		rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
	}
}
