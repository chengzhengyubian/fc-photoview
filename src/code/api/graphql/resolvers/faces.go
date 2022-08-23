package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"

	api "github.com/photoview/photoview/api/graphql"
	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/pkg/errors"
	//"gorm.io/gorm"
)

type imageFaceResolver struct {
	*Resolver
}

type faceGroupResolver struct {
	*Resolver
}

func (r *Resolver) ImageFace() api.ImageFaceResolver {
	return imageFaceResolver{r}
}

func (r *Resolver) FaceGroup() api.FaceGroupResolver {
	return faceGroupResolver{r}
}

func (r imageFaceResolver) FaceGroup(ctx context.Context, obj *models.ImageFace) (*models.FaceGroup, error) {
	if obj.FaceGroup != nil {
		return obj.FaceGroup, nil
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	var faceGroup models.FaceGroup
	sql_facegroup_se := fmt.Sprintf("select face_groups.* from face_groups left join image_faces on image_faces.face_group_id=face_groups.id where image_faces.id=%v", obj.ID)
	dataApi, _ := dataapi.NewDataApiClient()
	res, _ := dataApi.Query(sql_facegroup_se)
	if len(res) == 0 {
		return nil, nil
	}
	faceGroup.ID = dataapi.GetInt(res, 0, 0)
	faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, 0, 1)/1000, 0)
	faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, 0, 2)/1000, 0)
	faceGroup.Label = dataapi.GetStringP(res, 0, 3)
	obj.FaceGroup = &faceGroup

	return &faceGroup, nil
}

func (r imageFaceResolver) Media(ctx context.Context, obj *models.ImageFace) (*models.Media, error) {
	if err := obj.FillMedia(); err != nil {
		return nil, err
	}

	return &obj.Media, nil
}

func (r faceGroupResolver) ImageFaces(ctx context.Context, obj *models.FaceGroup, paginate *models.Pagination) ([]*models.ImageFace, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}

	var imageFaces []*models.ImageFace
	var limit int
	limit = *paginate.Limit
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_image_faces_se := fmt.Sprintf("select image_faces.* from image_faces left join media on image_faces.media_id=media.id where image_faces.face_group_id=%v and image_faces.album_id in (%v) limit %v", obj.ID, userAlbumids, limit)
	dataApi, _ := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_image_faces_se)
	if len(res) == 0 {
		return nil, err
	}
	for i := 0; i < len(res); i++ {
		var image models.ImageFace
		image.ID = dataapi.GetInt(res, i, 0)
		image.CreatedAt = time.Unix(dataapi.GetLong(res, i, 1)/1000, 0)
		image.UpdatedAt = time.Unix(dataapi.GetLong(res, i, 2)/1000, 0)
		image.FaceGroupID = dataapi.GetInt(res, i, 3)
		image.MediaID = dataapi.GetInt(res, i, 4)
		//image.descriptor=res[i][5].BlobValue.Read()
		//image.Rectangle=dataapi.GetStringP(res,i,6)
		imageFaces = append(imageFaces, &image)
	}
	return imageFaces, nil
}

func (r faceGroupResolver) ImageFaceCount(ctx context.Context, obj *models.FaceGroup) (int, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return -1, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return -1, errors.New("face detector not initialized")
	}

	if err := user.FillAlbums(); err != nil {
		return -1, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}
	var count int64
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_count_se := fmt.Sprintf("select count(image_faces.*) from image_faces left join media on image_faces.media_id=media.id where image_faces.face_group_id=%v and image_faces.albums_id in(%v)", obj.ID, userAlbumids)
	dataApi, _ := dataapi.NewDataApiClientJosn()
	res, err := dataApi.Query(sql_count_se)
	if len(res) == 0 {
		return -1, err
	}
	count = dataapi.GetLong(res, 0, 0)
	return int(count), nil
}

func (r *queryResolver) FaceGroup(ctx context.Context, id int) (*models.FaceGroup, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}
	var faceGroup models.FaceGroup
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_media_se := fmt.Sprintf("select face_groups.* from face_groups LEFT JOIN image_faces ON image_faces.face_group_id = face_groups.id LEFT JOIN media ON image_faces.media_id = media.id where face_groups.id =%v and media.album_id IN (%v)", id, userAlbumids)
	dataApi, _ := dataapi.NewDataApiClientJosn()
	res, err := dataApi.Query(sql_media_se)
	if len(res) == 0 {
		return nil, err
	}
	faceGroup.ID = dataapi.GetInt(res, 0, 0)
	faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, 0, 1)/1000, 0)
	faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, 0, 2)/1000, 0)
	faceGroup.Label = dataapi.GetStringP(res, 0, 3)
	return &faceGroup, nil
}

