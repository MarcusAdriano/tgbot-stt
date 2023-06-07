package audioconverter_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/marcusadriano/sound-stt-tgbot/internal/audioconverter"
	"github.com/marcusadriano/sound-stt-tgbot/internal/fileserver"
	"github.com/marcusadriano/sound-stt-tgbot/internal/mocks"
)

type cmdRunnerMock struct {
}

func (c *cmdRunnerMock) Run(name string, args ...string) error {
	return nil
}

func TestFfmpegToMp3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileServer := mocks.NewMockFileserver(ctrl)
	fileServer.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(&fileserver.FilePath{Path: "fileName"}, nil).
		Times(1)

	fileServer.
		EXPECT().
		Read(gomock.Any(), "fileName.mp3").
		Return(&fileserver.File{Data: []byte("fileData")}, nil).
		Times(1)

	fileServer.
		EXPECT().
		Delete(gomock.Any(), gomock.Eq("fileName.mp3")).
		Return(nil).
		Times(1)

	converter := audioconverter.NewFfmpegWithCmdRunner(fileServer, &cmdRunnerMock{})
	result, err := converter.ToMp3(context.TODO(), []byte("fileData"), "fileName")

	if err != nil {
		t.Fatalf("Error to convert file: %s", err)
	}

	if string(result.Data) != "fileData" {
		t.Fatalf("Expected fileData, got %s", result.Data)
	}
}
