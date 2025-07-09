package helpers

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"time"
	"github.com/Azer-ch/web-scraper/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func InitDB() error {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return errors.New("MYSQL_DSN not set in environment")
	}
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&types.Analysis{})
}

func HashURL(url string) []byte {
	hash := sha256.Sum256([]byte(url))
	return hash[:]
}

func GetCachedResult(url string) (*types.AnalyzeResponse, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}
	var cache types.Analysis
	urlHash := HashURL(url)
	result := db.Where("url_hash = ?", urlHash).First(&cache)
	if result.Error != nil {
		return nil, nil
	}
	if time.Since(cache.AnalyzedAt) > 24*time.Hour {
		return nil, nil
	}
	var resp types.AnalyzeResponse
	if err := json.Unmarshal([]byte(cache.Result), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func SetCachedResult(url string, result *types.AnalyzeResponse) error {
	if db == nil {
		if err := InitDB(); err != nil {
			return err
		}
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	urlHash := HashURL(url)
	cache := types.Analysis{
		URL:        url,
		URLHash:    urlHash,
		Result:     string(resultBytes),
		AnalyzedAt: time.Now(),
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{"result", "analyzed_at", "url"}),
	}).Create(&cache).Error
}
