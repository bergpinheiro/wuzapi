# Integração Wuzapi x Chatwoot

Esta integração permite conectar a Wuzapi com o Chatwoot para criar um sistema completo de atendimento ao cliente via WhatsApp.

## 🚀 Funcionalidades

### ✅ Implementadas
- **Envio automático de mensagens do WhatsApp para o Chatwoot**
- **Envio de mensagens do Chatwoot para o WhatsApp**
- **Suporte a todos os tipos de mídia** (imagem, vídeo, áudio, documento, sticker)
- **Criação automática de contatos e conversas**
- **Sincronização de status de mensagens**
- **Suporte a grupos do WhatsApp**
- **Interface de configuração no dashboard**
- **Validação de webhooks com assinatura**
- **Mapeamento de mensagens bidirecional**

### 📱 Tipos de Mensagem Suportados
- 📝 **Texto simples** - Mensagens de texto normais
- 🖼️ **Imagens** - Com legenda e upload automático
- 🎥 **Vídeos** - Com thumbnail e upload automático
- 🎵 **Áudios** - Notas de voz e arquivos de áudio
- 📄 **Documentos** - PDFs, DOCs e outros arquivos
- 📍 **Localização** - Coordenadas GPS com link do Google Maps
- 🎭 **Stickers** - Figurinhas do WhatsApp
- 👤 **Contatos** - Cartões de contato
- 🗑️ **Mensagens deletadas** - Notificação de exclusão
- 😀 **Reações** - Emojis de reação

## 🔧 Configuração

### 1. Configuração no Chatwoot

1. Acesse sua instância do Chatwoot
2. Vá em **Configurações > Contas**
3. Anote o **ID da Conta**
4. Vá em **Configurações > Caixas de Entrada**
5. Crie uma nova caixa de entrada ou use uma existente
6. Anote o **ID da Caixa de Entrada**
7. Gere um **Token de API** nas configurações da conta

### 2. Configuração na Wuzapi

1. Acesse o dashboard da Wuzapi
2. Clique em **Chatwoot Integration** na seção de configurações
3. Preencha os campos:
   - **URL Base do Chatwoot**: `https://app.chatwoot.com` (ou sua instância)
   - **ID da Conta**: ID da conta no Chatwoot
   - **Token da API**: Token gerado no Chatwoot
   - **ID da Caixa de Entrada**: ID da caixa de entrada
   - **Segredo do Webhook**: (opcional) Para validação de segurança
4. Clique em **Testar Conexão** para verificar
5. Clique em **Salvar Configuração**

### 3. Configuração do Webhook no Chatwoot

1. No Chatwoot, vá em **Configurações > Integrações**
2. Adicione uma nova integração webhook
3. Use a URL fornecida no dashboard da Wuzapi:
   ```
   https://sua-wuzapi.com/integrations/chatwoot/webhook
   ```
4. Configure os eventos para enviar: `message_created`
5. Salve a configuração

## 📋 API Endpoints

### Configuração
- `POST /integrations/chatwoot/config` - Salvar configuração
- `GET /integrations/chatwoot/config` - Obter configuração atual
- `POST /integrations/chatwoot/test` - Testar conexão

### Webhook
- `POST /integrations/chatwoot/webhook` - Receber webhooks do Chatwoot

## 🗄️ Estrutura do Banco de Dados

### Tabelas Criadas

#### `chatwoot_integrations`
Armazena as configurações de integração por usuário:
- `user_id` - ID do usuário da Wuzapi
- `base_url` - URL base do Chatwoot
- `account_id` - ID da conta no Chatwoot
- `api_token` - Token de API (criptografado)
- `inbox_id` - ID da caixa de entrada
- `webhook_secret` - Segredo para validação de webhooks
- `enabled` - Status da integração

#### `chatwoot_conversations`
Mapeia conversas entre WhatsApp e Chatwoot:
- `user_id` - ID do usuário da Wuzapi
- `chatwoot_conversation_id` - ID da conversa no Chatwoot
- `whatsapp_chat_jid` - JID do chat no WhatsApp
- `whatsapp_contact_jid` - JID do contato no WhatsApp

#### `chatwoot_message_mapping`
Mapeia mensagens entre os dois sistemas:
- `user_id` - ID do usuário da Wuzapi
- `whatsapp_message_id` - ID da mensagem no WhatsApp
- `chatwoot_message_id` - ID da mensagem no Chatwoot
- `direction` - Direção (inbound/outbound)

## 🔄 Fluxo de Funcionamento

### WhatsApp → Chatwoot
1. Mensagem recebida no WhatsApp
2. Sistema identifica o tipo de mensagem
3. Baixa mídia se necessário
4. Busca ou cria contato no Chatwoot
5. Busca ou cria conversa no Chatwoot
6. Envia mensagem com anexos
7. Salva mapeamento da mensagem

### Chatwoot → WhatsApp
1. Agente envia mensagem no Chatwoot
2. Webhook é recebido pela Wuzapi
3. Sistema valida assinatura do webhook
4. Busca conversa mapeada
5. Identifica tipo de mensagem
6. Envia para o WhatsApp via API
7. Salva mapeamento da mensagem

## 🛡️ Segurança

- **Validação de webhooks** com assinatura HMAC-SHA256
- **Tokens de API** armazenados de forma segura
- **Validação de origem** dos webhooks
- **Logs de auditoria** para todas as operações

## 🐛 Troubleshooting

### Problemas Comuns

#### Conexão falha
- Verifique se a URL base está correta
- Confirme se o token de API é válido
- Verifique se a conta e caixa de entrada existem

#### Mensagens não chegam no Chatwoot
- Verifique se a integração está habilitada
- Confirme se o webhook está configurado corretamente
- Verifique os logs da Wuzapi

#### Mensagens não chegam no WhatsApp
- Verifique se o webhook está configurado no Chatwoot
- Confirme se a URL do webhook está acessível
- Verifique se a instância WhatsApp está conectada

### Logs
Os logs da integração podem ser encontrados nos logs da Wuzapi com as seguintes tags:
- `chatwoot` - Operações gerais
- `chatwoot_webhook` - Processamento de webhooks
- `chatwoot_message` - Processamento de mensagens

## 📈 Monitoramento

### Métricas Importantes
- Número de mensagens processadas
- Taxa de sucesso de envio
- Tempo de resposta das APIs
- Erros de validação de webhook

### Alertas Recomendados
- Falhas de conexão com Chatwoot
- Erros de validação de webhook
- Mensagens não processadas
- Timeouts de API

## 🔧 Desenvolvimento

### Estrutura de Arquivos
```
├── chatwoot.go              # Cliente e lógica do Chatwoot
├── handlers.go              # Handlers HTTP
├── migrations.go            # Migrações do banco
├── static/dashboard/
│   └── (modais integrados no index.html)
│   └── js/app.js           # JavaScript do dashboard
└── CHATWOOT_INTEGRATION.md  # Esta documentação
```

### Adicionando Novos Tipos de Mensagem
1. Adicione o tipo em `ProcessWhatsAppMessage` em `chatwoot.go`
2. Implemente a lógica de processamento
3. Adicione suporte no webhook handler se necessário
4. Atualize a documentação

## 📞 Suporte

Para suporte técnico ou dúvidas sobre a integração:
1. Verifique esta documentação
2. Consulte os logs da aplicação
3. Abra uma issue no repositório
4. Entre em contato com a equipe de desenvolvimento

---

**Versão da Integração**: 1.0.0  
**Compatibilidade**: Wuzapi 1.0.3+  
**Última Atualização**: Dezembro 2024
