package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	proto "ipCamera/proto/ipcamera/proto"
	"log"
	"os"
)

type Service struct {
	log *log.Logger
}

func NewService(l *log.Logger) *Service {
	return &Service{l}
}

func (s *Service) Stream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {

	// Received message from client
	s.log.Println("Received message from client.")
	s.log.Println("Client:", req.Data)

	// Respond to client
	return &proto.StreamResponse{Response: "I got your message"}, nil
}

func (s *Service) StreamImage(ctx context.Context, req *proto.StreamImageRequest) (*proto.StreamImageResponse, error) {

	// Received image from client
	s.log.Println("Received image from client.")

	// save image
	serveFrames(req.Image)

	// Respond to client
	return &proto.StreamImageResponse{Response: "I got your image"}, nil
}

func (s *Service) StreamVideo(stream proto.IpCamera_StreamVideoServer) error {

	_, err := stream.Recv()

	if err != nil {
		s.log.Fatal("Error fetching stream")
	}

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		s.log.Println("Waiting for more chunks")

		req, err := stream.Recv()
		if err == io.EOF {
			s.log.Println("EOF Stop stream")
			break
		}
		if err != nil {
			return fmt.Errorf("Cannot receive h264 chunk: %v", err)
		}

		// process chunk
		s.log.Println("Received a h264 chunk")

		// h264 chunk
		chunk := req.GetVideo()

		// Process Chunk
		// todo..

		log.Println(chunk)

	}

	res := &proto.StreamVideoResponse{
		Response: "Finish Streaming.",
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return fmt.Errorf("Cannot send response: %v", err)
	}

	s.log.Println("Done Streaming")

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
