package mediatype

import "github.com/gotk/pg"

func GetById(db *pg.Session, id int) (*Model, error) {
	entity, err := db.FindOne(&Model{}, "id=$1", id)
	if err != nil {
		return nil, err
	}
	return entity.(*Model), nil
}
