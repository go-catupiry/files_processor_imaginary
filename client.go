package files_processor_imaginary

import (
	"fmt"
	"os"

	files_processor "github.com/go-catupiry/files/processor"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ClientConfiguration struct {
	URL string
}

type Client struct {
	HTTP *req.Client
	Cfg  *ClientConfiguration
}

func NewClient(cfg *ClientConfiguration) *Client {
	c := Client{
		HTTP: req.C(),
		Cfg:  cfg,
	}

	return &c
}

func (c *Client) Resize(sourcePath, destPath, fileName string, opts files_processor.Options) error {
	if _, ok := opts["url"]; ok {
		err := c.ResizeFromWeb(sourcePath, destPath, fileName, opts)
		if err != nil {
			return err
		}
	} else {
		err := c.ResizeFromLocalhost(sourcePath, destPath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) ResizeFromWeb(sourcePath, destPath, fileName string, opts files_processor.Options) error {
	url := c.Cfg.URL + "/resize"

	if f, ok := opts["format"]; ok {
		opts["type"] = f
	} else if _, ok := opts["type"]; !ok {
		opts["type"] = "webp"
	}

	if _, ok := opts["nocrop"]; !ok {
		opts["nocrop"] = "false"
	}

	res, err := c.HTTP.R().
		SetQueryParams(opts).
		SetOutputFile(destPath).
		Get(url)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("ResizeFromWeb error", err)
		return errors.Wrap(err, "error on access imaginary api")
	}

	if res.IsErrorState() {
		logrus.WithFields(logrus.Fields{
			"error":  fmt.Sprintf("%+v\n", err),
			"status": res.StatusCode,
			"dump":   res.Dump(),
		}).Error("ResizeFromWeb Response error", err)

		return errors.New(res.String())
	}

	return nil
}

func (c *Client) ResizeFromLocalhost(sourcePath string, destPath string, opts files_processor.Options) error {
	url := c.Cfg.URL + "/resize"

	f, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrap(err, "resizeFromLocalhost error on open file from sourcePath")
	}
	defer f.Close()

	if f, ok := opts["format"]; ok {
		opts["type"] = f
	} else if _, ok := opts["type"]; !ok {
		opts["type"] = "webp"
	}

	if _, ok := opts["nocrop"]; !ok {
		opts["nocrop"] = "false"
	}

	id := uuid.New()

	res, err := c.HTTP.R().
		SetQueryParams(opts).
		SetFileReader("file", id.String(), f).
		SetOutputFile(destPath).
		Post(url)

	// execution error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("resizeFromLocalhost error on access imaginary api", err)
		return errors.Wrap(err, "resizeFromLocalhost error on access imaginary api")
	}

	// http error
	if res.IsError() {
		logrus.WithFields(logrus.Fields{
			"error":  fmt.Sprintf("%+v\n", err),
			"status": res.StatusCode,
			"dump":   res.Dump(),
		}).Error("resizeFromLocalhost response error", err)
		return errors.New(res.String())
	}

	return nil
}

// DownloadFile
// Usage:
// originalPath := path.Join(os.TempDir(), fileName) + "_original"
// defer os.Remove(originalPath)
// DownloadFile(fileURL, originalPath, fileName string) (error)
func (c *Client) DownloadFile(fileURL, donwloadedFilePath, fileName string) error {
	res, err := c.HTTP.R().
		SetOutputFile(donwloadedFilePath).
		Get(fileURL)

	// execution error
	if err != nil {
		return errors.Wrap(err, "resizeFromLocalhost error on download original image")
	}
	// http error
	if res.IsError() {
		logrus.WithFields(logrus.Fields{
			"error":  fmt.Sprintf("%+v\n", err),
			"status": res.StatusCode,
			"body":   res.String(),
		}).Error("resizeFromLocalhost Response error", err)

		return errors.New(res.String())
	}

	return nil
}

// // TODO!

// crop - Same as /crop endpoint.
// smartcrop - Same as /smartcrop endpoint.
// enlarge - Same as /enlarge endpoint.
// extract - Same as /extract endpoint.
// rotate - Same as /rotate endpoint.
// autorotate - Same as /autorotate endpoint.
// flip - Same as /flip endpoint.
// flop - Same as /flop endpoint.
// thumbnail - Same as /thumbnail endpoint.
// zoom - Same as /zoom endpoint.
// convert - Same as /convert endpoint.
// watermark - Same as /watermark endpoint.
// watermarkimage - Same as /watermarkimage endpoint.
// blur - Same as /blur endpoint.
