package goty

import "gorm.io/gorm"

type Response struct {
	Error error
	RowsAffected int64
}

func migrate(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(&model)

	if err != nil {
		panic(err)
	}
}

func getDB(model interface{}) *gorm.DB {
	db, config := getConnection()
	if config!=nil && config.migrate {
		migrate(db, model)
	}

	return db

}

func Save(model interface{}) Response {
	db := getDB(model)
	response := db.Create(model)
	return dbResponse(response)

}

func dbResponse(response *gorm.DB) Response {
	return Response{
		Error: response.Error,
		RowsAffected: response.RowsAffected,
	}
}

func SaveOmit(model interface{}, omit []string) Response {
	db := getDB(model)
	response := db.Omit(omit...).Create(model)
	return dbResponse(response)
}

func SaveSelected(model interface{}, selected []string) Response {
	db := getDB(model)
	response := db.Select(model, selected).Create(model)
	return dbResponse(response)
}

func BatchSave(models []interface{}, batchSize int) Response {
	db := getDB(models)
	var response *gorm.DB
	if batchSize>0 {
		response = db.Create(models)
	}
	response = db.CreateInBatches(models, batchSize)
	return dbResponse(response)
}