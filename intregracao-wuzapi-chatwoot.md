# 🧠 Prompt para o Cursor — Integração Nativa Wuzapi x Chatwoot

## Objetivo
Criar uma **integração nativa dentro do dashboard da Wuzapi API**, que permita **enviar e receber mensagens** entre **Wuzapi ↔ Chatwoot**, de forma bidirecional e em tempo real.

A integração deve ser modular, segura e com suporte completo a todos os tipos de mensagens suportadas pelo WhatsApp via Wuzapi.

---

## 💡 Contexto Geral

A **Wuzapi** é uma API de comunicação baseada no WhatsApp, que recebe e envia eventos via **webhooks** e **endpoints REST**.

O **Chatwoot** é uma plataforma de atendimento omnichannel open source, que possui APIs para **mensagens, contatos, contas e conversas**.

A proposta é desenvolver uma **integração nativa**, diretamente **no painel da Wuzapi**, que permita:
- Conectar uma conta Chatwoot a uma instância Wuzapi.
- Sincronizar automaticamente as mensagens e status entre os dois sistemas.
- Renderizar corretamente todos os tipos de mensagens recebidas do WhatsApp dentro do Chatwoot.
- Manter a consistência de status (entregue, lido, falha, etc).

---

## ⚙️ Requisitos Técnicos

### 1. Estrutura da Integração
- Criar um módulo dedicado dentro do painel da Wuzapi:  
  `Integrações > Chatwoot`
- Permitir configuração por instância (multiusuário).
- Campos configuráveis:
  - `CHATWOOT_BASE_URL`
  - `CHATWOOT_ACCOUNT_ID`
  - `CHATWOOT_API_TOKEN`
  - `CHATWOOT_INBOX_ID`
  - `CHATWOOT_WEBOOK_SECRET` (opcional, para validação)
  - Botão “Testar Conexão”

---

## 🔄 Fluxo de Comunicação

### 🟢 1. Wuzapi → Chatwoot
Quando a Wuzapi receber mensagens do WhatsApp, ela deve:
1. Receber o **webhook** de evento (`message`, `status`, `event`, etc).
2. Normalizar o payload no formato esperado pela API do Chatwoot:
   - Endpoint Chatwoot:  
     `POST /api/v1/accounts/:account_id/conversations/:conversation_id/messages`
3. Identificar se o remetente já existe:
   - Buscar contato via Chatwoot (`GET /contacts?inbox_id=...&phone_number=...`)
   - Se não existir, criar um novo (`POST /contacts`)
4. Criar/recuperar a conversa (`POST /conversations`).
5. Enviar mensagem com base no tipo detectado (ver abaixo).

#### Tipos de Mensagens a Tratar
| Tipo | Descrição | Renderização no Chatwoot |
|------|------------|---------------------------|
| `text` | Mensagem de texto simples | Texto padrão |
| `image` | Imagem com legenda | Tipo `image` com `attachment_url` |
| `audio` | Áudio ou nota de voz | Tipo `audio` com player embutido |
| `video` | Vídeo com thumbnail | Tipo `video` |
| `document` | Arquivo PDF, DOC, etc | Tipo `file` |
| `location` | Localização enviada | Usar mensagem com link do Google Maps |
| `sticker` | Figurinha | Exibir imagem como sticker |
| `reaction` | Reação a mensagem | Adicionar emoji inline |
| `poll` | Enquete | Exibir texto e opções votadas |
| `product` | Produto do catálogo | Mostrar imagem, nome, preço e link |
| `ad_message` | Mensagem de anúncio | Exibir como card com CTA |
| `payment` / `pix` | Mensagem de pagamento | Exibir status e valor |
| `event` | Mensagem de evento do WhatsApp | Log interno / card informativo |
| `group` | Mensagens em grupo | Exibir nome do grupo e participante |

#### Status de Mensagem
- Integrar confirmação de:
  - Enviado (`sent`)
  - Entregue (`delivered`)
  - Lido (`read`)
  - Falha (`failed`)
- Atualizar status no Chatwoot via:
  - `PATCH /api/v1/accounts/:account_id/conversations/:conversation_id/messages/:message_id`

---

### 🔵 2. Chatwoot → Wuzapi
Quando um agente enviar mensagem no Chatwoot:
1. Interceptar o **webhook de saída** do Chatwoot.
   - Endpoint configurado: `/integrations/chatwoot/webhook`
2. Identificar o tipo de mensagem (texto, imagem, arquivo, etc).
3. Reencaminhar para o endpoint da Wuzapi:
   - `POST /sessions/:session/messages`
4. Wuzapi enviará ao WhatsApp do cliente final.

#### Suporte Bidirecional:
- Texto (mensagem padrão)
- Imagens e vídeos (enviados por upload)
- Arquivos (documentos, PDFs)
- Mensagens de produto (template de catálogo)
- Áudios gravados
- Mensagens de localização
- Respostas a mensagens (reply)
- Mensagens de grupo

---

## 🔐 Segurança e Autenticação
- Todos os webhooks devem validar assinatura (`X-Hub-Signature` ou `X-Chatwoot-Signature`).
- Usar HTTPS obrigatório.
- Tokens de API armazenados de forma criptografada (ex: AES-256).
- Logs de eventos e falhas com auditoria.

---


---

## 🧠 Extras (Diferenciais)
- Exibir log visual de eventos (Chatwoot ↔ Wuzapi) no painel da Wuzapi.
- Botão “Sincronizar Histórico” (últimas 50 mensagens).
- Opção de **reenvio automático** de mensagens falhas.
- Detecção automática do tipo de mídia.
- Suporte para **grupos WhatsApp**, exibindo nome e avatar do grupo.

---

## 📚 Documentação de Referência
- **Chatwoot API** → https://www.chatwoot.com/developers/api/
- **Wuzapi API** → (usar documentação interna / Swagger)
- **WhatsApp Business Message Types** → https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages

---

## ✅ Objetivo Final
Ao concluir, a integração deve permitir:
- Enviar e receber mensagens **em tempo real** entre Chatwoot e Wuzapi.
- Renderizar corretamente todos os tipos de mensagens WhatsApp no Chatwoot.
- Sincronizar status e leitura.
- Permitir configuração simples e visual no painel da Wuzapi.

---

**Instrução final para o Cursor:**
> Gere o código completo dessa integração seguindo a arquitetura acima.  
> Priorize clareza, modularidade e segurança.  
> Use TypeScript + Node.js (Express ou NestJS).  
> Mantenha o estilo e design consistentes com o dashboard da Wuzapi.