func (r *queryResolver) MyFaceGroups(ctx context.Context, paginate *models.Pagination) ([]*models.FaceGroup, error) {

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}

	var faceGroups []*models.FaceGroup
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_face_groups_se := fmt.Sprintf("SELECT `face_groups`.`id`,`face_groups`.`created_at`,`face_groups`.`updated_at`,`face_groups`.`label` FROM `face_groups` JOIN image_faces ON image_faces.face_group_id = face_groups.id WHERE image_faces.media_id IN (SELECT media.id FROM `media` WHERE media.album_id IN (%v)) GROUP BY `image_faces`.`face_group_id`,`face_groups`.`id` ORDER BY CASE WHEN label IS NULL THEN 1 ELSE 0 END,COUNT(image_faces.id) DESC LIM", userAlbumids)
	dataApi, _ := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_face_groups_se)
	if len(res) == 0 {
		return nil, err
	}
	for i := 0; i < len(res); i++ {
		var faceGroup models.FaceGroup
		faceGroup.ID = dataapi.GetInt(res, i, 0)
		faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, i, 1)/1000, 0)
		faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, i, 2)/1000, 0)
		faceGroup.Label = dataapi.GetStringP(res, i, 3)
		faceGroups = append(faceGroups, &faceGroup)
	}
	return faceGroups, nil
}

func (r *mutationResolver) SetFaceGroupLabel(ctx context.Context, faceGroupID int, label *string) (*models.FaceGroup, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	faceGroup, err := userOwnedFaceGroup( /*db, */ user, faceGroupID)
	if err != nil {
		return nil, err
	}
	sql_face_groups_up := fmt.Sprintf("update face_groups set label='%v' where face_groups.id=%v", faceGroup.ID, label)
	dataApi, _ := dataapi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_face_groups_up)
	return faceGroup, nil
}

func (r *mutationResolver) CombineFaceGroups(ctx context.Context, destinationFaceGroupID int, sourceFaceGroupID int) (*models.FaceGroup, error) {

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	destinationFaceGroup, err := userOwnedFaceGroup( /*db,*/ user, destinationFaceGroupID)
	if err != nil {
		return nil, err
	}

	sourceFaceGroup, err := userOwnedFaceGroup( /*db, */ user, sourceFaceGroupID)
	if err != nil {
		return nil, err
	}

	dataApi, _ := dataapi.NewDataApiClientJosn()
	sql_image_faces_up := fmt.Sprintf("update image_faces set face_group_id=%v where id=%v", destinationFaceGroup.ID, models.ImageFace{}.ID)
	dataApi.ExecuteSQl(sql_image_faces_up)
	sql_face_groups_de := fmt.Sprintf("delete from face_groups where id=%v", sourceFaceGroup.ID)
	dataApi.ExecuteSQl(sql_face_groups_de)

	face_detection.GlobalFaceDetector.MergeCategories(int32(sourceFaceGroupID), int32(destinationFaceGroupID))

	return destinationFaceGroup, nil
}

func (r *mutationResolver) MoveImageFaces(ctx context.Context, imageFaceIDs []int, destinationFaceGroupID int) (*models.FaceGroup, error) {

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	userOwnedImageFaceIDs := make([]int, 0)
	var destFaceGroup *models.FaceGroup
	dataApi, _ := dataapi.NewDataApiClient()
	//var err error
	destFaceGroup, err := userOwnedFaceGroup( /*tx, */ user, destinationFaceGroupID)
	if err != nil {
		return nil, err
	}

	userOwnedImageFaces, err := getUserOwnedImageFaces( /*tx,*/ user, imageFaceIDs)

	for _, imageFace := range userOwnedImageFaces {
		userOwnedImageFaceIDs = append(userOwnedImageFaceIDs, imageFace.ID)
	}

	var sourceFaceGroups []*models.FaceGroup
	userOwnedImageFaceID, _ := json.Marshal(userOwnedImageFaceIDs)
	userOwnedImageFaceids := strings.Trim(string(userOwnedImageFaceID), "[]")
	sql_face_groups_se := fmt.Sprintf("select face_groups.* from face_groups LEFT JOIN image_faces ON image_faces.face_group_id = face_groups.id where image_faces.id IN (%v)", userOwnedImageFaceids)
	res, err := dataApi.Query(sql_face_groups_se)
	for i := 0; i < len(res); i++ {
		var faceGroup models.FaceGroup
		faceGroup.ID = dataapi.GetInt(res, i, 0)
		faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, i, 1)/1000, 0)
		faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, i, 2)/1000, 0)
		faceGroup.Label = dataapi.GetStringP(res, i, 3)
		sourceFaceGroups = append(sourceFaceGroups, &faceGroup)
	}
	userOwnedImageFaceID, _ = json.Marshal(userOwnedImageFaceIDs)
	userOwnedImageFaceids = strings.Trim(string(userOwnedImageFaceID), "[]")
	sql_face_groups_up := fmt.Sprintf("update face_groups set face_group_id=%v where id in(%v)", userOwnedImageFaceids)
	dataApi.ExecuteSQl(sql_face_groups_up)
	// delete face groups if they have become empty
	for _, faceGroup := range sourceFaceGroups {
		var count int64
		sql_face_groups_se := fmt.Sprintf("select count(*) from face_groups where face_group_id=%v", faceGroup.ID)
		res, _ := dataApi.Query(sql_face_groups_se)
		count = dataapi.GetLong(res, 0, 0)
		if count == 0 {
			sql_delete := fmt.Sprintf("delete from face_groups where face_groups.id=%v", faceGroup.ID)
			dataApi.ExecuteSQl(sql_delete)
		}
	}

	face_detection.GlobalFaceDetector.MergeImageFaces(userOwnedImageFaceIDs, int32(destFaceGroup.ID))

	return destFaceGroup, nil
}

