package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.connectToGemini()
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
	stream, err := client.Live.Connect(a.ctx, "gemini-2.5-flash-native-audio-preview-12-2025", &config)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	for {
		resp, err := stream.Receive()
		if err != nil {
			log.Fatal(err)
		}
		if resp.SetupComplete != nil {
			fmt.Printf("Setup complete: %#v\n", resp.SetupComplete)
			continue
		}
		if resp.GoAway != nil {
			fmt.Printf("Go away: %s\n", resp.GoAway.TimeLeft)
			break
		}
		fmt.Printf("unhandled message: %#v\n", resp)
		break
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
