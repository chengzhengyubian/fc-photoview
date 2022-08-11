package scanner_tasks

/*修改完,未测试*/
import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"log"
	"strconv"
	"strings"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/scanner/scanner_tasks/processing_tasks"
	"github.com/pkg/errors"
)

//未修改
type VideoMetadataTask struct {
	scanner_task.ScannerTaskBase
}

func (t VideoMetadataTask) AfterMediaFound(ctx scanner_task.TaskContext, media *models.Media, newMedia bool) error {

	if !newMedia || media.Type != models.MediaTypeVideo {
		return nil
	}

	err := ScanVideoMetadata( /*ctx.GetDB(), */ media)
	if err != nil {
		log.Printf("WARN: ScanVideoMetadata for %s failed: %s\n", media.Title, err)
	}

	return nil
}

func ScanVideoMetadata( /*tx *gorm.DB, */ video *models.Media) error {

	data, err := processing_tasks.ReadVideoMetadata(video.Path)
	if err != nil {
		return errors.Wrapf(err, "scan video metadata failed (%s)", video.Title)
	}

	stream := data.FirstVideoStream()
	if stream == nil {
		return errors.New(fmt.Sprintf("could not get video stream from metadata (%s)", video.Path))
	}

	audio := data.FirstAudioStream()
	var audioText string
	if audio == nil {
		audioText = "No audio"
	} else {
		switch audio.Channels {
		case 0:
			audioText = "No audio"
		case 1:
			audioText = "Mono audio"
		case 2:
			audioText = "Stereo audio"
		default:
			audioText = fmt.Sprintf("Audio (%d channels)", audio.Channels)
		}
	}

	var framerate *float64 = nil
	if stream.AvgFrameRate != "" {
		parts := strings.Split(stream.AvgFrameRate, "/")
		if len(parts) == 2 {
			if numerator, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
				if denominator, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					result := float64(numerator) / float64(denominator)
					framerate = &result
				}
			}
		}
	}

	videoMetadata := models.VideoMetadata{
		Width:        stream.Width,
		Height:       stream.Height,
		Duration:     data.Format.DurationSeconds,
		Codec:        &stream.CodecLongName,
		Framerate:    framerate,
		Bitrate:      &stream.BitRate,
		ColorProfile: &stream.Profile,
		Audio:        &audioText,
	}

	video.VideoMetadata = &videoMetadata

	//未修改
	//if err := tx.Save(video).Error; err != nil {
	//	return errors.Wrapf(err, "failed to add video metadata to database (%s)", video.Title)
	//}
	var codec string
	if video.VideoMetadata.Codec == nil {
		codec = "NULL"
	} else {
		codec = *video.VideoMetadata.Codec
	}
	var fram string
	if video.VideoMetadata.Framerate == nil {
		fram = "NULL"
	} else {
		num := int(*video.VideoMetadata.Framerate)
		fram = strconv.Itoa(num)
	}
	var bitrate string
	if video.VideoMetadata.Bitrate == nil {
		bitrate = "NULL"
	} else {
		bitrate = *video.VideoMetadata.Bitrate
	}
	var colorProfile string
	if video.VideoMetadata.ColorProfile == nil {
		colorProfile = "NULL"
	} else {
		colorProfile = *video.VideoMetadata.ColorProfile
	}
	var aud string
	if video.VideoMetadata.Audio == nil {
		aud = "NULL"
	} else {
		aud = *video.VideoMetadata.Audio
	}
	sql_video_metadata := fmt.Sprintf("update video_metadata set updated_at=NOW(),width=%v,height=%v,duration=%v,codec='%v',framerate=%v,bitrate='%v',color_profile='%v',audio='%v'", video.VideoMetadata.Width, video.VideoMetadata.Height, video.VideoMetadata.Duration, codec, fram, bitrate, colorProfile, aud)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_video_metadata)
	return nil
}
