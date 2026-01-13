package chat

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Space-Software-LTDA/owncast/config"
	"github.com/Space-Software-LTDA/owncast/core/chat/events"
	"github.com/Space-Software-LTDA/owncast/core/webhooks"
	"github.com/Space-Software-LTDA/owncast/persistence/chatmessagerepository"
	"github.com/Space-Software-LTDA/owncast/persistence/configrepository"
	"github.com/Space-Software-LTDA/owncast/persistence/userrepository"
	"github.com/Space-Software-LTDA/owncast/utils"
	log "github.com/sirupsen/logrus"
)

func (s *Server) userNameChanged(eventData chatClientEvent) {
	var receivedEvent events.NameChangeEvent
	if err := json.Unmarshal(eventData.data, &receivedEvent); err != nil {
		log.Errorln("error unmarshalling to NameChangeEvent", err)
		return
	}

	configRepository := configrepository.Get()

	proposedUsername := receivedEvent.NewName

	// Check if name is on the blocklist
	blocklist := configRepository.GetForbiddenUsernameList()

	// Names have a max length
	proposedUsername = utils.MakeSafeStringOfLength(proposedUsername, config.MaxChatDisplayNameLength)

	// Check if the sanitized name is empty or just whitespace
	if strings.TrimSpace(proposedUsername) == "" {
		log.Debugln(logSanitize(eventData.client.User.DisplayName), "attempted to change name to empty or whitespace-only name")
		message := "Display name cannot be empty or contain only whitespace."
		s.sendActionToClient(eventData.client, message)

		// Resend the client's user so their username is in sync.
		eventData.client.sendConnectedClientInfo()

		return
	}

	for _, blockedName := range blocklist {
		normalizedName := strings.TrimSpace(blockedName)
		normalizedName = strings.ToLower(normalizedName)
		if strings.Contains(normalizedName, proposedUsername) {
			// Denied.
			log.Debugln(logSanitize(eventData.client.User.DisplayName), "blocked from changing name to", logSanitize(proposedUsername), "due to blocked name", normalizedName)
			message := fmt.Sprintf("You cannot change your name to **%s**.", proposedUsername)
			s.sendActionToClient(eventData.client, message)

			// Resend the client's user so their username is in sync.
			eventData.client.sendConnectedClientInfo()

			return
		}
	}

	userRepository := userrepository.Get()

	// Check if the name is not already assigned to a registered user.
	if available, err := userRepository.IsDisplayNameAvailable(proposedUsername); err != nil {
		log.Errorln("error checking if name is available", err)
		return
	} else if !available {
		message := fmt.Sprintf("The name **%s** has already been registered. If this is your name, please authenticate.", proposedUsername)
		s.sendActionToClient(eventData.client, message)

		// Resend the client's user so their username is in sync.
		eventData.client.sendConnectedClientInfo()

		return
	}

	savedUser := userRepository.GetUserByToken(eventData.client.accessToken)
	oldName := savedUser.DisplayName

	// Check that the new name is different from old.
	if proposedUsername == oldName {
		eventData.client.sendConnectedClientInfo()
		return
	}

	// Save the new name
	if err := userRepository.ChangeUsername(eventData.client.User.ID, proposedUsername); err != nil {
		log.Errorln("error changing username", err)
	}

	// Update the connected clients associated user with the new name
	now := time.Now()
	eventData.client.User = savedUser
	eventData.client.User.NameChangedAt = &now

	// Send chat event letting everyone about about the name change
	savedUser.DisplayName = proposedUsername

	broadcastEvent := events.NameChangeBroadcast{
		Oldname: oldName,
	}
	broadcastEvent.User = savedUser
	broadcastEvent.SetDefaults()
	payload := broadcastEvent.GetBroadcastPayload()
	if err := s.Broadcast(payload); err != nil {
		log.Errorln("error broadcasting NameChangeEvent", err)
		return
	}

	// Send chat user name changed webhook
	receivedEvent.User = savedUser
	receivedEvent.ClientID = eventData.client.Id
	webhooks.SendChatEventUsernameChanged(receivedEvent)

	// Resend the client's user so their username is in sync.
	eventData.client.sendConnectedClientInfo()
}

func (s *Server) userColorChanged(eventData chatClientEvent) {
	userRepository := userrepository.Get()

	var receivedEvent events.ColorChangeEvent
	if err := json.Unmarshal(eventData.data, &receivedEvent); err != nil {
		log.Errorln("error unmarshalling to ColorChangeEvent", err)
		return
	}

	// Verify this color is valid
	if receivedEvent.NewColor > config.MaxUserColor {
		log.Errorln("invalid color requested when changing user display color")
		return
	}

	// Save the new color
	if err := userRepository.ChangeUserColor(eventData.client.User.ID, receivedEvent.NewColor); err != nil {
		log.Errorln("error changing user display color", err)
	}

	// Resend client's user info with new color, otherwise the name change dialog would still show the old color
	eventData.client.User.DisplayColor = receivedEvent.NewColor
	eventData.client.sendConnectedClientInfo()
}

