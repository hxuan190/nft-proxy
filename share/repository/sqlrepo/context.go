package sqlrepo

type DbContext interface {
	GetDBConnection() *gorm.DB
}

type dbContext struct {
	db *gorm.DB
}

func NewDbContext(rootDbConn *gorm.DB) DbContext {
	return &dbContext{rootDbConn}
}

func (dbc *dbContext) GetDBConnection() *gorm.DB {
	return dbc.db.Session(&gorm.Session{NewDB: true})
}
