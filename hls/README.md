# Simple File Server

## Preparing
### Installing FFmpeg
```bash
# For mac
brew install ffmpeg

# For linux
apt install ffmpeg
```
### Convert MP3 and MP4
```bash
# Prepare mp3
ffmpeg -i Shahmen-Mark.mp3 -c:a libmp3lame -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist.m3u8 -segment_format mpegts output%03d.ts

# Prepare mp4
ffmpeg -i Shahmen-Mark.mp4 -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls outputlist.m3u8
```

## Run apps
```bash
# Run server
go run hls/main.go
```

## Resources
 - https://www.cdnvideo.ru/blog/protokoly_dlya_striminga/
 - https://habr.com/ru/companies/odnoklassniki/articles/467669/