func (s *Server) userMessageSent(eventData chatClientEvent) {
	userRepository := userrepository.Get()
	configRepository := configrepository.Get()

	if configRepository.GetChatDisabled() {
		log.Debugln("chat is disabled, ignoring user message")
		return
	}

	var event events.UserMessageEvent
	if err := json.Unmarshal(eventData.data, &event); err != nil {
		log.Errorln("error unmarshalling to UserMessageEvent", err)
		return
	}

	event.SetDefaults()
	event.ClientID = eventData.client.Id

	// Ignore empty messages
	if event.Empty() {
		return
	}

	// Ignore if the stream has been offline
	if !getStatus().Online && getStatus().LastDisconnectTime != nil {
		disconnectedTime := getStatus().LastDisconnectTime.Time
		if time.Since(disconnectedTime) > 5*time.Minute {
			return
		}
	}

	event.User = userRepository.GetUserByToken(eventData.client.accessToken)

	// Guard against nil users
	if event.User == nil {
		return
	}

	prohibitedTerms := []string{
		"golpe",
		"golpes",
		"golpista",
		"golpistas",
		"golp",
		"g0lpe",
		"g0lpista",
		"scam",
		"scammer",
		"fraude",
		"fraude financeira",
		"esquema",
		"esquema ilegal",
		"piramide",
		"esquema de piramide",
		"pirâmide",
		"esquema de pirâmide",
		"estelionato",
		"crime",
		"criminoso",
		"quadrilha",
		"safado",
		"safados",
		"ladrão",
		"ladrao",
		"ladrões",
		"ladroes",
		"vagabundo",
		"vagabundos",
		"pilantra",
		"pilantragem",
		"canalha",
		"mau caráter",
		"mau-caráter",
		"bandido",
		"bandidagem",
		"charlatão",
		"charlatao",
		"charlatões",
		"picareta",
		"picaretagem",
		"sacar",
		"saque",
		"não consigo sacar",
		"nao consigo sacar",
		"nao consigo fazer saque",
		"saque travado",
		"saque bloqueado",
		"saque pendente",
		"dinheiro preso",
		"dinheiro bloqueado",
		"dinheiro sumiu",
		"cadê meu dinheiro",
		"cade meu dinheiro",
		"cadê o dinheiro",
		"cade o dinheiro",
		"perdi dinheiro",
		"perdi meu dinheiro",
		"não recebi",
		"nao recebi",
		"não caiu",
		"nao caiu",
		"saldo travado",
		"saldo bloqueado",
		"estou com problemas",
		"estou com problema",
		"problema",
		"problemas",
		"isso é problema",
		"isso e problema",
		"não funciona",
		"nao funciona",
		"não está funcionando",
		"nao esta funcionando",
		"não resolve",
		"nao resolve",
		"ninguém responde",
		"ninguem responde",
		"suporte não responde",
		"suporte nao responde",
		"ninguém ajuda",
		"ninguem ajuda",
		"mentira",
		"mentiroso",
		"mentirosos",
		"enganação",
		"enganacao",
		"enganar",
		"enganando",
		"ilusão",
		"ilusao",
		"ilusão financeira",
		"ilusao financeira",
		"promessa falsa",
		"falsas promessas",
		"propaganda enganosa",
		"furada",
		"fria",
		"isso é fria",
		"isso e fria",
		"perda de tempo",
		"não recomendo",
		"nao recomendo",
		"cuidado",
		"tomem cuidado",
		"alerta",
		"aviso",
		"denúncia",
		"denuncia",
		"vou denunciar",
		"denunciar",
		"reclame aqui",
		"procon",
		"polícia",
		"policia",
		"isso é golpe",
		"isso e golpe",
		"é golpe",
		"e golpe",
		"isso é fraude",
		"isso e fraude",
		"vocês são golpistas",
		"voces sao golpistas",
		"empresa golpista",
		"site golpista",
		"roubo",
		"roubaram meu dinheiro",
		"perdi tudo",
		"fui enganado",
		"fui enganada",
	}

	var containsProhibitedTerm bool = false
	for _, term := range prohibitedTerms {
		term = strings.ToLower(term)

		if strings.Contains(strings.ToLower(event.Body), term) {
			containsProhibitedTerm = true
		}
	}

	if !containsProhibitedTerm {
		payload := event.GetBroadcastPayload()
		if err := s.Broadcast(payload); err != nil {
			log.Errorln("error broadcasting UserMessageEvent payload", err)
			return
		}
	} else {
		hidden := time.Now()
		event.HiddenAt = &hidden
	}

	// Send chat message sent webhook
	webhooks.SendChatEvent(&event)
	chatMessagesSentCounter.Inc()
	chatMessageRepository := chatmessagerepository.Get()
	chatMessageRepository.SaveUserMessage(event)
	eventData.client.MessageCount++
}

func logSanitize(userValue string) string {
	// strip carriage return and newline from user-submitted values to prevent log injection
	sanitizedValue := strings.ReplaceAll(userValue, "\n", "")
	sanitizedValue = strings.ReplaceAll(sanitizedValue, "\r", "")

	return fmt.Sprintf("userSuppliedValue(%s)", sanitizedValue)
}
