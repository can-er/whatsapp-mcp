# WhatsApp MCP Server

This is a Model Context Protocol (MCP) server for WhatsApp.

With this you can search and read your personal Whatsapp messages (including images, videos, documents, and audio messages), search your contacts and send messages to either individuals or groups. You can also send media files including images, videos, documents, and audio messages.

It connects to your **personal WhatsApp account** directly via the Whatsapp web multidevice API (using the [whatsmeow](https://github.com/tulir/whatsmeow) library). All your messages are stored locally in a SQLite database and only sent to an LLM (such as Claude) when the agent accesses them through tools (which you control).

Here's an example of what you can do when it's connected to Claude.

![WhatsApp MCP](./example-use.png)

> To get updates on this and other projects I work on [enter your email here](https://docs.google.com/forms/d/1rTF9wMBTN0vPfzWuQa2BjfGKdKIpTbyeKxhPMcEzgyI/preview)

> *Caution:* as with many MCP servers, the WhatsApp MCP is subject to [the lethal trifecta](https://simonwillison.net/2025/Jun/16/the-lethal-trifecta/). This means that project injection could lead to private data exfiltration.
> Any message you receive is untrusted input that the model reads: someone can send you a message containing instructions (prompt injection) that try to make the agent forward your private conversations or files to them. If you only need to read and search messages, enable [read-only mode](#read-only-mode).

## Installation

### Prerequisites

- Go 1.26 or newer
- Python 3.6+
- Anthropic Claude Desktop app (or Cursor)
- UV (Python package manager), install with `curl -LsSf https://astral.sh/uv/install.sh | sh`
- FFmpeg (_optional_) - Only needed for audio messages. If you want to send audio files as playable WhatsApp voice messages, they must be in `.ogg` Opus format. With FFmpeg installed, the MCP server will automatically convert non-Opus audio files. Without FFmpeg, you can still send raw audio files using the `send_file` tool.

### Steps

1. **Clone this repository**

   ```bash
   git clone https://github.com/lharries/whatsapp-mcp.git
   cd whatsapp-mcp
   ```

2. **Run the WhatsApp bridge**

   Navigate to the whatsapp-bridge directory and run the Go application:

   ```bash
   cd whatsapp-bridge
   go run main.go
   ```

   The first time you run it, you will be prompted to scan a QR code. Scan the QR code with your WhatsApp mobile app to authenticate.

   After approximately 20 days, you will might need to re-authenticate.

3. **Connect to the MCP server**

   Copy the below json with the appropriate {{PATH}} values:

   ```json
   {
     "mcpServers": {
       "whatsapp": {
         "command": "{{PATH_TO_UV}}", // Run `which uv` and place the output here
         "args": [
           "--directory",
           "{{PATH_TO_SRC}}/whatsapp-mcp/whatsapp-mcp-server", // cd into the repo, run `pwd` and enter the output here + "/whatsapp-mcp-server"
           "run",
           "main.py"
         ]
       }
     }
   }
   ```

   For **Claude**, save this as `claude_desktop_config.json` in your Claude Desktop configuration directory at:

   ```
   ~/Library/Application Support/Claude/claude_desktop_config.json
   ```

   For **Cursor**, save this as `mcp.json` in your Cursor configuration directory at:

   ```
   ~/.cursor/mcp.json
   ```

4. **Restart Claude Desktop / Cursor**

   Open Claude Desktop and you should now see WhatsApp as an available integration.

   Or restart Cursor.

### Windows Compatibility

If you're running this project on Windows, be aware that `go-sqlite3` requires **CGO to be enabled** in order to compile and work properly. By default, **CGO is disabled on Windows**, so you need to explicitly enable it and have a C compiler installed.

#### Steps to get it working:

1. **Install a C compiler**  
   We recommend using [MSYS2](https://www.msys2.org/) to install a C compiler for Windows. After installing MSYS2, make sure to add the `ucrt64\bin` folder to your `PATH`.  
   → A step-by-step guide is available [here](https://code.visualstudio.com/docs/cpp/config-mingw).

2. **Enable CGO and run the app**

   ```bash
   cd whatsapp-bridge
   go env -w CGO_ENABLED=1
   go run main.go
   ```

Without this setup, you'll likely run into errors like:

> `Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work.`

## Security

### Read-only mode

Set `WHATSAPP_MCP_READ_ONLY=1` in the MCP server's environment to remove the tools that send messages (`send_message`, `send_file`, `send_audio_message`). The model can then still search and read your messages, but cannot send anything, which removes the exfiltration path of a prompt injection hidden in a received message:

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "{{PATH_TO_UV}}",
      "args": ["--directory", "{{PATH_TO_SRC}}/whatsapp-mcp/whatsapp-mcp-server", "run", "main.py"],
      "env": { "WHATSAPP_MCP_READ_ONLY": "1" }
    }
  }
}
```

### Bridge REST API

The bridge exposes a small REST API (`/api/send`, `/api/download`) used by the MCP server. It is locked down by default:

- It listens on `127.0.0.1:8080` only. Set `WHATSAPP_BRIDGE_ADDR` to change the address; avoid exposing it beyond loopback.
- Every request needs an `Authorization: Bearer <token>` header. On first start the bridge generates a random token in `whatsapp-bridge/store/bridge_token` (owner-only permissions), and the MCP server reads it from there. To use your own token (for example when the MCP server runs from another directory), set `WHATSAPP_BRIDGE_TOKEN` to the same value for both processes.
- Requests must use `Content-Type: application/json`, bodies are limited to 1 MiB, and any request carrying an `Origin` header (i.e. sent by a web page) is rejected.
- If the bridge listens on another address, point the MCP server at it with `WHATSAPP_API_BASE_URL` (default `http://127.0.0.1:8080/api`).

### Sending files

`send_file` and `send_audio_message` can only send files located in the folder named by `WHATSAPP_MEDIA_DIR` (symlinks are resolved, so a link pointing outside the folder is refused). Set the same value for the bridge and the MCP server. When it is not set, sending files is disabled; text messages still work. Audio converted by ffmpeg is written to that folder temporarily and deleted after sending.

### Restricting chats

By default every chat (including the synced history) is stored and readable by the MCP server. To limit the bridge to specific chats, set `WHATSAPP_ALLOWED_CHATS` to a comma-separated list of JIDs or phone numbers with country code, for both the bridge and the MCP server:

```bash
export WHATSAPP_ALLOWED_CHATS="15550000003,+1 555 000 0004,120363000000000000@g.us"
```

Other chats are then not stored, are hidden from the MCP tools (even if they were stored before), and cannot be used with `/api/send` or `/api/download`.

### Downloaded media

File names of received documents are chosen by the sender, so the bridge reduces them to a single safe file name (no directories, `..`, control characters or Windows reserved names) before saving media under `whatsapp-bridge/store/<chat>/`.

## Message events (fork)

For local automation, the bridge can emit every real-time message of the allowed chats (see `WHATSAPP_ALLOWED_CHATS`; history sync is not emitted). Both outputs are optional and independent:

| Variable | Meaning |
|---|---|
| `WHATSAPP_INBOX_DIR` | Durable inbox. Media is downloaded first to `media/<message id>.<ext>` (full quality, generated name), then the event is written to `events/<YYYYMMDDTHHMMSSZ>_<message id>.json`. Both are written to `tmp/` and renamed, so watch `events/` for `IN_MOVED_TO`. |
| `WHATSAPP_EVENT_WEBHOOK_URL` | Wake-up call: the same JSON is POSTed there (3 attempts, 10 s timeout). Only `localhost` / loopback URLs are accepted. |
| `WHATSAPP_EVENT_WEBHOOK_TOKEN` | Sent as `Authorization: Bearer <token>` to the webhook. |
| `WHATSAPP_INBOX_RETENTION_DAYS` | Files older than this are deleted by the bridge (default 30, `0` disables the purge). |

Sub-directories are created with mode `0750` and files with `0640`. The inbox root is never chmod-ed, so to share it with another Unix user, create it as `chown <bridge-user>:<group> inbox && chmod 2750 inbox`: the setgid bit makes every file inherit the group. Keep the inbox outside `whatsapp-bridge/store/`.

Event format (version 1):

```json
{
  "version": 1,
  "type": "message",
  "id": "3EB0C0FFEE",
  "chat_jid": "15550000001@s.whatsapp.net",
  "chat_lid": null,
  "chat_name": "Display name (informative only)",
  "sender_jid": "15550000001@s.whatsapp.net",
  "sender_phone": "15550000001",
  "sender_lid": null,
  "sender_name": "Push name (informative only)",
  "is_from_me": false,
  "timestamp": "2026-10-29T09:30:00Z",
  "received_at": "2026-10-29T09:30:02Z",
  "text": "Message text or media caption",
  "quoted_message_id": null,
  "media": {
    "status": "ok",
    "type": "audio",
    "mime": "audio/ogg; codecs=opus",
    "path": "/srv/inbox/media/3EB0C0FFEE.ogg",
    "size": 12345,
    "sha256": "…",
    "voice": true,
    "duration_s": 7
  }
}
```

`chat_jid` uses the phone-number form when a LID chat can be resolved (`chat_lid` then holds the LID). `sender_phone` is always the real number, or `null` when only the LID is known (`sender_lid`). If the media cannot be downloaded, `media.status` is `"failed"`, `media.error` explains why and `media.path` is `null`. `media` is `null` for text messages.

## Send API hardening (fork)

- `WHATSAPP_SEND_TOKEN` (at least 32 characters, different from the bridge token): when set, `/api/send` only accepts this token, so a process that holds the bridge token can read (`/api/download`, `/api/health`) but cannot send. The MCP server uses it for its send tools when it is in its environment.
- `client_msg_id` (optional field of `/api/send`): a request repeated with the same key does not send a second message; the response is the original one with `"duplicate": true`. Responses include the WhatsApp `message_id`.
- Every send is appended to `whatsapp-bridge/store/send_audit.jsonl` (`0600`): time, recipient, `client_msg_id`, SHA-256 of the text (never the text itself), success and message ID. The file is replayed on start, so idempotency survives restarts.
- `GET /api/health` (bridge token) returns `{"connected": true, "logged_in": true, "last_event_at": "…"}`.

## Architecture Overview

This application consists of two main components:

1. **Go WhatsApp Bridge** (`whatsapp-bridge/`): A Go application that connects to WhatsApp's web API, handles authentication via QR code, and stores message history in SQLite. It serves as the bridge between WhatsApp and the MCP server.

2. **Python MCP Server** (`whatsapp-mcp-server/`): A Python server implementing the Model Context Protocol (MCP), which provides standardized tools for Claude to interact with WhatsApp data and send/receive messages.

### Data Storage

- All message history is stored in a SQLite database within the `whatsapp-bridge/store/` directory
- The database maintains tables for chats and messages
- Messages are indexed for efficient searching and retrieval
- `whatsapp-bridge/store/` holds your WhatsApp session keys (`whatsapp.db`, which give full access to your account) and your message history. The bridge restricts it to your user: the directory is set to `0700` (also when it already exists) and the databases and downloaded media to `0600`.
- On Windows these permission bits do not apply; the folder inherits the access rules of its parent. Keep the repository inside your user profile, or restrict the folder explicitly, e.g. `icacls whatsapp-bridge\store /inheritance:r /grant:r "%USERNAME%:(OI)(CI)F"`.

## Usage

Once connected, you can interact with your WhatsApp contacts through Claude, leveraging Claude's AI capabilities in your WhatsApp conversations.

### MCP Tools

Claude can access the following tools to interact with WhatsApp:

- **search_contacts**: Search for contacts by name or phone number
- **list_messages**: Retrieve messages with optional filters and context
- **list_chats**: List available chats with metadata
- **get_chat**: Get information about a specific chat
- **get_direct_chat_by_contact**: Find a direct chat with a specific contact
- **get_contact_chats**: List all chats involving a specific contact
- **get_last_interaction**: Get the most recent message with a contact
- **get_message_context**: Retrieve context around a specific message
- **send_message**: Send a WhatsApp message to a specified phone number or group JID
- **send_file**: Send a file (image, video, raw audio, document) to a specified recipient
- **send_audio_message**: Send an audio file as a WhatsApp voice message (requires the file to be an .ogg opus file or ffmpeg must be installed)
- **download_media**: Download media from a WhatsApp message and get the local file path

### Media Handling Features

The MCP server supports both sending and receiving various media types:

#### Media Sending

You can send various media types to your WhatsApp contacts:

- **Images, Videos, Documents**: Use the `send_file` tool to share any supported media type.
- **Voice Messages**: Use the `send_audio_message` tool to send audio files as playable WhatsApp voice messages.
  - For optimal compatibility, audio files should be in `.ogg` Opus format.
  - With FFmpeg installed, the system will automatically convert other audio formats (MP3, WAV, etc.) to the required format.
  - Without FFmpeg, you can still send raw audio files using the `send_file` tool, but they won't appear as playable voice messages.

#### Media Downloading

By default, just the metadata of the media is stored in the local database. The message will indicate that media was sent. To access this media you need to use the download_media tool which takes the `message_id` and `chat_jid` (which are shown when printing messages containing the meda), this downloads the media and then returns the file path which can be then opened or passed to another tool.

## Technical Details

1. Claude sends requests to the Python MCP server
2. The MCP server queries the Go bridge for WhatsApp data or directly to the SQLite database
3. The Go accesses the WhatsApp API and keeps the SQLite database up to date
4. Data flows back through the chain to Claude
5. When sending messages, the request flows from Claude through the MCP server to the Go bridge and to WhatsApp

## Troubleshooting

- If you encounter permission issues when running uv, you may need to add it to your PATH or use the full path to the executable.
- Make sure both the Go application and the Python server are running for the integration to work properly.

### Authentication Issues

- **QR Code Not Displaying**: If the QR code doesn't appear, try restarting the authentication script. If issues persist, check if your terminal supports displaying QR codes.
- **WhatsApp Already Logged In**: If your session is already active, the Go bridge will automatically reconnect without showing a QR code.
- **Device Limit Reached**: WhatsApp limits the number of linked devices. If you reach this limit, you'll need to remove an existing device from WhatsApp on your phone (Settings > Linked Devices).
- **No Messages Loading**: After initial authentication, it can take several minutes for your message history to load, especially if you have many chats.
- **WhatsApp Out of Sync**: If your WhatsApp messages get out of sync with the bridge, delete both database files (`whatsapp-bridge/store/messages.db` and `whatsapp-bridge/store/whatsapp.db`) and restart the bridge to re-authenticate.

For additional Claude Desktop integration troubleshooting, see the [MCP documentation](https://modelcontextprotocol.io/quickstart/server#claude-for-desktop-integration-issues). The documentation includes helpful tips for checking logs and resolving common issues.
