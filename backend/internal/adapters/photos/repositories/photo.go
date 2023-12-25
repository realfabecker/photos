package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/realfabecker/photos/internal/adapters/common/dynamo"
	"github.com/realfabecker/photos/internal/adapters/common/validator"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

type PhotoDynamodDBRepository struct {
	db    *dynamodb.Client
	table string
	app   string
}

func NewWalletDynamoDBRepository(db *dynamodb.Client, table string, app string) (corpts.PhotoRepository, error) {
	return &PhotoDynamodDBRepository{db, table, app}, nil
}

func (u *PhotoDynamodDBRepository) ListPhotos(user string, q cordom.PhotoPagedDTOQuery) (*cordom.PagedDTO[cordom.Photo], error) {
	cipher := fmt.Sprintf("%s%d", user, q.Limit)
	k, err := dynamo.DecodePageToken(q.PageToken, cipher)
	if err != nil {
		return nil, err
	}

	var out *dynamodb.QueryOutput
	out, err = u.db.Query(context.TODO(), &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :v and begins_with(SK, :x)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#USER#" + user,
			},
			":x": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#PHOTO#" + q.GetPeriod(),
			},
		},
		ScanIndexForward:  aws.Bool(false),
		TableName:         aws.String(u.table),
		Limit:             &q.Limit,
		ExclusiveStartKey: k,
	})

	if err != nil {
		return nil, err
	}

	var lst []cordom.Photo
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &lst); err != nil {
		return nil, err
	}

	dto := cordom.PagedDTO[cordom.Photo]{}
	dto.PageCount = out.ScannedCount
	dto.Items = lst
	dto.HasMore = out.LastEvaluatedKey != nil

	if out.LastEvaluatedKey != nil {
		if dto.PageToken, err = dynamo.EncodePageToken(
			out.LastEvaluatedKey,
			cipher,
		); err != nil {
			return nil, err
		}
	}
	return &dto, nil
}

func (u *PhotoDynamodDBRepository) CreatePhoto(p *cordom.Photo) (*cordom.Photo, error) {
	pd := photo{Photo: p}

	pd.PhotoId = time.Now().Format("20060102") + validator.NewULID(time.Now())
	pd.CreatedAt = time.Now().Format("2006-01-02T15:04:05-07:00")

	pd.PK = "APP#" + u.app + "#USER#" + p.UserId
	pd.SK = "APP#" + u.app + "#PHOTO#" + p.PhotoId

	avs, err := attributevalue.MarshalMap(pd)
	if err != nil {
		return nil, err
	}

	if _, err := u.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(u.table),
		Item:      avs,
	}); err != nil {
		return nil, err
	}

	return pd.Photo, nil
}

func (u *PhotoDynamodDBRepository) GetPhotoById(user string, photo string) (*cordom.Photo, error) {
	out, err := u.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#USER#" + user,
			},
			"SK": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#PHOTO#" + photo,
			},
		},
		TableName: aws.String(u.table),
	})

	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var dto cordom.Photo
	if err := attributevalue.UnmarshalMap(out.Item, &dto); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (u *PhotoDynamodDBRepository) DeletePhoto(user string, photo string) error {
	_, err := u.db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#USER#" + user,
			},
			"SK": &types.AttributeValueMemberS{
				Value: "APP#" + u.app + "#PHOTO#" + photo,
			},
		},
		TableName: aws.String(u.table),
	})
	return err
}

func (u *PhotoDynamodDBRepository) PutPhoto(p *cordom.Photo) (*cordom.Photo, error) {
	r, err := u.GetPhotoById(p.UserId, p.PhotoId)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return u.CreatePhoto(p)
	}

	pd := photo{Photo: p}
	pd.PK = "APP#" + u.app + "#USER#" + p.UserId
	pd.SK = "APP#" + u.app + "#PHOTO#" + p.PhotoId

	avs, err := attributevalue.MarshalMap(pd)
	if err != nil {
		return nil, err
	}

	if _, err := u.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(u.table),
		Item:      avs,
	}); err != nil {
		return nil, err
	}

	return p, nil
}
