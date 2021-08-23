package azure

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
)

var gAzure *Azure

type Azure struct {
	accountName string
	accountKey  string
	azureDomain string
	//containerNft     string
	//containerProduct string

	credential          *azblob.SharedKeyCredential
	pipeline            pipeline.Pipeline
	containerNftUrl     azblob.ContainerURL
	containerProductUrl azblob.ContainerURL
}

func GetAzure() *Azure {
	return gAzure
}

func InitAzure(azureAccount, azureAccessKey, azureDomain, containerNft, containerProduct string) (*Azure, error) {
	log.Info("InitAzure start...")

	gAzure = &Azure{
		accountName: azureAccount,
		accountKey:  azureAccessKey,
		azureDomain: azureDomain,
	}
	if len(azureAccount) == 0 || len(azureAccessKey) == 0 {
		err := errors.New("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
		log.Fatal(err)
		return nil, err
	}

	if credential, err := azblob.NewSharedKeyCredential(gAzure.accountName, gAzure.accountKey); err != nil {
		return nil, err
	} else {
		gAzure.credential = credential
	}

	gAzure.pipeline = azblob.NewPipeline(gAzure.credential, azblob.PipelineOptions{})

	//make nft, product container url
	gAzure.MakeNftContainerUrl(containerNft)
	gAzure.MakeProductContainerUrl(containerProduct)

	log.Info("InitAzure success!")
	return gAzure, nil
}

func (o *Azure) MakeNftContainerUrl(containerName string) {
	url, _ := url.Parse(fmt.Sprintf("%s%s", o.azureDomain, containerName))
	gAzure.containerNftUrl = azblob.NewContainerURL(*url, o.pipeline)
}

func (o *Azure) MakeProductContainerUrl(containerName string) {
	url, _ := url.Parse(fmt.Sprintf("%s%s", o.azureDomain, containerName))
	gAzure.containerProductUrl = azblob.NewContainerURL(*url, o.pipeline)
}

func (o *Azure) UploadNftInfoBuffer(b []byte, remoteFileName string) error {
	ctx := context.Background()

	blobURL := o.containerNftUrl.NewBlockBlobURL(remoteFileName)
	options := azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16}
	_, err := azblob.UploadBufferToBlockBlob(ctx, b, blobURL, options)

	return err
}

func (o *Azure) UploadFile() {

}
