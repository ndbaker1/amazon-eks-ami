package configprovider

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	internalapi "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	apibridge "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api/bridge"
)

const (
	eksBootstrapConfigKey   = "name"
	eksBootstrapConfigValue = "eks-bootstrap-config"
	userDataMediaType       = "multipart/mixed"
)

type imdsConfigProvider struct{}

func NewIMDSConfigProvider() ConfigProvider {
	return &imdsConfigProvider{}
}

func (ics *imdsConfigProvider) Provide() (*internalapi.NodeConfig, error) {
	userDataOutput, err := imds.New(imds.Options{}).GetUserData(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	userData, err := io.ReadAll(userDataOutput.Content)
	if err != nil {
		return nil, err
	}

	data, err := parseMIMEUserData(userData)
	if err != nil {
		return nil, err
	}

	return apibridge.DecodeNodeConfig(data)
}

func parseMIMEUserData(mimedata []byte) ([]byte, error) {
	tp := textproto.NewReader(bufio.NewReader(bytes.NewReader(mimedata)))
	header, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	mediaType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if mediaType != userDataMediaType {
		return nil, fmt.Errorf("UserData MIME Content-Type was %s and needs to be %s", mediaType, userDataMediaType)
	}

	boundary := params["boundary"]
	userDataReader := multipart.NewReader(bytes.NewReader(mimedata), boundary)

	for {
		part, err := userDataReader.NextPart()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("Could not find %s data within the IMDS UserData", eksBootstrapConfigValue)
			}

			return nil, err
		}

		if partHeader := part.Header.Get("Content-Disposition"); partHeader != "" {
			_, params, err := mime.ParseMediaType(partHeader)
			if err != nil {
				return nil, err
			}

			if params[eksBootstrapConfigKey] == eksBootstrapConfigValue {
				configData, err := io.ReadAll(part)
				if err != nil {
					return nil, err
				}

				return configData, nil
			}
		}
	}
}
