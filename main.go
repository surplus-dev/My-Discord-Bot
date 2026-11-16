package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN 환경 변수를 설정해주세요")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discord 세션 생성 실패: %v", err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author == nil || m.Author.Bot {
			return
		}

		if m.Content == "!ping" {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
				log.Printf("메시지 전송 실패: %v", err)
			}
		}
	})

	session.AddHandler(func(_ *discordgo.Session, r *discordgo.Ready) {
		log.Printf("봇 로그인 완료: %s", r.User.Username)
	})

	if err := session.Open(); err != nil {
		log.Fatalf("Discord 연결 실패: %v", err)
	}
	defer session.Close()

	log.Println("봇 실행 중 (!ping으로 테스트)")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("봇을 종료합니다")
}
