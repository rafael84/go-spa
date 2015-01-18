package media

type MediaForm struct {
	Name        string `json:"name"`
	MediatypeId int    `json:"mediatypeId"`
	LocationId  int    `json:"locationId"`
	Path        string `json:"path"`
}
