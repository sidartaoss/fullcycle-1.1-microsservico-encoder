package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/sidartaoss/fullcycle/encoder/application/repositories"
	"github.com/sidartaoss/fullcycle/encoder/domain"
)

type VideoService struct {
	*domain.Video
	repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (s *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(s.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.mp4", os.Getenv("localStoragePath"), s.Video.ID))
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("video %s has been stored", s.Video.ID)
	return nil
}

func (v *VideoService) Fragment() error {

	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 file", v.Video.ID, ".mp4")
		return err
	}
	err = os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Println("error removing frag file", v.Video.ID, ".frag")
		return err
	}
	err = os.RemoveAll(os.Getenv("localStoragePath") + "/" + v.Video.ID)
	if err != nil {
		log.Println("error removing folder", v.Video.ID)
		return err
	}
	log.Println("files have been removed", v.Video.ID)
	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}
	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("==========> Output: %s\n", string(out))
	}
}
