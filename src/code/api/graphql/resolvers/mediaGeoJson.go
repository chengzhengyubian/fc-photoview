package resolvers

import (
	"context"
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"os"
	"path"

	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/utils"
)

type geoMedia struct {
	MediaID         int
	MediaTitle      string
	ThumbnailName   string
	ThumbnailWidth  int
	ThumbnailHeight int
	Latitude        float64
	Longitude       float64
}

type geoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []geoJSONFeature `json:"features"`
}

type geoJSONFeature struct {
	Type       string                 `json:"type"`
	Properties interface{}            `json:"properties"`
	Geometry   geoJSONFeatureGeometry `json:"geometry"`
}

type geoJSONMediaProperties struct {
	MediaID    int    `json:"media_id"`
	MediaTitle string `json:"media_title"`
	Thumbnail  struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
}

type geoJSONFeatureGeometry struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

func makeGeoJSONFeatureCollection(features []geoJSONFeature) geoJSONFeatureCollection {
	return geoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

func makeGeoJSONFeature(properties interface{}, geometry geoJSONFeatureGeometry) geoJSONFeature {
	return geoJSONFeature{
		Type:       "Feature",
		Properties: properties,
		Geometry:   geometry,
	}
}

func makeGeoJSONFeatureGeometryPoint(lat float64, long float64) geoJSONFeatureGeometry {
	coordinates := [2]float64{long, lat}

	return geoJSONFeatureGeometry{
		Type:        "Point",
		Coordinates: coordinates,
	}
}

func (r *queryResolver) MyMediaGeoJSON(ctx context.Context) (interface{}, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	var media []*geoMedia

	sql_media_se := fmt.Sprintf("SELECT media.id AS media_id, media.title AS media_title, media_urls.media_name AS thumbnail_name, media_urls.width AS thumbnail_width, media_urls.height AS thumbnail_height, media_exif.gps_latitude AS latitude, media_exif.gps_longitude AS longitude FROM `media` INNER JOIN media_exif ON media.exif_id = media_exif.id INNER JOIN media_urls ON media.id = media_urls.media_id INNER JOIN user_albums ON media.album_id = user_albums.album_id WHERE media_exif.gps_latitude IS NOT NULL AND media_exif.gps_longitude IS NOT NULL AND media_urls.purpose = 'thumbnail' AND user_albums.user_id = %v", user.ID)
	dataApi, _ := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	if len(res) == 0 {
		return nil, err
	}
	for i := 0; i < num; i++ {
		var Media geoMedia
		Media.MediaID = dataapi.GetInt(res, 0, 0)
		Media.MediaTitle = dataapi.GetString(res, 0, 1)
		Media.ThumbnailName = dataapi.GetString(res, 0, 2)
		Media.ThumbnailWidth = dataapi.GetInt(res, 0, 3)
		Media.ThumbnailHeight = dataapi.GetInt(res, 0, 4)
		var lati float64
		lati = float64(dataapi.GetLong(res, 0, 5))
		Media.Latitude = lati
		var long float64
		long = float64(dataapi.GetLong(res, 0, 6))
		Media.Longitude = long
		media = append(media, &Media)
	}
	if err != nil {
		return nil, err
	}

	features := make([]geoJSONFeature, 0)

	for _, item := range media {
		geoPoint := makeGeoJSONFeatureGeometryPoint(item.Latitude, item.Longitude)

		thumbnailURL := utils.ApiEndpointUrl()
		thumbnailURL.Path = path.Join(thumbnailURL.Path, "photo", item.ThumbnailName)

		properties := geoJSONMediaProperties{
			MediaID:    item.MediaID,
			MediaTitle: item.MediaTitle,
			Thumbnail: struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			}{
				URL:    thumbnailURL.String(),
				Width:  item.ThumbnailWidth,
				Height: item.ThumbnailHeight,
			},
		}

		features = append(features, makeGeoJSONFeature(properties, geoPoint))
	}

	featureCollection := makeGeoJSONFeatureCollection(features)
	return featureCollection, nil
}

func (r *queryResolver) MapboxToken(ctx context.Context) (*string, error) {
	mapboxTokenEnv := os.Getenv("MAPBOX_TOKEN")
	if mapboxTokenEnv == "" {
		return nil, nil
	}

	return &mapboxTokenEnv, nil
}
