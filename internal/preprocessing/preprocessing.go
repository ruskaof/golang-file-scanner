package preprocessing

import (
	"biocadTask/internal/storage"
	"encoding/csv"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
)

type FilePreprocessor struct {
	outputDirectoryPath string
	deviceDao           storage.DeviceDao
	errorDao            storage.ErrorDao
}

func (fp *FilePreprocessor) CreateOutputDir() error {
	err := os.MkdirAll(fp.outputDirectoryPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func NewFilePreprocessor(
	outputDirectoryPath string,
	deviceDao storage.DeviceDao,
	errorDao storage.ErrorDao,
) *FilePreprocessor {
	return &FilePreprocessor{outputDirectoryPath: outputDirectoryPath, deviceDao: deviceDao, errorDao: errorDao}
}

func (fp *FilePreprocessor) PreprocessFile(
	filePath string,
) error {

	log.Printf("preprocessing file %v", filePath)

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	tsvReader := csv.NewReader(file)
	tsvReader.Comma = '\t'

	data, err := tsvReader.ReadAll()

	if err != nil {
		return err
	}

	var entitiesByGuid = make(map[uuid.UUID][]storage.DeviceEntity)
	var errorsById = make(map[int64]error)

	for i := 0; i < len(data); i++ {
		entity, entityParseError := rowToEntity(data[i])
		if entityParseError != nil {
			id, insertErrorErr := fp.errorDao.Add(entityParseError.Error())
			if insertErrorErr != nil {
				return insertErrorErr
			}
			errorsById[id] = entityParseError
			continue
		}

		daoInsertError := fp.deviceDao.AddDevice(entity)
		if daoInsertError != nil {
			id, insertErrorErr := fp.errorDao.Add(daoInsertError.Error())
			if insertErrorErr != nil {
				return insertErrorErr
			}
			errorsById[id] = daoInsertError
			continue
		}

		entitiesByGuid[entity.UnitGUID] = append(entitiesByGuid[entity.UnitGUID], entity)
	}

	for guid, entities := range entitiesByGuid {
		err = GenerateProcessedFileReports(entities, guid, fp.outputDirectoryPath)
		if err != nil {
			return err
		}
	}

	for id, errorMessage := range errorsById {
		err = GenerateErrorReport(errorMessage, id, fp.outputDirectoryPath)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func rowToEntity(row []string) (storage.DeviceEntity, error) {
	intNum, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	uuidUnitGuid, err := uuid.Parse(row[3])
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	intLevel, err := strconv.ParseInt(row[8], 10, 64)
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	boolBlock, err := strconv.ParseBool(row[11])
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	intBit, err := strconv.ParseInt(row[13], 10, 64)
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	boolInvertBit, err := strconv.ParseBool(row[14])
	if err != nil {
		return storage.DeviceEntity{}, err
	}

	return storage.DeviceEntity{
		ID:        0,
		Num:       intNum,
		Mqtt:      row[1],
		Invid:     row[2],
		UnitGUID:  uuidUnitGuid,
		MsgID:     row[4],
		Text:      row[5],
		Context:   row[6],
		Class:     row[7],
		Level:     intLevel,
		Area:      row[9],
		Addr:      row[10],
		Block:     boolBlock,
		Type:      row[12],
		Bit:       intBit,
		InvertBit: boolInvertBit,
	}, nil
}
