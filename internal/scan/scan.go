package scan

import (
	"biocadTask/internal/queue"
	"biocadTask/internal/storage"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StartScanner scans the inputDirectory for new files in infinite loop.
func StartScanner(
	inputDirectory string,
	fileQueue queue.FileMessageQueue,
	preprocessedFileDao storage.PreprocessedFileDao,
) {
	for {
		time.Sleep(5 * time.Second)
		log.Print("scanning for new files")

		files, err := os.ReadDir(inputDirectory)

		if err != nil {
			log.Printf("could not scan for files: %v", err)
			continue
		}

		for _, file := range files {
			filePath := inputDirectory + "/" + file.Name()

			preprocessed, queryError := preprocessedFileDao.WasPreprocessed(filePath)
			if queryError != nil {
				log.Printf("could not check if file was preprocessed: %v", queryError)
				continue
			}

			if filepath.Ext(file.Name()) == ".tsv" && !preprocessed {
				log.Printf("sending file to rabbit: %s", file.Name())
				produceError := fileQueue.Produce(filePath)
				if produceError != nil {
					log.Printf("could not send to rabbit:  %v", produceError)
					continue
				} else {
					queryError = preprocessedFileDao.Add(filePath)
					if queryError != nil {
						log.Printf("could not save to preprocessed_file: :  %v", queryError)
						continue
					}
				}
			}
		}
	}
}
