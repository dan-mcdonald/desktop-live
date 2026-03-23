package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

// App struct
type App struct {
	ctx   context.Context
	fe2be chan FrontendMessage
	be2fe chan *genai.LiveServerMessage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fe2be: make(chan FrontendMessage),
		be2fe: make(chan *genai.LiveServerMessage),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.connectToGemini()
}

func geminiReceiveServerMessages(session *genai.Session, be2fe chan *genai.LiveServerMessage) {
	for {
		resp, err := session.Receive()
		if err != nil {
			log.Fatal(err)
		}
		be2fe <- resp
	}
}

type FrontendMessage struct {
	Text string
}

func (a *App) PostMessage(msg FrontendMessage) {
	a.fe2be <- msg
}

func (a *App) PollMessage() *genai.LiveServerMessage {
	return <-a.be2fe
}

func (a *App) connectToGemini() {
	client, err := genai.NewClient(a.ctx, &genai.ClientConfig{})
	if err != nil {
		log.Fatal(err)
	}

	config := genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{
			genai.ModalityAudio,
		},
		MediaResolution: genai.MediaResolutionMedium,
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: "Zubenelgenubi",
				},
			},
		},
		ContextWindowCompression: &genai.ContextWindowCompressionConfig{
			TriggerTokens: new(int64(104857)),
			SlidingWindow: &genai.SlidingWindow{
				TargetTokens: new(int64(52428)),
			},
		},
	}
	session, err := client.Live.Connect(a.ctx, "gemini-2.5-flash-native-audio-preview-12-2025", &config)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	go geminiReceiveServerMessages(session, a.be2fe)

	for feMsg := range a.fe2be {
		session.SendRealtimeInput(genai.LiveRealtimeInput{
			Text: feMsg.Text,
		})
	}
}

// domReady is called when the DOM is ready
func (a *App) domReady(ctx context.Context) {
	fmt.Println("DOM is ready")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
