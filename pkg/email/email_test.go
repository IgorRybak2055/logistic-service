package email

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testEmailService struct {
	testDir      string
	templatePath string
	messageCh    chan MessageData
}

func (t testEmailService) Run() {
	var ctx = context.TODO()

LOOP:
	for {
		select {
		case msg := <-t.messageCh:
			go handleSend(t, msg, "testURL", t.templatePath)
		case <-ctx.Done():
			break LOOP
		}
	}
}

func (t testEmailService) send(e MessageData) error {
	file, err := os.Create(filepath.Join(t.testDir, "letter.txt"))
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Println("closing test file:", err)
		}
	}()

	if _, err = file.WriteString(e.message); err != nil {
		return err
	}

	return nil
}

// SendToQueue send data to rabbitMQ channel
func SendToQueue(recvEmail, hashID string, ch chan MessageData) {
	ch <- MessageData{
		RecvEmail: recvEmail,
		UserID:    hashID,
		UserToken: "testUserToken",
	}
}

func TestMailService(t *testing.T) {
	var (
		srv = testEmailService{
			messageCh: make(chan MessageData, 1),
			templatePath: filepath.Join("."),
		}
		err error
	)

	srv.testDir, err = ioutil.TempDir("testdata", "email")
	require.NoError(t, err)

	go srv.Run()

	SendToQueue("testemail@mail.com", "hashID", srv.messageCh)

	time.Sleep(1 * time.Second)

	defer func() {
		require.NoError(t, os.RemoveAll(srv.testDir))
	}()

	expectedByte, err := ioutil.ReadFile(filepath.Join("testdata", "letter.golden"))
	require.NoError(t, err)

	file, err := os.Open(filepath.Join(srv.testDir, "letter.txt"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, file.Close())
	}()

	data := make([]byte, len(expectedByte))

	n, err := file.Read(data)
	require.NoError(t, err)

	strFromFile := data[:n]
	require.Equal(t, expectedByte, strFromFile)
}
