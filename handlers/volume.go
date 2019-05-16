package handlers

import (
	"Kart/database"
	"Kart/global"
	"Kart/storage"
	"Kart/utils"
	"fmt"
	"net/http"
)

// AddVolumeHandler 添加 Volume
func AddVolumeHandler(w http.ResponseWriter, r *http.Request) {
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	dirPath := params["dirPath"].(string)
	maxSize := params["maxSize"].(float64)
	volume := storage.NewVolumeModel(dirPath, int64(maxSize), int64(maxSize))
	global.DBConn.Insert("VolumeModel", volume)
	section := storage.NewFreeSectionModel(volume.ID.Value, 0, int64(maxSize))
	global.DBConn.Insert("FreeSectionModel", section)
	utils.JSONResponse(database.ModelToMap(volume), w)
}
