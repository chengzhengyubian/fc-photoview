package models

//修改完
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"
)

type Album struct {
	Model
	Title         string `gorm:"not null"`
	ParentAlbumID *int   `gorm:"index"`
	ParentAlbum   *Album `gorm:"constraint:OnDelete:SET NULL;"`
	// OwnerID       int `gorm:"not null"`
	// Owner         User
	Owners   []User `gorm:"many2many:user_albums;constraint:OnDelete:CASCADE;"`
	Path     string `gorm:"not null"`
	PathHash string `gorm:"unique"`
	CoverID  *int
}

func (a *Album) FilePath() string {
	return a.Path
}

func (a *Album) BeforeSave( /*tx *gorm.DB*/ ) (err error) {
	hash := md5.Sum([]byte(a.Path))
	a.PathHash = hex.EncodeToString(hash[:])
	return nil
}

// GetChildren performs a recursive query to get all the children of the album.
// An optional filter can be provided that can be used to modify the query on the children.
func (a *Album) GetChildren( /*db *gorm.DB*/ filter func(sql string) string) (children []*Album, err error) {
	return GetChildrenFromAlbums(filter, []int{a.ID})
}

//修改完，待测试
func GetChildrenFromAlbums( /*db *gorm.DB*/ filter func(sql string) string, albumIDs []int) (children []*Album, err error) {
	//query := db.Model(&Album{}).Table("sub_albums")
	sql_albums_se := "select * from sub_albums"
	//标记一下，应该是递归
	if filter != nil {
		//query = filter(query)
		sql_albums_se = filter(sql_albums_se)
	}
	//err = db.Raw(`
	//WITH recursive sub_albums AS (
	//	SELECT * FROM albums AS root WHERE id IN (?)
	//	UNION ALL
	//	SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id
	//)
	//
	//?
	//`, albumIDs, query).Find(&children).Error
	/* WITH recursive sub_albums AS (
	           SELECT * FROM albums AS root WHERE id IN (110)
	           UNION ALL
	           SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id
	   )

	   SELECT * FROM `sub_albums`
	*/
	albumID, _ := json.Marshal(albumIDs)
	albumids := strings.Trim(string(albumID), "[]")
	sql_albums_se = fmt.Sprintf("WITH recursive sub_albums AS (SELECT * FROM albums AS root WHERE id IN (%v)  UNION ALL  SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id) %v", albumids, sql_albums_se)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var album Album
		album.ID = DataApi.GetInt(res, i, 0)
		album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, i, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
		album.Path = DataApi.GetString(res, i, 5)
		album.PathHash = DataApi.GetString(res, i, 6)
		album.CoverID = DataApi.GetIntP(res, i, 7)
		children = append(children, &album)
	}
	return children, err
}

func (a *Album) GetParents(filter func(sql string) string) (parents []*Album, err error) {
	return GetParentsFromAlbums(filter, a.ID)
}

//修改完，待测试
func GetParentsFromAlbums(filter func(sql string) string, albumID int) (parents []*Album, err error) {
	//query := db.Model(&Album{}).Table("super_albums")
	sql_albums_se := "select * from super_albums"
	if filter != nil {
		sql_albums_se = filter(sql_albums_se)
	}

	//err = db.Raw(`
	//WITH recursive super_albums AS (
	//	SELECT * FROM albums AS leaf WHERE id = ?
	//	UNION ALL
	//	SELECT parent.* from albums AS parent JOIN super_albums ON parent.id = super_albums.parent_album_id
	//)
	//
	//?
	//`, albumID, sql_albums_se).Find(&parents).Error
	/*WITH recursive super_albums AS (
	          SELECT * FROM albums AS leaf WHERE id = 1
	          UNION ALL
	          SELECT parent.* from albums AS parent JOIN super_albums ON parent.id = super_albums.parent_album_id
	  )

	  SELECT * FROM `super_albums` WHERE id IN (1)*/
	sql_albums_se = fmt.Sprintf("WITH recursive super_albums AS (SELECT * FROM albums AS leaf WHERE id = %v UNION ALL SELECT parent.* from albums AS parent JOIN super_albums ON parent.id = super_albums.parent_album_id) %v", albumID, sql_albums_se)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var album Album
		album.ID = DataApi.GetInt(res, i, 0)
		album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, i, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
		album.Path = DataApi.GetString(res, i, 5)
		album.PathHash = DataApi.GetString(res, i, 6)
		album.CoverID = DataApi.GetIntP(res, i, 7)
		parents = append(parents, &album)
	}
	return parents, err
}

