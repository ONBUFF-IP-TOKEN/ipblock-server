package azure

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

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

// cdn byte 버퍼로 업로드
// remotefileName에 경로를 포함하면 azure cdn에 가상 폴더를 자동으로 생성해준다.
func (o *Azure) UploadNftInfoBuffer(b []byte, remoteFileName string) error {
	blobURL := o.containerNftUrl.NewBlockBlobURL(remoteFileName)
	options := azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16}
	_, err := azblob.UploadBufferToBlockBlob(context.Background(), b, blobURL, options)

	return err
}

// cdn file로 업로드
// remotefileName에 경로를 포함하면 azure cdn에 가상 폴더를 자동으로 생성해준다.
func (o *Azure) UploadNftFile(file *os.File, remoteFileName string) error {
	blobURL := o.containerNftUrl.NewBlockBlobURL(remoteFileName)
	options := azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16}
	_, err := azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, options)
	return err
}

// nft 컨테이너 삭제
func (o *Azure) DeleteNftContainer() {
	o.containerNftUrl.Delete(context.Background(), azblob.ContainerAccessConditions{})
}

// product 컨테이너 삭제
func (o *Azure) DeleteProductContainer() {
	o.containerProductUrl.Delete(context.Background(), azblob.ContainerAccessConditions{})
}

// cdn 파일 삭제
func (o *Azure) DeleteNftFile(remoteFileName string) {
	blobURL := o.containerNftUrl.NewBlockBlobURL(remoteFileName)
	blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
}
