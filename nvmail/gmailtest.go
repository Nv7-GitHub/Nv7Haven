package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
)

//go:embed serviceAccount.json
var serviceAccount []byte

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("nvmail/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	/*r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}*/

	/*r2, err := srv.Users.Messages.List(user).Q("label:UNREAD").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r2.Messages) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Messages:")
	for _, l := range r2.Messages {
		msg, err := srv.Users.Messages.Get(user, l.Id).Do()
		if err != nil {
			log.Fatalf("Unable to load message: %v", err)
		}
		for _, header := range msg.Payload.Headers {
			if header.Name == "Subject" {
				fmt.Println(header.Value)
			}
		}
		/*if msg.Payload.Body.Data != "" {
			data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(string(data))
		}
	}*/

	// Pubsub Create Topic
	fmt.Println("Creating pubsub client...")
	clt, err := pubsub.NewClient(context.Background(), "nvmail-1611539053087", option.WithCredentialsJSON(serviceAccount))
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer clt.Close()

	fmt.Println("Creating topic...")
	var mailTopic *pubsub.Topic
	topics := clt.Topics(context.Background())
	for {
		topic, err := topics.Next()
		if err != nil {
			break
		}
		if topic.ID() == "mail" {
			mailTopic = topic
		}
	}

	if mailTopic == nil {
		fmt.Println("Creating new topic...")
		var err error
		mailTopic, err = clt.CreateTopic(context.Background(), "mail")
		if err != nil {
			panic(err)
		}
	}

	// Listen to mail topic
	fmt.Println("Creating pubsub listener...")
	var sub *pubsub.Subscription
	subs := clt.Subscriptions(context.Background())
	for {
		subS, err := subs.Next()
		if err != nil {
			break
		}
		if subS.ID() == "mailSub" {
			sub = subS
		}
	}

	if sub == nil {
		var err error
		sub, err = clt.CreateSubscription(context.Background(), "mailSub", pubsub.SubscriptionConfig{
			Topic: mailTopic,
		})
		if err != nil {
			panic(err)
		}
	}

	// Listen
	fmt.Println("Adding publisher...")
	srv.Users.Watch(user, &gmail.WatchRequest{
		LabelIds:  []string{"UNREAD"},
		TopicName: "projects/nvmail-1611539053087/topics/mail",
	}).Do()
	if err != nil {
		log.Fatalf("Unable to listen for mail: %v", err)
	}

	// Add Listener
	fmt.Println("Listening for messages...")
	ctx, cancel := context.WithCancel(context.Background())
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Println("Received message!")
		fmt.Println(string(m.Data))
	})
	if err != nil {
		panic(err)
	}
	defer cancel()
}
