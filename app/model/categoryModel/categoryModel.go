package categoryModel

import (
	"time"

	"github.com/go-pg/pg/v10/types"
)

type Category struct {
	tableName struct{} `pg:"categories,alias:c"`

	ID        int                 `pg:"type:serial,pk" json:"id,omitempty"`
	Name      string              `pg:"type:varchar(255),unique,notnull" json:"name,omitempty"`
	Vectors   types.ValueAppender `pg:"type:tsvector,notnull,default:''" json:"-"`
	CreatedAt *time.Time          `pg:"type:timestamp,notnull,default:current_timestamp" json:"createdAt,omitempty"`
	UpdatedAt *time.Time          `pg:"type:timestamp,notnull,default:current_timestamp" json:"-"`
}

const (
	CategoryTable = "categories"
	CategoryAlias = "c"
)
