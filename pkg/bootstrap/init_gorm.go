package bootstrap

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		c.User, c.Password, c.Host, c.Port, c.DBname,
	)
}

type ReadDB struct {
	DB    *gorm.DB
	Label string
}

type DBGroup struct {
	Write *gorm.DB
	Reads []ReadDB
}

func NewDBGroup(v *viper.Viper) *DBGroup {
	cfg := DBConfig{
		Host:     v.GetString("postgresql.write.host"),
		Port:     v.GetString("postgresql.write.port"),
		User:     v.GetString("postgresql.write.user"),
		Password: v.GetString("postgresql.write.password"),
		DBname:   v.GetString("postgresql.write.dbname"),
	}
	write, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("write db error: %v", err)
	}

	var reads []ReadDB
	readConfigs := v.Get("postgresql.reads").([]interface{})

	for _, cfg := range readConfigs {
		rc := cfg.(map[string]interface{})
		conf := DBConfig{
			Host:     rc["host"].(string),
			Port:     rc["port"].(string),
			User:     rc["user"].(string),
			Password: rc["password"].(string),
			DBname:   rc["dbname"].(string),
		}
		db, err := gorm.Open(postgres.Open(conf.DSN()), &gorm.Config{})
		if err != nil {
			log.Fatalf("read db error: %v", err)
		}
		reads = append(reads, ReadDB{DB: db, Label: conf.Host})
	}

	return &DBGroup{Write: write, Reads: reads}
}

func (d *DBGroup) PickDBLeastConnRead() *gorm.DB {
	var selected *gorm.DB
	minConn := int(^uint(0) >> 1)
	for _, read := range d.Reads {
		sqlDB, err := read.DB.DB()
		if err != nil {
			continue
		}
		stats := sqlDB.Stats()
		if stats.OpenConnections < minConn {
			minConn = stats.OpenConnections
			selected = read.DB
		}
	}
	return selected
}
