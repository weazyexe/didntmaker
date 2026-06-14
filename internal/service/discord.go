package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"weazyexe.dev/didntmaker/internal/i18n"
	"weazyexe.dev/didntmaker/internal/repository"

	tele "gopkg.in/telebot.v3"
)

type DiscordService interface {
	Start() error
	Stop()
}

type discordService struct {
	session     *discordgo.Session
	bindingRepo repository.DiscordBindingRepository
	teleBot     *tele.Bot
	msg         *i18n.Messages
}

func NewDiscordService(
	token string,
	bindingRepo repository.DiscordBindingRepository,
	teleBot *tele.Bot,
	msg *i18n.Messages,
) (*discordService, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	sess.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	svc := &discordService{
		session:     sess,
		bindingRepo: bindingRepo,
		teleBot:     teleBot,
		msg:         msg,
	}

	sess.AddHandler(svc.handleVoiceStateUpdate)

	return svc, nil
}

func (s *discordService) Start() error {
	if err := s.session.Open(); err != nil {
		return err
	}
	slog.Info("discord session opened")
	return nil
}

func (s *discordService) Stop() {
	if err := s.session.Close(); err != nil {
		slog.Error("failed to close discord session", "error", err)
	}
	slog.Info("discord session closed")
}

func (s *discordService) handleVoiceStateUpdate(_ *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	if e.Member == nil || e.Member.User == nil {
		return
	}

	// Skip bots
	if e.Member.User.Bot {
		return
	}

	beforeChannelID := ""
	if e.BeforeUpdate != nil {
		beforeChannelID = e.BeforeUpdate.ChannelID
	}
	afterChannelID := e.ChannelID

	var eventMsg string

	switch {
	case beforeChannelID == "" && afterChannelID != "":
		// Join
		channelName := s.resolveChannelName(afterChannelID)
		displayName := s.resolveDisplayName(e.Member)
		eventMsg = fmt.Sprintf(s.msg.DiscordVoiceJoin, displayName, channelName)
	case beforeChannelID != "" && afterChannelID == "":
		// Leave
		channelName := s.resolveChannelName(beforeChannelID)
		displayName := s.resolveDisplayName(e.Member)
		eventMsg = fmt.Sprintf(s.msg.DiscordVoiceLeave, displayName, channelName)
	default:
		// Channel switch, mute/deafen, etc. — ignore
		return
	}

	bindings, err := s.bindingRepo.GetByGuildID(context.Background(), e.GuildID)
	if err != nil {
		slog.Error("failed to get discord bindings",
			"guild_id", e.GuildID,
			"error", err,
		)
		return
	}

	for _, binding := range bindings {
		chat := tele.ChatID(binding.ChatID)
		if _, err := s.teleBot.Send(chat, eventMsg); err != nil {
			slog.Error("failed to send discord voice notification",
				"chat_id", binding.ChatID,
				"guild_id", binding.GuildID,
				"error", err,
			)
		}
	}
}

func (s *discordService) resolveChannelName(channelID string) string {
	// Try cache first
	ch, err := s.session.State.Channel(channelID)
	if err == nil {
		return ch.Name
	}

	// Fallback to API
	ch, err = s.session.Channel(channelID)
	if err != nil {
		slog.Warn("failed to resolve channel name",
			"channel_id", channelID,
			"error", err,
		)
		return channelID
	}
	return ch.Name
}

func (s *discordService) resolveDisplayName(member *discordgo.Member) string {
	if member.Nick != "" {
		return member.Nick
	}
	if member.User.GlobalName != "" {
		return member.User.GlobalName
	}
	return member.User.Username
}
