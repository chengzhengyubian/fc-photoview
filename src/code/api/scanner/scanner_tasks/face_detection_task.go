package scanner_tasks

import (
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/scanner_task"
)

type FaceDetectionTask struct {
	scanner_task.ScannerTaskBase
}

func (t FaceDetectionTask) AfterProcessMedia(ctx scanner_task.TaskContext, mediaData *media_encoding.EncodeMediaData, updatedURLs []*models.MediaURL, mediaIndex int, mediaTotal int) error {
	didProcess := len(updatedURLs) > 0

	if didProcess && mediaData.Media.Type == models.MediaTypePhoto {
		go func(media *models.Media) {
		}(mediaData.Media)
	}

	return nil
}
