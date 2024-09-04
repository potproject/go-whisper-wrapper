# go-whisper-wrapper

Whisper(OpenAI / Create transcription)を使用して、音声ファイルをテキストに変換する
OpenAI Create transcription APIまたはWhisper(OSS)を使用可能

## Setting

Create .env file

```bash
CHATGPT_WHISPER_API=https://api.openai.com/v1/audio/transcriptions
CHATGPT_API_KEY=YOUR_API_KEY
TRANSCRIPTION_LANGUAGE=ja
EXPORT_EXTENTION=srt
SERVICE=openai
```

## Usage

```bash
go run main.go input.mp3
```

## License
MIT