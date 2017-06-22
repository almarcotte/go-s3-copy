// Package gossip
package gossip

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gnumast/gossip/config"
	"github.com/gnumast/gossip/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File represents a local file with all the options required for upload and cleanup
type File struct {
	Source      string
	Path        config.Path
	Destination string
}

// UploadQueue is an array of Files to be uploaded to S3
type UploadQueue []File

// Watcher contains everything needed for the program to run
type Watcher struct {
	Watcher    *fsnotify.Watcher
	Config     *config.Parsed
	Logger     log.StdOut
	Queue      UploadQueue
	CleanPaths map[string]config.Path
}

// NewWatcher initializes a watcher
func NewWatcher(config *config.Parsed, logger log.StdOut) *Watcher {
	return &Watcher{
		Config: config,
		Logger: logger,
		Queue:  UploadQueue{},
	}
}

// Start is the main loop where most of the action takes place
func (w *Watcher) Start() (err error) {
	w.Watcher, err = fsnotify.NewWatcher()
	defer w.Watcher.Close()

	if err != nil {
		return
	}

	done := make(chan bool)

	go func() {
		var delay <-chan time.Time

		for {
			select {
			case event := <-w.Watcher.Events:
				delay = time.After(w.Config.Global.Delay * time.Second)
				if shouldBreak := w.HandleEvent(event); shouldBreak {
					break
				}
			case <-delay: // We waited long enough, upload everything in the queue
				w.UploadQueue()
			}
		}
	}()

	// Navigate the registered paths, find their subfolders and add them to the cleaned paths so we can have access
	// to their parent's configs as needed
	w.CleanPaths = w.StartWatching()

	<-done

	return
}

// StartWatching navigates the configured paths and adds their subfolders if the top level path is recursive then,
// starts the watch process
func (w *Watcher) StartWatching() (cleaned map[string]config.Path) {
	cleaned = make(map[string]config.Path)

	for _, v := range w.Config.Paths {
		w.Logger.Printf("Adding watcher on %s\n", v.Root)

		if err := w.Watcher.Add(v.Root); err != nil {
			w.Logger.Error(err)
			continue
		}

		cleaned[v.Root] = v

		if v.Recursive {
			w.Logger.Printf("  Recursive, walking the tree to add subfolders")
			filepath.Walk(v.Root, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() && path != v.Root {
					w.Logger.Printf("  Adding watcher on sub-folders %s", path)
					cleaned[path] = v
					w.Watcher.Add(path)
				}

				return nil
			})
		}
	}

	return
}

func (w *Watcher) UploadQueue() {
	w.Logger.Printf("About to upload the queue")

	for _, v := range w.Queue {
		dir := filepath.Dir(v.Source)

		v.Destination = strings.TrimPrefix(strings.Replace(v.Source, w.CleanPaths[dir].Root, "", 1), "/")
		v.Path = w.CleanPaths[dir]

		err := w.UploadToS3(v)

		if err != nil {
			w.Logger.Error(err)
		}
	}

	// Empty out the queue
	w.Queue = UploadQueue{}
}

// HandleEvent handles a file event coming from fsnotify
func (w *Watcher) HandleEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return true
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		w.Logger.Printf("Found new file: %s", event.Name)
		w.Queue = append(w.Queue, File{Source: event.Name})
	}

	return false
}

// UploadToS3 uploads a file to the appropriate S3 bucket
func (w Watcher) UploadToS3(file File) (err error) {
	w.Logger.Printf("Uploading %s", file.Path)

	source, err := os.Open(file.Source)
	defer source.Close()

	if err != nil {
		return
	}

	fileInfo, err := source.Stat()

	if err != nil {
		return
	}

	size := fileInfo.Size()

	buffer := make([]byte, size)
	source.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(file.Path.Bucket),
		Key:           aws.String(file.Destination),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	sess, err := session.NewSession(w.configureAws())

	if err != nil {
		return
	}

	svc := s3.New(sess)
	resp, err := svc.PutObject(params)

	if err != nil {
		return
	}

	w.Logger.Printf("Response: %s", awsutil.StringValue(resp))

	if file.Path.Delete {
		w.Logger.Printf("Delete %s because parent directory has delete set to true", file.Source)
		err = os.Remove(file.Source)
	}

	return
}

func (w Watcher) configureAws() *aws.Config {
	creds := credentials.NewStaticCredentials(w.Config.Credentials.Access, w.Config.Credentials.Secret, "")
	_, err := creds.Get()

	if err != nil {
		w.Logger.Fatal(err)
	}

	return aws.NewConfig().WithRegion(w.Config.Credentials.Region).WithCredentials(creds)
}
