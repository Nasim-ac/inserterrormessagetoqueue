package inserterrormessagetoqueue

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // assuming you're using SQL Server
)

func insertErrorMessageToQueue(db *sql.DB, recipients, subject, body string) error {
	// First, check if a message with the same subject and body exists on the same day
	var exists int
	currentDate := time.Now().Format("2006-01-02")

	checkQuery := `
		SELECT COUNT(*) 
		FROM dbo.MessageQueue 
		WHERE CONVERT(VARCHAR, CreationTime, 23) = @CreationDate 
		AND Subject = @Subject 
		AND Body = @Body`

	err := db.QueryRow(checkQuery, sql.Named("CreationDate", currentDate), sql.Named("Subject", subject), sql.Named("Body", body)).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existing message: %v", err)
	}

	// If a message with the same subject and body exists on the same day, return
	if exists > 0 {
		log.Println("Message with the same subject and body already exists today. Skipping insert.")
		return nil
	}

	// If no message exists, insert the new message into the queue
	insertQuery := `
		INSERT INTO dbo.MessageQueue (
			MessageID, CreationTime, Recipients, CC, BCC, Subject, Body, 
			InQueue, Retries, LastTry, LastError, Sent, DatabaseID, MessageType, 
			JobID, EmployeeID, ThreadID, InTransitTimestamp, filename, filetype, HTML
		) 
		VALUES (
			NEWID(), GETDATE(), @Recipients, NULL, NULL, @Subject, @Body, 
			1, 5, NULL, NULL, 0, NULL, NULL, 
			NULL, NULL, NULL, NULL, NULL, 1
		)`

	_, err = db.Exec(insertQuery,
		sql.Named("Recipients", recipients),
		sql.Named("Subject", subject),
		sql.Named("Body", body),
	)

	if err != nil {
		return fmt.Errorf("error inserting message into MessageQueue: %v", err)
	}

	log.Println("Message inserted successfully into MessageQueue.")
	return nil
}
