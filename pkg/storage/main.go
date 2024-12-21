package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureStorage struct {
	containerClient *azb	lob
	accountName     string
	accountKey      string
}

func NewAzureStorage(accountName, accountKey string) (*AzureStorage, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	containerClient, err := azblob.NewContainerClient(fmt.Sprintf("https://%s.blob.core.windows.net/messages", accountName), credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client: %v", err)
	}

	return &AzureStorage{
		containerClient: containerClient,
		accountName:     accountName,
		accountKey:      accountKey,
	}, nil
}

func (a *AzureStorage) UploadProfilePhoto(ctx context.Context, userID string, photo io.Reader) (string, error) {
	blobClient := a.containerClient.NewBlobClient(fmt.Sprintf("profiles/%s.jpg", userID))
	_, err := blobClient.Upload(ctx, photo, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload profile photo: %v", err)
	}
	return blobClient.URL(), nil
}

func (a *AzureStorage) UploadChatAttachment(ctx context.Context, chatID string, filename string, content io.Reader) (string, error) {
	blobClient := a.containerClient.NewBlobClient(fmt.Sprintf("chats/%s/%s", chatID, filename))
	_, err := blobClient.Upload(ctx, content, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload attachment: %v", err)
	}
	return blobClient.URL(), nil
}

func (a *AzureStorage) GetSignedURL(blobPath string, duration time.Duration) (string, error) {
	blobClient := a.containerClient.NewBlobClient(blobPath)
	sasURL, err := blobClient.GenerateSasURL(azblob.BlobSASPermissions{Read: true}, time.Now().Add(duration), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate SAS URL: %v", err)
	}
	return sasURL, nil
}
