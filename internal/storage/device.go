package storage

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// DeviceEntity is a dao for table of main data in the application
type DeviceEntity struct {
	ID        int64
	Num       int64
	Mqtt      string
	Invid     string
	UnitGUID  uuid.UUID
	MsgID     string
	Text      string
	Context   string
	Class     string
	Level     int64
	Area      string
	Addr      string
	Block     bool
	Type      string
	Bit       int64
	InvertBit bool
}

type DeviceDao interface {
	AddDevice(file DeviceEntity) error
	GetDevices(page int64, pageSize int64, unitGuid uuid.UUID) ([]DeviceEntity, error)
}

type PostgresDeviceDao struct {
	db *sql.DB
}

func NewPostgresFileDao(db *sql.DB) *PostgresDeviceDao {
	return &PostgresDeviceDao{db}
}

func (e *DeviceEntity) String() string {
	return fmt.Sprintf(
		`ID: %d
Num: %d
Mqtt: %s
Invid: %s
UnitGUID: %s
MsgID: %s
Text: %s
Context: %s
Class: %s
Level: %d
Area: %s
Addr: %s
Block: %t
Type: %s
Bit: %d
InvertBit: %t`,
		e.ID,
		e.Num,
		e.Mqtt,
		e.Invid,
		e.UnitGUID.String(),
		e.MsgID,
		e.Text,
		e.Context,
		e.Class,
		e.Level,
		e.Area,
		e.Addr,
		e.Block,
		e.Type,
		e.Bit,
		e.InvertBit,
	)
}

func (dao *PostgresDeviceDao) GetDevices(page int64, pageSize int64, unitGuid uuid.UUID) ([]DeviceEntity, error) {
	offset := (page - 1) * pageSize
	query := `
        SELECT 
            id, num, mqtt, invid, unit_guid, msg_id, text,
            context, class, level, area, addr, block, type, bit, invert_bit
        FROM 
            device
        WHERE 
            unit_guid=$3
        ORDER BY 
            id
        OFFSET $1
        LIMIT $2
    `
	rows, err := dao.db.Query(query, offset, pageSize, unitGuid)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var devices []DeviceEntity
	for rows.Next() {
		var device DeviceEntity
		err = rows.Scan(
			&device.ID,
			&device.Num,
			&device.Mqtt,
			&device.Invid,
			&device.UnitGUID,
			&device.MsgID,
			&device.Text,
			&device.Context,
			&device.Class,
			&device.Level,
			&device.Area,
			&device.Addr,
			&device.Block,
			&device.Type,
			&device.Bit,
			&device.InvertBit,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (dao *PostgresDeviceDao) AddDevice(entity DeviceEntity) error {
	query := `
        INSERT INTO 
            device (num, mqtt, invid, unit_guid, msg_id, text, context, class, level, area, addr, block, type, bit, invert_bit)
        VALUES 
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `

	_, err := dao.db.Exec(query,
		entity.Num,
		entity.Mqtt,
		entity.Invid,
		entity.UnitGUID,
		entity.MsgID,
		entity.Text,
		entity.Context,
		entity.Class,
		entity.Level,
		entity.Area,
		entity.Addr,
		entity.Block,
		entity.Type,
		entity.Bit,
		entity.InvertBit,
	)

	if err != nil {
		return err
	}

	return nil
}
