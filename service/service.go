package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
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

func (s *Service) StreamVideo(stream proto.IpCamera_StreamVideoServer) error {

	_, err := stream.Recv()

	if err != nil {
		log.Fatal("Error fetching stream")
	}

	videoData := bytes.Buffer{}

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return fmt.Errorf("cannot receive chunk data: %v", err)
		}

		chunk := req.GetVideo()

		// process chunk

		_, err = videoData.Write(chunk)
		if err != nil {
			return fmt.Errorf("cannot write chunk data: %v", err)
		}
	}

	res := &proto.StreamVideoResponse{
		Response: "Finish Streaming.",
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return fmt.Errorf("cannot send response: %v", err)
	}

	fmt.Println("Done Streaming")

	return nil

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

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return fmt.Errorf("canceled")
	case context.DeadlineExceeded:
		return fmt.Errorf("DeadlineExceeded")
	default:
		return nil
	}
}
