package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"image"
	"strconv"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/photoview/photoview/api/scanner/media_encoding/media_utils"
	//"gorm.io/gorm"
)

type FaceGroup struct {
	Model
	Label      *string
	ImageFaces []ImageFace `gorm:"constraint:OnDelete:CASCADE;"`
}

type ImageFace struct {
	Model
	FaceGroupID int `gorm:"not null;index"`
	FaceGroup   *FaceGroup
	MediaID     int            `gorm:"not null;index"`
	Media       Media          `gorm:"constraint:OnDelete:CASCADE;"`
	Descriptor  FaceDescriptor `gorm:"not null"`
	Rectangle   FaceRectangle  `gorm:"not null"`
}

func (f *ImageFace) FillMedia() error {
	if f.Media.ID != 0 {
		// media already exists
		return nil
	}
	sql_media_se := fmt.Sprintf("select * from media left join image_faces on image_faces.media_id=media.id where image_faces.id=%v", f.ID)
	dataApi, _ := dataapi.NewDataApiClientJosn()
	res, _ := dataApi.Query(sql_media_se)
	f.Media.ID = dataapi.GetInt(res, 0, 0)
	f.Media.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	f.Media.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	f.Media.Title = *res[0][3].StringValue
	f.Media.Path = *res[0][4].StringValue
	f.Media.PathHash = *res[0][5].StringValue
	f.Media.AlbumID = int(*res[0][6].LongValue)
	f.Media.ExifID = dataapi.GetIntP(res, 0, 7)
	f.Media.DateShot = time.Unix(*res[0][8].LongValue/1000, 0)
	if *res[0][9].StringValue == "photo" {
		f.Media.Type = MediaTypePhoto
	} else {
		f.Media.Type = MediaTypeVideo
	}
	f.Media.VideoMetadataID = dataapi.GetIntP(res, 0, 10)
	f.Media.SideCarPath = dataapi.GetStringP(res, 0, 11)
	f.Media.SideCarHash = dataapi.GetStringP(res, 0, 12)
	f.Media.Blurhash = dataapi.GetStringP(res, 0, 13)
	f.Media.AlbumID = dataapi.GetInt(res, 0, 14)
	f.Media.Album.ID = dataapi.GetInt(res, 0, 14)
	f.Media.Album.CreatedAt = time.Unix(*res[0][15].LongValue/1000, 0)
	f.Media.Album.UpdatedAt = time.Unix(*res[0][16].LongValue/1000, 0)
	f.Media.Album.Title = dataapi.GetString(res, 0, 17)
	f.Media.Album.ParentAlbumID = dataapi.GetIntP(res, 0, 18)
	f.Media.Album.Path = dataapi.GetString(res, 0, 19)
	f.Media.Album.PathHash = dataapi.GetString(res, 0, 20)
	f.Media.Album.CoverID = dataapi.GetIntP(res, 0, 21)
	return nil
}

type FaceDescriptor face.Descriptor

// Scan tells GORM how to convert database data to Go format
func (fd *FaceDescriptor) Scan(value interface{}) error {
	byteValue := value.([]byte)
	reader := bytes.NewReader(byteValue)
	binary.Read(reader, binary.LittleEndian, fd)
	return nil
}

// Value tells GORM how to save into the database
func (fd FaceDescriptor) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, fd); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type FaceRectangle struct {
	MinX, MaxX float64
	MinY, MaxY float64
}

// ToDBFaceRectangle converts a pixel absolute rectangle to a relative FaceRectangle to be saved in the database
func ToDBFaceRectangle(imgRec image.Rectangle, imagePath string) (*FaceRectangle, error) {
	size, err := media_utils.GetPhotoDimensions(imagePath)
	if err != nil {
		return nil, err
	}

	return &FaceRectangle{
		MinX: float64(imgRec.Min.X) / float64(size.Width),
		MaxX: float64(imgRec.Max.X) / float64(size.Width),
		MinY: float64(imgRec.Min.Y) / float64(size.Height),
		MaxY: float64(imgRec.Max.Y) / float64(size.Height),
	}, nil
}

// GormDataType datatype used in database
func (fr FaceRectangle) GormDataType() string {
	return "VARCHAR(64)"
}

// Scan tells GORM how to convert database data to Go format
func (fr *FaceRectangle) Scan(value interface{}) error {
	stringArray, ok := value.(string)
	if !ok {
		byteArray := value.([]uint8)
		stringArray = string(byteArray)
	}

	slices := strings.Split(stringArray, ":")

	if len(slices) != 4 {
		return fmt.Errorf("Invalid face rectangle format, expected 4 values, got %d", len(slices))
	}

	var err error

	fr.MinX, err = strconv.ParseFloat(slices[0], 32)
	if err != nil {
		return err
	}

	fr.MaxX, err = strconv.ParseFloat(slices[1], 32)
	if err != nil {
		return err
	}

	fr.MinY, err = strconv.ParseFloat(slices[2], 32)
	if err != nil {
		return err
	}

	fr.MaxY, err = strconv.ParseFloat(slices[3], 32)
	if err != nil {
		return err
	}

	return nil
}

// Value tells GORM how to save into the database
func (fr FaceRectangle) Value() (driver.Value, error) {
	result := fmt.Sprintf("%f:%f:%f:%f", fr.MinX, fr.MaxX, fr.MinY, fr.MaxY)
	return result, nil
}
