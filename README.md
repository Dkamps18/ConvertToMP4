# ConvertToMP4
This is a small little tool to convert all video files in a given folder to MP4

## Requirements
- FFmpeg

## Usage
| Flag      | Description                                                                                                                            | Default                   |
|:--------- | -------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- |
| \-p       | Path to the folder from where to convert from                                                                                          | Current Working Directory |
| \-v       | Verbose output logging                                                                                                                 | false                     |
| \-r       | Recursively check any subfolders for any video files                                                                                   | false                     |
| \-c       | Copy the existing frames for formats that most likely support it, this is much faster as it saves re-encoding where possible           | false                     |
| \-d       | Delete old files after successful conversion                                                                                           | false                     |
| \-m       | Move old files after successful conversion                                                                                             | false                     |
| \-mt      | Folder to move old files to                                                                                                            |                           |
| \-exec    | Overwrite the FFmpeg executable, this is useful if FFmpeg is not in your path or the same directory                                    | ffmpeg                    |
| \-args    | Args to pass to FFmpeg, make sure to specify these as a string, example: `-args "-c:v h264_nvenc -preset p7 -tune hq -profile:v high"` |                           |
| \-e       | Exit on FFmpeg error                                                                                                                   | false                     |
| \-help    | Display help menu                                                                                                                      |                           |
| \-version | Display version number                                                                                                                 |                           |