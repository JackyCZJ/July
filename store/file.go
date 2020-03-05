package store

import (
	"io"

	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func Upload(reader io.Reader, filename string) (string, error) {
	bucket, err := gridfs.NewBucket(Client.db)
	if err != nil {
		return "", err
	}
	fileid, err := bucket.UploadFromStream(filename, reader)
	if err != nil {
		return "", err
	}

	return fileid.Hex(), err
}

func Download(id primitive.ObjectID, writer io.Writer) error {
	bucket, err := gridfs.NewBucket(Client.db)
	if err != nil {
		return err
	}
	_, err = bucket.DownloadToStream(id, writer)
	if err != nil {
		return err
	}
	return nil
}

func Public() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		file := ctx.Get(":filename")
		return ctx.File("/public" + file.(string))
	}
}
