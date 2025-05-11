package config

import (
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "os"
)

func ConnDB(logLogrus *logrus.Logger) *gorm.DB {
    err := godotenv.Load(".env")
    if err != nil {
        logLogrus.WithFields(logrus.Fields{
            "error":   err,
            "message": "Error loading .env file",
        }).Error("Error loading .env file for database connection")

        return nil
    }

    db, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_USERNAME")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DATABASE")+"?charset=utf8mb4&parseTime=true&loc=Asia%2FJakarta"), &gorm.Config{})
    if err != nil {
        logLogrus.WithFields(logrus.Fields{
            "error":   err,
            "message": "Error connecting to mysql",
        }).Error("Error connecting to mysql")

        return nil
    }

    return db
}
