// Biocad task API
//
// The purpose of this API is to provide information over the devices stored in app_db.devices table
//
// Version: 1.0.0
//
// Contact: Dmitrii Rusinov<199-41@mail.ru>
//
// swagger:meta
package web

import (
	"biocadTask/internal/storage"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// swagger:model
type DeviceDto struct {
	ID        int64  `json:"id"`
	Num       int64  `json:"num"`
	Mqtt      string `json:"mqtt"`
	Invid     string `json:"invid"`
	UnitGUID  string `json:"unitGUID"`
	MsgID     string `json:"msgID"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	Class     string `json:"class"`
	Level     int64  `json:"level"`
	Area      string `json:"area"`
	Addr      string `json:"addr"`
	Block     bool   `json:"block"`
	Type      string `json:"type"`
	Bit       int64  `json:"bit"`
	InvertBit bool   `json:"invertBit"`
}

func entityToDTO(e storage.DeviceEntity) DeviceDto {
	return DeviceDto{
		ID:        e.ID,
		Num:       e.Num,
		Mqtt:      e.Mqtt,
		Invid:     e.Invid,
		UnitGUID:  e.UnitGUID.String(),
		MsgID:     e.MsgID,
		Text:      e.Text,
		Context:   e.Context,
		Class:     e.Class,
		Level:     e.Level,
		Area:      e.Area,
		Addr:      e.Addr,
		Block:     e.Block,
		Type:      e.Type,
		Bit:       e.Bit,
		InvertBit: e.InvertBit,
	}
}

type ApiHandler struct {
	messageDao storage.DeviceDao
}

func (ah *ApiHandler) ListenAndServe() error {
	router := mux.NewRouter()
	router.HandleFunc("/messages/{unit_guid}", ah.handleGetMessages).Methods("GET")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		return err
	}
	return nil
}

// swagger:route GET /messages/{unit_guid} messages getMessages
//
// Get messages by unit GUID.
//
// This will return all messages associated with the given unit GUID.
//
// responses:
//
//	200: []DeviceDto
//	400: BadRequestError
//	500: InternalServerError
func (ah *ApiHandler) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitGuid, err := uuid.Parse(vars["unit_guid"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageEntities, err := ah.messageDao.GetDevices(page, limit, unitGuid)
	if err != nil {
		log.Printf("error during request handling: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messageDtos []DeviceDto
	for _, entity := range messageEntities {
		messageDtos = append(messageDtos, entityToDTO(entity))
	}

	jsonData, err := json.Marshal(messageDtos)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		return
	}
}

func NewApiHandler(messageDao storage.DeviceDao) *ApiHandler {
	return &ApiHandler{messageDao: messageDao}
}

// swagger:model
type BadRequestError struct {
	// A short error code or description
	//
	// example: Invalid parameter
	Message string `json:"message"`
}

// swagger:model
type InternalServerError struct {
	// A short error code or description
	//
	// example: Server error
	Message string `json:"message"`
}

// swagger:parameters getMessages
type MessagesParams struct {
	// The GUID of the unit to get messages for
	//
	// in: path
	// required: true
	UnitGuid string `json:"unit_guid"`

	// The page number to retrieve
	//
	// in: query
	// required: true
	Page int64 `json:"page"`

	// The maximum number of messages to retrieve per page
	//
	// in: query
	// required: true
	Limit int64 `json:"limit"`
}