//修改完，暂时没问题

func (a *Album) Thumbnail( /*db *gorm.DB*/ ) (*Media, error) {
	var media Media

	if a.CoverID == nil {
		//if err := db.Raw(`
		//	WITH recursive sub_albums AS (
		//		SELECT * FROM albums AS root WHERE id = ?
		//		UNION ALL
		//		SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id
		//	)
		//
		//	SELECT * FROM media WHERE media.album_id IN (
		//		SELECT id FROM sub_albums
		//	) AND media.id IN (
		//		SELECT media_id FROM media_urls WHERE media_urls.media_id = media.id
		//	) ORDER BY id LIMIT 1
		//`, a.ID).Find(&media).Error; err != nil {
		//	return nil, err
		//}
		/* WITH recursive sub_albums AS (
		           SELECT * FROM albums AS root WHERE id = 1
		           UNION ALL
		           SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id
		   )

		   SELECT * FROM media WHERE media.album_id IN (
		           SELECT id FROM sub_albums
		   ) AND media.id IN (
		           SELECT media_id FROM media_urls WHERE media_urls.media_id = media.id
		   ) ORDER BY id LIMIT 1
		*/
		sql_media_se := fmt.Sprintf("WITH recursive sub_albums AS (SELECT * FROM albums AS root WHERE id = %v UNION ALL SELECT child.* FROM albums AS child JOIN sub_albums ON child.parent_album_id = sub_albums.id) SELECT * FROM media WHERE media.album_id IN (SELECT id FROM sub_albums) AND media.id IN (SELECT media_id FROM media_urls WHERE media_urls.media_id = media.id) ORDER BY id LIMIT 1", a.ID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_se)
		if len(res) == 0 {
			return nil, err
		}
		media.ID = DataApi.GetInt(res, 0, 0)
		media.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
		media.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
		media.Title = *res[0][3].StringValue
		media.Path = *res[0][4].StringValue
		media.PathHash = *res[0][5].StringValue
		media.AlbumID = int(*res[0][6].LongValue)
		media.ExifID = DataApi.GetIntP(res, 0, 7)
		media.DateShot = time.Unix(*res[0][8].LongValue/1000, 0)
		if *res[0][9].StringValue == "photo" {
			media.Type = MediaTypePhoto
		} else {
			media.Type = MediaTypeVideo
		}
		media.VideoMetadataID = DataApi.GetIntP(res, 0, 10)
		media.SideCarPath = DataApi.GetStringP(res, 0, 11)
		media.SideCarHash = DataApi.GetStringP(res, 0, 12)
		media.Blurhash = DataApi.GetStringP(res, 0, 13)
		if err != nil {
			return nil, err
		}
	} else {
		//if err := db.Where("id = ?", a.CoverID).Find(&media).Error; err != nil {
		//	return nil, err
		//}
		id := *a.CoverID
		sql_media_se := fmt.Sprintf("select * from media where id=%v", id)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_se)
		media.ID = DataApi.GetInt(res, 0, 0)
		media.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
		media.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
		media.Title = *res[0][3].StringValue
		media.Path = *res[0][4].StringValue
		media.PathHash = *res[0][5].StringValue
		media.AlbumID = int(*res[0][6].LongValue)
		media.ExifID = DataApi.GetIntP(res, 0, 7)
		media.DateShot = time.Unix(*res[0][8].LongValue/1000, 0)
		if *res[0][9].StringValue == "photo" {
			media.Type = MediaTypePhoto
		} else {
			media.Type = MediaTypeVideo
		}
		media.VideoMetadataID = DataApi.GetIntP(res, 0, 10)
		media.SideCarPath = DataApi.GetStringP(res, 0, 11)
		media.SideCarHash = DataApi.GetStringP(res, 0, 12)
		media.Blurhash = DataApi.GetStringP(res, 0, 13)
		if err != nil {
			return nil, err
		}
	}

	return &media, nil
}
