package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	proto "ipCamera/proto/ipcamera/proto"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Stream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {

	fmt.Println("Client-", req.Data)
	return &proto.StreamResponse{Response: "We got your message"}, nil
}

func (s *Service) StreamImage(ctx context.Context, req *proto.StreamImageRequest) (*proto.StreamImageResponse, error) {

	fmt.Println("Client- Sent picture")
	// save picture
	serveFrames(req.Image)
	return &proto.StreamImageResponse{Response: "We got your picture"}, nil
}

func serveFrames(imgByte []byte) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create("./img.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Println(err)
	}

}
