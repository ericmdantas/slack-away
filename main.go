package main

import (
	"fmt"
	"time"
	"os"
	"log"
	"math/rand"
	"strings"

	"github.com/nlopes/slack"
)

const (
	botID = "U9N01EY01"
	botIDMsg = "<@U9N01EY01>"
)

var (
	wordsThatTriggerReply = []string{
		"<@U9N01EY01>", // botIDMsg
		"away",
		"petropolis",
		"brother",
		"ronaldinho",
		"smith",
		"legal",
		"Boa sorte e que os jogos comecem", // jigsaw
		"Poucos builds hoje.", // cumulo.nimbus
		"Lamento por ter que ler isto", // cumulo.nimbus
		"Sai dae, doido! Cê é mó playboy!", // felipe.smith
	}
	
	freeOffenses = []string{
		"tu não aguenta dez minutos de porrada comigo",
		"teu pau ta expurgando, morô",
		"tu foi criado a leitinho com pêra, morô",
		"fica na tua aí",
		"legal é meu pau de paraquedas",
		"fica na tua aí senão vou te deitar na porrada",
		"FILHAAAAA DA PUTAAAAAA",
		"PARA COM ESSA PORRA, MERMÃO",
		"CAAAAAAAALA BOCA",		
		"MONICA DO CARALHO",
		"porra nenhuma, morô",
		"Ronaldinho, tu ta excroto",
		"tu sabe de nada, maluco",
		"tu é um inocente junvenil, morô",
		"Um bando de badernista, tudo uns aluno criado a leite com pêra, a ovo maltino, a pão com mortandela",
		"cumigu num tem caratê, num tem jorgite e nem cunguifu. A minha chinfra, maluco, aí, é botá os 5 dedo na cara e cair tremendo. Tô convidano QUALQUER UM PRA PORRADA, PRA PORRADA MERMO",
		"O aumento da gasolina que se foda-se mermão eu num tenho carro",
		"porque meu dinheiro, suor da minha testa, não foi feito, pra mim dar pra filho de político ir pras Ilhas Galáctas",
		"cala boca meu aluno, você está defecando pela boca!",
		"esse caso é complicado. É muita falcatrua com um pouquinho assim de vadiagem",
		"vou passar a lambida no pescoço",
		"cê qué ficá peitudo? Bundudo? A gatinha bunitinha vira pra você e cê vira a bundinha!",
		"já cozinhei pra Hebe Camargo com aquela cara pelancuda",
		"você tá pensanu que café dá em árvori, rapá?",
		"testada no pau do nariz, chapa nos peito, soco no coração… E soco na cabeça pra desentupir o célebro.",
		"você, garotinho inocente, garotinho juvenil, criado a leite com pêra, a Ovomaltino na geladeira…",
		"poxa, vâmo almoçá?",
		"estou perprecto",
		"então vai tomar no cu tranquilo",
		"calma aí… tô aqui desenvolvendo um raciocínio",
		"tudo por causa da… vadiagem!",
		"que Ritney Ritsu o caralho mermão!",
		"Vo te bota cinco dedo na cara e tu vai cair tremendo!",
		"seje feliz hein viadinho",
		"foda-se, mermão",
		"e daí?",
		"vai porra nenhuma, mermão",
		"tenho nada com isso não, morô",
		"fecha o cu pra falar comigo, morô",
		"vai criar porquinho na ilha, até cai dentinho",
		"cumpadi é o caralho! Num batizei teus filho porra!",
		"não tem dotô que dá jeito",
		"e essa rapaziada q ta ai te acompanhando ta tudo com o cu dando bote!",
		"tu ta com o cu dando bote",
		"não, meu amor... estou aqui desinvolvendo o meu raçocínio aqui no meu lépi toki... não tenho tempo, meu amor....",
		"o cachacero toma no cu hoje, e amanhã ele toma no cu de novo",
		"vo t'stala o coco rapá",
		"esse caso é complicado. É muita falcatrua com um pouquinho assim de vadiagem",
		"os indio tava dançando igual um filhadaputa todo barrigudo lá na aldeia",
		"...E você DIZ, antes que eu te passe a LAMBIDA: Viva a morte do meu pau, viva... a morte... do meu pau! Anabolizante filho da puta",
		"que que tem?",
		"ai delicia",
		"isso daí é poblema teu, mermão",
		"vou te levar pro meu sítio em São José e te colocar pra mamar nas cabrita",
		"bota uma dentadura no cu e sorri pro caralho",
		"dá um tempo, morô?",
	}
)

func replyToUser(rtm *slack.RTM, ev *slack.MessageEvent) {
	if fmt.Sprintf("<@%s>", ev.User) == botIDMsg {
		return
	}

	userFormatted := fmt.Sprintf("<@%s>", ev.User)
	
	if userFormatted == "<@>" {
		userFormatted = ""
	}
	
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
	
	fmt.Printf("Checking if there should be an interaction, %d possible numbers with %f.2 chance of happening -- result: %f.2, should interact: %t\n", max, max /div, result, shouldInteract)
	
	if shouldInteract {
		log.Println("------------> Should interact with someone")
		
		channel, err := client.GetChannelInfo(ev.Channel)		
		
		log.Printf("------------> Interacting with channel: %s\n", channel.Name)
		
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

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	client := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot-away: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	client.SetDebug(true)
	
	rtm := client.NewRTM()	
	go rtm.ManageConnection()
	
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go mightCreateNewConversationAfterTime(client, rtm, ev)
			
			for _, trigger := range wordsThatTriggerReply {
				if caseInsensitiveContains(ev.Text, trigger) {
					replyToUser(rtm, ev)
					break
				} else {
					mightReplyToConversation(client, rtm, ev)
				}
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}
	}
}