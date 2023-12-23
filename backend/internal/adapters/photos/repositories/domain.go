package repositories

import (
	cordom "github.com/realfabecker/photos/internal/core/domain"
)

type photo struct {
	*cordom.Photo
	PK     string `dynamodbav:"PK,omitempty" json:"-"`
	SK     string `dynamodbav:"SK,omitempty" json:"-"`
	GSI1PK string `dynamodbav:"GSI1_PK,omitempty" json:"-"`
	GSI1SK string `dynamodbav:"GSI1_SK,omitempty" json:"-"`
}
