package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// ChatwootIntegration representa a configuração de integração do Chatwoot
type ChatwootIntegration struct {
	ID            int       `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	BaseURL       string    `json:"base_url" db:"base_url"`
	AccountID     string    `json:"account_id" db:"account_id"`
	APIToken      string    `json:"api_token" db:"api_token"`
	InboxID       string    `json:"inbox_id" db:"inbox_id"`
	WebhookSecret string    `json:"webhook_secret" db:"webhook_secret"`
	Enabled       bool      `json:"enabled" db:"enabled"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ChatwootConversation representa o mapeamento entre conversas do WhatsApp e Chatwoot
type ChatwootConversation struct {
	ID                      int       `json:"id" db:"id"`
	UserID                  string    `json:"user_id" db:"user_id"`
	ChatwootConversationID  int       `json:"chatwoot_conversation_id" db:"chatwoot_conversation_id"`
	WhatsAppChatJID         string    `json:"whatsapp_chat_jid" db:"whatsapp_chat_jid"`
	WhatsAppContactJID      string    `json:"whatsapp_contact_jid" db:"whatsapp_contact_jid"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

// ChatwootMessageMapping representa o mapeamento entre mensagens do WhatsApp e Chatwoot
type ChatwootMessageMapping struct {
	ID                  int       `json:"id" db:"id"`
	UserID              string    `json:"user_id" db:"user_id"`
	WhatsAppMessageID   string    `json:"whatsapp_message_id" db:"whatsapp_message_id"`
	ChatwootMessageID   int       `json:"chatwoot_message_id" db:"chatwoot_message_id"`
	Direction           string    `json:"direction" db:"direction"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// Estruturas para API do Chatwoot
type ChatwootContact struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Identifier  string `json:"identifier"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ChatwootConversationAPI struct {
	ID          int    `json:"id"`
	AccountID   int    `json:"account_id"`
	InboxID     int    `json:"inbox_id"`
	Status      string `json:"status"`
	Contact     *ChatwootContact `json:"contact"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ChatwootMessage struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	MessageType  int    `json:"message_type"`
	ContentType  string `json:"content_type"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Sender       *ChatwootSender `json:"sender"`
	Attachments  []ChatwootAttachment `json:"attachments"`
}

type ChatwootSender struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

type ChatwootAttachment struct {
	ID           int    `json:"id"`
	FileType     string `json:"file_type"`
	FileURL      string `json:"file_url"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	ContentType  string `json:"content_type"`
}

type ChatwootWebhookPayload struct {
	Account     ChatwootAccount     `json:"account"`
	Conversation ChatwootConversationAPI `json:"conversation"`
	Message     ChatwootMessage     `json:"message"`
	Inbox       ChatwootInbox       `json:"inbox"`
	Event       string              `json:"event"`
}

type ChatwootAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ChatwootInbox struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ChatwootClient representa o cliente para comunicação com a API do Chatwoot
type ChatwootClient struct {
	BaseURL   string
	APIToken  string
	AccountID string
	InboxID   string
	Client    *resty.Client
}

// NewChatwootClient cria uma nova instância do cliente Chatwoot
func NewChatwootClient(baseURL, apiToken, accountID, inboxID string) *ChatwootClient {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeader("Authorization", "Bearer "+apiToken)
	client.SetHeader("Content-Type", "application/json")

	return &ChatwootClient{
		BaseURL:   strings.TrimSuffix(baseURL, "/"),
		APIToken:  apiToken,
		AccountID: accountID,
		InboxID:   inboxID,
		Client:    client,
	}
}

// TestConnection testa a conexão com a API do Chatwoot
func (c *ChatwootClient) TestConnection() error {
	resp, err := c.Client.R().
		Get(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/inboxes/" + c.InboxID)
	
	if err != nil {
		return fmt.Errorf("erro ao testar conexão: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	return nil
}

// GetContact busca um contato pelo número de telefone
func (c *ChatwootClient) GetContact(phoneNumber string) (*ChatwootContact, error) {
	resp, err := c.Client.R().
		SetQueryParam("inbox_id", c.InboxID).
		SetQueryParam("phone_number", phoneNumber).
		Get(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/contacts")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contato: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload []ChatwootContact `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	if len(result.Payload) > 0 {
		return &result.Payload[0], nil
	}
	
	return nil, nil
}

// CreateContact cria um novo contato
func (c *ChatwootClient) CreateContact(phoneNumber, name string) (*ChatwootContact, error) {
	payload := map[string]interface{}{
		"name":         name,
		"phone_number": phoneNumber,
		"identifier":   phoneNumber,
	}
	
	resp, err := c.Client.R().
		SetBody(payload).
		Post(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/contacts")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao criar contato: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload ChatwootContact `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	return &result.Payload, nil
}

// GetConversation busca uma conversa pelo ID do contato
func (c *ChatwootClient) GetConversation(contactID int) (*ChatwootConversationAPI, error) {
	resp, err := c.Client.R().
		Get(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/contacts/" + strconv.Itoa(contactID) + "/conversations")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conversa: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload []ChatwootConversationAPI `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	if len(result.Payload) > 0 {
		return &result.Payload[0], nil
	}
	
	return nil, nil
}

// CreateConversation cria uma nova conversa
func (c *ChatwootClient) CreateConversation(contactID int) (*ChatwootConversationAPI, error) {
	payload := map[string]interface{}{
		"source_id": contactID,
		"inbox_id":  c.InboxID,
		"contact_id": contactID,
	}
	
	resp, err := c.Client.R().
		SetBody(payload).
		Post(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/conversations")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao criar conversa: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload ChatwootConversationAPI `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	return &result.Payload, nil
}

// SendMessage envia uma mensagem para uma conversa
func (c *ChatwootClient) SendMessage(conversationID int, content, messageType string, attachments []ChatwootAttachment) (*ChatwootMessage, error) {
	payload := map[string]interface{}{
		"content": content,
		"message_type": messageType,
	}
	
	if len(attachments) > 0 {
		payload["attachments"] = attachments
	}
	
	resp, err := c.Client.R().
		SetBody(payload).
		Post(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/conversations/" + strconv.Itoa(conversationID) + "/messages")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar mensagem: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload ChatwootMessage `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	return &result.Payload, nil
}

// UploadAttachment faz upload de um anexo
func (c *ChatwootClient) UploadAttachment(fileData []byte, fileName, contentType string) (*ChatwootAttachment, error) {
	resp, err := c.Client.R().
		SetFileReader("file", fileName, strings.NewReader(string(fileData))).
		SetFormData(map[string]string{
			"file_type": contentType,
		}).
		Post(c.BaseURL + "/api/v1/accounts/" + c.AccountID + "/conversations/attachments")
	
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do anexo: %w", err)
	}
	
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("erro na API do Chatwoot: %s", resp.String())
	}
	
	var result struct {
		Payload ChatwootAttachment `json:"payload"`
	}
	
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	return &result.Payload, nil
}

// ValidateWebhookSignature valida a assinatura do webhook do Chatwoot
func ValidateChatwootWebhookSignature(payload []byte, signature, secret string) bool {
	if secret == "" {
		return true // Se não há secret configurado, aceita
	}
	
	expectedSignature := "sha256=" + hex.EncodeToString(hmac.New(sha256.New, []byte(secret)).Sum(payload))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Métodos de banco de dados para integração Chatwoot

// GetChatwootIntegration busca a configuração de integração do Chatwoot para um usuário
func (s *server) GetChatwootIntegration(userID string) (*ChatwootIntegration, error) {
	var integration ChatwootIntegration
	query := `SELECT id, user_id, base_url, account_id, api_token, inbox_id, webhook_secret, enabled, created_at, updated_at 
			  FROM chatwoot_integrations WHERE user_id = $1`
	
	if s.db.DriverName() == "sqlite" {
		query = `SELECT id, user_id, base_url, account_id, api_token, inbox_id, webhook_secret, enabled, created_at, updated_at 
				 FROM chatwoot_integrations WHERE user_id = ?`
	}
	
	err := s.db.Get(&integration, query, userID)
	if err != nil {
		return nil, err
	}
	
	return &integration, nil
}

// SaveChatwootIntegration salva ou atualiza a configuração de integração do Chatwoot
func (s *server) SaveChatwootIntegration(integration *ChatwootIntegration) error {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `INSERT INTO chatwoot_integrations (user_id, base_url, account_id, api_token, inbox_id, webhook_secret, enabled, updated_at)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				 ON CONFLICT (user_id) DO UPDATE SET
				 base_url = EXCLUDED.base_url,
				 account_id = EXCLUDED.account_id,
				 api_token = EXCLUDED.api_token,
				 inbox_id = EXCLUDED.inbox_id,
				 webhook_secret = EXCLUDED.webhook_secret,
				 enabled = EXCLUDED.enabled,
				 updated_at = EXCLUDED.updated_at`
	} else {
		query = `INSERT OR REPLACE INTO chatwoot_integrations (user_id, base_url, account_id, api_token, inbox_id, webhook_secret, enabled, updated_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	}
	
	integration.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query, integration.UserID, integration.BaseURL, integration.AccountID, 
		integration.APIToken, integration.InboxID, integration.WebhookSecret, integration.Enabled, integration.UpdatedAt)
	
	return err
}

// GetChatwootConversation busca uma conversa mapeada
func (s *server) GetChatwootConversation(userID, whatsappChatJID string) (*ChatwootConversation, error) {
	var conversation ChatwootConversation
	query := `SELECT id, user_id, chatwoot_conversation_id, whatsapp_chat_jid, whatsapp_contact_jid, created_at, updated_at
			  FROM chatwoot_conversations WHERE user_id = $1 AND whatsapp_chat_jid = $2`
	
	if s.db.DriverName() == "sqlite" {
		query = `SELECT id, user_id, chatwoot_conversation_id, whatsapp_chat_jid, whatsapp_contact_jid, created_at, updated_at
				 FROM chatwoot_conversations WHERE user_id = ? AND whatsapp_chat_jid = ?`
	}
	
	err := s.db.Get(&conversation, query, userID, whatsappChatJID)
	if err != nil {
		return nil, err
	}
	
	return &conversation, nil
}

// SaveChatwootConversation salva um mapeamento de conversa
func (s *server) SaveChatwootConversation(conversation *ChatwootConversation) error {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `INSERT INTO chatwoot_conversations (user_id, chatwoot_conversation_id, whatsapp_chat_jid, whatsapp_contact_jid, updated_at)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (user_id, whatsapp_chat_jid) DO UPDATE SET
				 chatwoot_conversation_id = EXCLUDED.chatwoot_conversation_id,
				 whatsapp_contact_jid = EXCLUDED.whatsapp_contact_jid,
				 updated_at = EXCLUDED.updated_at`
	} else {
		query = `INSERT OR REPLACE INTO chatwoot_conversations (user_id, chatwoot_conversation_id, whatsapp_chat_jid, whatsapp_contact_jid, updated_at)
				 VALUES (?, ?, ?, ?, ?)`
	}
	
	conversation.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query, conversation.UserID, conversation.ChatwootConversationID, 
		conversation.WhatsAppChatJID, conversation.WhatsAppContactJID, conversation.UpdatedAt)
	
	return err
}

// SaveChatwootMessageMapping salva um mapeamento de mensagem
func (s *server) SaveChatwootMessageMapping(mapping *ChatwootMessageMapping) error {
	var query string
	if s.db.DriverName() == "postgres" {
		query = `INSERT INTO chatwoot_message_mapping (user_id, whatsapp_message_id, chatwoot_message_id, direction)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (user_id, whatsapp_message_id) DO UPDATE SET
				 chatwoot_message_id = EXCLUDED.chatwoot_message_id,
				 direction = EXCLUDED.direction`
	} else {
		query = `INSERT OR REPLACE INTO chatwoot_message_mapping (user_id, whatsapp_message_id, chatwoot_message_id, direction)
				 VALUES (?, ?, ?, ?)`
	}
	
	_, err := s.db.Exec(query, mapping.UserID, mapping.WhatsAppMessageID, 
		mapping.ChatwootMessageID, mapping.Direction)
	
	return err
}

// GetChatwootMessageMapping busca um mapeamento de mensagem
func (s *server) GetChatwootMessageMapping(userID, whatsappMessageID string) (*ChatwootMessageMapping, error) {
	var mapping ChatwootMessageMapping
	query := `SELECT id, user_id, whatsapp_message_id, chatwoot_message_id, direction, created_at
			  FROM chatwoot_message_mapping WHERE user_id = $1 AND whatsapp_message_id = $2`
	
	if s.db.DriverName() == "sqlite" {
		query = `SELECT id, user_id, whatsapp_message_id, chatwoot_message_id, direction, created_at
				 FROM chatwoot_message_mapping WHERE user_id = ? AND whatsapp_message_id = ?`
	}
	
	err := s.db.Get(&mapping, query, userID, whatsappMessageID)
	if err != nil {
		return nil, err
	}
	
	return &mapping, nil
}

// ProcessWhatsAppMessage processa uma mensagem do WhatsApp e envia para o Chatwoot
func (s *server) ProcessWhatsAppMessage(userID, chatJID, senderJID, messageID, messageType, textContent, mediaLink string) error {
	// Verificar se a integração está habilitada
	integration, err := s.GetChatwootIntegration(userID)
	if err != nil {
		log.Debug().Err(err).Str("user_id", userID).Msg("Integração Chatwoot não configurada")
		return nil // Não é um erro, apenas não processa
	}
	
	if !integration.Enabled {
		log.Debug().Str("user_id", userID).Msg("Integração Chatwoot desabilitada")
		return nil
	}
	
	// Criar cliente Chatwoot
	client := NewChatwootClient(integration.BaseURL, integration.APIToken, integration.AccountID, integration.InboxID)
	
	// Extrair número de telefone do JID
	phoneNumber := extractPhoneFromJID(senderJID)
	if phoneNumber == "" {
		log.Warn().Str("sender_jid", senderJID).Msg("Não foi possível extrair número de telefone")
		return nil
	}
	
	// Buscar ou criar contato
	contact, err := client.GetContact(phoneNumber)
	if err != nil {
		log.Error().Err(err).Str("phone_number", phoneNumber).Msg("Erro ao buscar contato no Chatwoot")
		return err
	}
	
	if contact == nil {
		// Criar novo contato
		contact, err = client.CreateContact(phoneNumber, extractNameFromJID(senderJID))
		if err != nil {
			log.Error().Err(err).Str("phone_number", phoneNumber).Msg("Erro ao criar contato no Chatwoot")
			return err
		}
		log.Info().Int("contact_id", contact.ID).Str("phone_number", phoneNumber).Msg("Contato criado no Chatwoot")
	}
	
	// Buscar ou criar conversa
	conversation, err := client.GetConversation(contact.ID)
	if err != nil {
		log.Error().Err(err).Int("contact_id", contact.ID).Msg("Erro ao buscar conversa no Chatwoot")
		return err
	}
	
	if conversation == nil {
		// Criar nova conversa
		conversation, err = client.CreateConversation(contact.ID)
		if err != nil {
			log.Error().Err(err).Int("contact_id", contact.ID).Msg("Erro ao criar conversa no Chatwoot")
			return err
		}
		log.Info().Int("conversation_id", conversation.ID).Int("contact_id", contact.ID).Msg("Conversa criada no Chatwoot")
	}
	
	// Salvar mapeamento da conversa
	chatwootConv := &ChatwootConversation{
		UserID:                 userID,
		ChatwootConversationID: conversation.ID,
		WhatsAppChatJID:        chatJID,
		WhatsAppContactJID:     senderJID,
	}
	
	if err := s.SaveChatwootConversation(chatwootConv); err != nil {
		log.Error().Err(err).Msg("Erro ao salvar mapeamento de conversa")
	}
	
	// Processar mensagem baseada no tipo
	var chatwootMessage *ChatwootMessage
	var attachments []ChatwootAttachment
	
	switch messageType {
	case "text":
		chatwootMessage, err = client.SendMessage(conversation.ID, textContent, "incoming", nil)
	case "image":
		// Upload da imagem
		if mediaLink != "" {
			imageData, err := s.downloadMedia(mediaLink)
			if err != nil {
				log.Error().Err(err).Msg("Erro ao baixar imagem")
				return err
			}
			
			attachment, err := client.UploadAttachment(imageData, "image.jpg", "image/jpeg")
			if err != nil {
				log.Error().Err(err).Msg("Erro ao fazer upload da imagem")
				return err
			}
			attachments = append(attachments, *attachment)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, textContent, "incoming", attachments)
	case "video":
		if mediaLink != "" {
			videoData, err := s.downloadMedia(mediaLink)
			if err != nil {
				log.Error().Err(err).Msg("Erro ao baixar vídeo")
				return err
			}
			
			attachment, err := client.UploadAttachment(videoData, "video.mp4", "video/mp4")
			if err != nil {
				log.Error().Err(err).Msg("Erro ao fazer upload do vídeo")
				return err
			}
			attachments = append(attachments, *attachment)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, textContent, "incoming", attachments)
	case "audio":
		if mediaLink != "" {
			audioData, err := s.downloadMedia(mediaLink)
			if err != nil {
				log.Error().Err(err).Msg("Erro ao baixar áudio")
				return err
			}
			
			attachment, err := client.UploadAttachment(audioData, "audio.ogg", "audio/ogg")
			if err != nil {
				log.Error().Err(err).Msg("Erro ao fazer upload do áudio")
				return err
			}
			attachments = append(attachments, *attachment)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, textContent, "incoming", attachments)
	case "document":
		if mediaLink != "" {
			docData, err := s.downloadMedia(mediaLink)
			if err != nil {
				log.Error().Err(err).Msg("Erro ao baixar documento")
				return err
			}
			
			attachment, err := client.UploadAttachment(docData, "document.pdf", "application/pdf")
			if err != nil {
				log.Error().Err(err).Msg("Erro ao fazer upload do documento")
				return err
			}
			attachments = append(attachments, *attachment)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, textContent, "incoming", attachments)
	case "location":
		// Para localização, enviar como texto com link do Google Maps
		locationText := "📍 Localização enviada"
		if textContent != "" {
			locationText += ": " + textContent
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, locationText, "incoming", nil)
	case "sticker":
		// Para sticker, enviar como imagem
		if mediaLink != "" {
			stickerData, err := s.downloadMedia(mediaLink)
			if err != nil {
				log.Error().Err(err).Msg("Erro ao baixar sticker")
				return err
			}
			
			attachment, err := client.UploadAttachment(stickerData, "sticker.webp", "image/webp")
			if err != nil {
				log.Error().Err(err).Msg("Erro ao fazer upload do sticker")
				return err
			}
			attachments = append(attachments, *attachment)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, "🎭 Sticker", "incoming", attachments)
	default:
		// Para outros tipos, enviar como texto
		content := textContent
		if content == "" {
			content = fmt.Sprintf("Mensagem do tipo: %s", messageType)
		}
		chatwootMessage, err = client.SendMessage(conversation.ID, content, "incoming", nil)
	}
	
	if err != nil {
		log.Error().Err(err).Str("message_type", messageType).Msg("Erro ao enviar mensagem para Chatwoot")
		return err
	}
	
	// Salvar mapeamento da mensagem
	mapping := &ChatwootMessageMapping{
		UserID:            userID,
		WhatsAppMessageID: messageID,
		ChatwootMessageID: chatwootMessage.ID,
		Direction:         "inbound",
	}
	
	if err := s.SaveChatwootMessageMapping(mapping); err != nil {
		log.Error().Err(err).Msg("Erro ao salvar mapeamento de mensagem")
	}
	
	log.Info().
		Str("user_id", userID).
		Str("whatsapp_message_id", messageID).
		Int("chatwoot_message_id", chatwootMessage.ID).
		Str("message_type", messageType).
		Msg("Mensagem enviada para Chatwoot com sucesso")
	
	return nil
}

// downloadMedia baixa uma mídia de uma URL
func (s *server) downloadMedia(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar mídia: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro HTTP ao baixar mídia: %d", resp.StatusCode)
	}
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler dados da mídia: %w", err)
	}
	
	return data, nil
}

// extractPhoneFromJID extrai o número de telefone de um JID do WhatsApp
func extractPhoneFromJID(jid string) string {
	// JID do WhatsApp tem formato: 5511999999999@s.whatsapp.net
	// ou 5511999999999:5511888888888@g.us para grupos
	parts := strings.Split(jid, "@")
	if len(parts) == 0 {
		return ""
	}
	
	// Para grupos, pegar apenas o primeiro número
	phoneParts := strings.Split(parts[0], ":")
	return phoneParts[0]
}

// extractNameFromJID extrai um nome básico de um JID
func extractNameFromJID(jid string) string {
	phone := extractPhoneFromJID(jid)
	if phone == "" {
		return "Usuário WhatsApp"
	}
	return "WhatsApp " + phone
}
