package bootstrap

import (
	"context"
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

type WriteDB struct {
	DB    *gorm.DB
	Label string
}

type DBGroup struct {
	Write          *WriteDB
	Reads          []*ReadDB
	PreviousWriter *ReadDB // 存放上次被降級的writer
}

func NewDBGroup(v *viper.Viper) *DBGroup {
	cfg := DBConfig{
		Host:     v.GetString("postgresql.write.host"),
		Port:     v.GetString("postgresql.write.port"),
		User:     v.GetString("postgresql.write.user"),
		Password: v.GetString("postgresql.write.password"),
		DBname:   v.GetString("postgresql.write.dbname"),
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("write db error: %v", err)
	}
	write := &WriteDB{DB: db, Label: cfg.Host}

	var reads []*ReadDB
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
		reads = append(reads, &ReadDB{DB: db, Label: conf.Host})
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

func (d *DBGroup) HeadlthCheck(ctx context.Context) {
	// 檢查 write 狀態
	if err := pingNode(ctx, d.Write.DB); err != nil {
		log.Printf("[failover] write db failed: %v", err)

		// 將失效的writer 站存為 previousWriter
		d.PreviousWriter = &ReadDB{DB: d.Write.DB, Label: d.Write.Label}

		// 尋找第一個健康的 reader 做 promote
		for i, read := range d.Reads {
			if err := pingNode(ctx, read.DB); err == nil {
				log.Printf("[failover] promoting reader %s to writer", read.Label)
				d.Write.DB = read.DB
				d.Reads = append(d.Reads[:i], d.Reads[i+1:]...)
				break
			}
		}
	} else {
		if d.PreviousWriter != nil {
			if err := pingNode(ctx, d.PreviousWriter.DB); err == nil {
				log.Printf("[recover] previous writer recoverd, added as reader")
				d.Reads = append(d.Reads, d.PreviousWriter)
				d.PreviousWriter = nil
			}
		}
	}

}

func pingNode(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