func (r *mutationResolver) RecognizeUnlabeledFaces(ctx context.Context) ([]*models.ImageFace, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	var updatedImageFaces []*models.ImageFace

	updatedImageFaces, err := face_detection.GlobalFaceDetector.RecognizeUnlabeledFaces( /*tx,*/ user)
	if err != nil {
		return nil, err
	}
	return updatedImageFaces, nil
}

func (r *mutationResolver) DetachImageFaces(ctx context.Context, imageFaceIDs []int) (*models.FaceGroup, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil, errors.New("face detector not initialized")
	}

	userOwnedImageFaceIDs := make([]int, 0)
	newFaceGroup := models.FaceGroup{}

	userOwnedImageFaces, err := getUserOwnedImageFaces( /*tx,*/ user, imageFaceIDs)
	if err != nil {
		return nil, err
	}

	for _, imageFace := range userOwnedImageFaces {
		userOwnedImageFaceIDs = append(userOwnedImageFaceIDs, imageFace.ID)
	}
	sql_face_groups_up := fmt.Sprintf("update face_groups set updated_at=NOW() set lable='%v' id= %v", newFaceGroup.Label, newFaceGroup.ID)
	dataApi, _ := dataapi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_face_groups_up)
	userOwnedImageFaceID, _ := json.Marshal(userOwnedImageFaceIDs)
	userOwnedImageFaceids := strings.Trim(string(userOwnedImageFaceID), "[]")
	sql_image_faces_up := fmt.Sprintf("update image_faces set face_group_id=%v where image_faces.id in (%v)", newFaceGroup.ID, userOwnedImageFaceids)
	//dataApi, _ := dataapi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_image_faces_up)

	face_detection.GlobalFaceDetector.MergeImageFaces(userOwnedImageFaceIDs, int32(newFaceGroup.ID))

	return &newFaceGroup, nil
}

func userOwnedFaceGroup(user *models.User, faceGroupID int) (*models.FaceGroup, error) {
	dataApi, _ := dataapi.NewDataApiClient()
	if user.Admin {
		var faceGroup models.FaceGroup
		sql_face_groups_se := fmt.Sprintf("select * from face_groups where id=%v", faceGroupID)
		res, err := dataApi.Query(sql_face_groups_se)
		if len(res) == 0 {
			return nil, err
		}
		faceGroup.ID = dataapi.GetInt(res, 0, 0)
		faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, 0, 1)/1000, 0)
		faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, 0, 2)/1000, 0)
		faceGroup.Label = dataapi.GetStringP(res, 0, 3)
		return &faceGroup, nil
	}

	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}
	var faceGroup models.FaceGroup
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_image_faces_se := fmt.Sprintf("select image_faces.* from image_faces left join media on image_faces.media_id=media.id where  image_faces.album_id in (%v) limit %v", userAlbumids)
	sql_face_groups_se := fmt.Sprintf("select * from face_groups.* from face_groups JOIN image_faces ON face_groups.id = image_faces.face_group_id where face_groups.id = %v and image_faces.id IN(%v)", faceGroupID, sql_image_faces_se)

	res, err := dataApi.Query(sql_face_groups_se)
	if len(res) == 0 {
		return nil, err
	}
	faceGroup.ID = dataapi.GetInt(res, 0, 0)
	faceGroup.CreatedAt = time.Unix(dataapi.GetLong(res, 0, 1)/1000, 0)
	faceGroup.UpdatedAt = time.Unix(dataapi.GetLong(res, 0, 2)/1000, 0)
	faceGroup.Label = dataapi.GetStringP(res, 0, 3)
	return &faceGroup, nil
}

func getUserOwnedImageFaces( /*tx *gorm.DB, */ user *models.User, imageFaceIDs []int) ([]*models.ImageFace, error) {
	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}

	var userOwnedImageFaces []*models.ImageFace
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	imageFaceID, _ := json.Marshal(imageFaceIDs)
	imageFaceids := strings.Trim(string(imageFaceID), "[]")
	sql_image_faces_se := fmt.Sprintf("select image_faces.* from image_faces JOIN media ON media.id = image_faces.media_id and media.album_id IN (%v) and image_faces.id IN(%v)", userAlbumids, imageFaceids)
	dataApi, _ := dataapi.NewDataApiClient()
	dataApi.Query(sql_image_faces_se)
	return userOwnedImageFaces, nil
}
