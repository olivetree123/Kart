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
	maxSize := params["maxSize"].(int64)
	volume := storage.NewVolumeModel(dirPath, maxSize, maxSize)
	global.DBConn.Insert("VolumeModel", volume)
	storage.NewFreeSectionModel(volume.ID.Value, 0, maxSize)
	utils.JSONResponse(database.ModelToMap(volume), w)
}
