package media

type MediaForm struct {
	Name        string `json:"name"`
	MediaTypeId int    `json:"mediaTypeId"`
	LocationId  int    `json:"locationId"`
	Path        string `json:"path"`
}
