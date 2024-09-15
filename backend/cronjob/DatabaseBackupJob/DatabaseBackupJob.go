package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/araddon/dateparse"
	"github.com/joho/godotenv"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// return filename, stdout, stderr, err
func DumpDatabase(envMap map[string]string, filename string) error {
	cmd := exec.Command("pg_dump", envMap["PG_DUMP_CONNECTION_STRING"], "-f", filename)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func UploadFile(service *drive.Service, file *os.File, parentFolders ...string) (string, error) {
	driveFile := &drive.File{
		Name:    filepath.Base(file.Name()),
		Parents: parentFolders,
	}
	fileSrv := service.Files
	call := fileSrv.Create(driveFile).Media(file)
	createdFile, err := call.Do()
	if err != nil {
		return "", err
	}
	return createdFile.Id, nil
}

type BackupFileEntry struct {
	Path  string
	Id    string
	CTime time.Time
}

func RemoveOldestBackups(folderPath, driveBackupFolderId string, fileService *drive.FilesService) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// remove from local
	backupsFileEntries := make([]BackupFileEntry, 0)

	for _, e := range entries {
		filePath := folderPath + "/" + e.Name()
		fi, err := os.Stat(filePath)
		if err != nil {
			return nil
		}

		stat := fi.Sys().(*syscall.Stat_t)
		ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))

		backupsFileEntries = append(backupsFileEntries, BackupFileEntry{Path: filePath, CTime: ctime})
	}

	if len(backupsFileEntries) <= 10 {
		log.Println("Not enough file entries to erase extra one.")
		return nil
	}

	sort.Slice(backupsFileEntries, func(i, j int) bool { return backupsFileEntries[i].CTime.Before(backupsFileEntries[j].CTime) })

	for i := range len(backupsFileEntries) - 10 {
		err = os.Remove(backupsFileEntries[i].Path)
		if err != nil {
			return err
		}
	}

	// remove from gdrive
	backupsFileEntries = make([]BackupFileEntry, 0)

	f, err := fileService.List().Fields("files(id,name,createdTime,parents)").Do()
	if err != nil {
		return err
	}

	for _, e := range f.Files {
		isInFolder := false

		for _, parent := range e.Parents {
			if parent == driveBackupFolderId {
				isInFolder = true
				break
			}
		}

		if !isInFolder {
			continue
		}

		parsedTime, err := dateparse.ParseAny(e.CreatedTime)
		if err != nil {
			return err
		}

		backupsFileEntries = append(backupsFileEntries, BackupFileEntry{Id: e.Id, CTime: parsedTime})
	}

	sort.Slice(backupsFileEntries, func(i, j int) bool { return backupsFileEntries[i].CTime.Before(backupsFileEntries[j].CTime) })
	for i := range len(backupsFileEntries) - 10 {
		err := fileService.Delete(backupsFileEntries[i].Id).SupportsAllDrives(true).Do()
		if err != nil {
			return nil
		}
	}

	return nil
}

func main() {
	log.Println("Starting database backup procedure...")
	log.Println("Loading environment variables...")
	envMap, err := godotenv.Read(fmt.Sprintf(".env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")))
	if err != nil {
		log.Printf("Unable to load enviromental variables from file: %v\n", err)
		return
	}

	log.Println("Environment variables successfully loaded.")

	log.Println("Creating new drive service...")
	credsFilePath := fmt.Sprintf("drive-credentials-%s.json", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV"))
	service, err := drive.NewService(context.Background(), option.WithCredentialsFile(credsFilePath))
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
		return
	}

	log.Println("New drive service successfully created.")
	log.Println("Dumping database into file...")

	filename := envMap["DB_BACKUPS_FOLDER_PATH"] + "/" + time.Now().Format("2006-01-02_15-04-05") + ".sql"
	err = DumpDatabase(envMap, filename)
	if err != nil {
		log.Println("ERR in DumpDatabase: " + err.Error())
		return
	}

	log.Println("Database dump file successfully created.")

	log.Println("Opening dump file...")
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Failed to open file.")
		return
	}

	log.Println("Dump file successfully opened.")
	log.Println("Uploading dump file to Google Drive...")

	_, err = UploadFile(service, f, envMap["DRIVE_BACKUP_FOLDER_ID"])
	if err != nil {
		log.Println("ERR in UploadFile: " + err.Error())
		return
	}

	log.Println("Dump file uploaded successfully to Google Drive.")
	log.Println("Closing dump file in os...")

	err = f.Close()
	if err != nil {
		log.Println("ERR in f.Close(): " + err.Error())
		return
	}

	log.Println("Dump file in os successfully closed.")
	log.Println("Removing oldest backup dump in $DB_BACKUPS_FOLDER_PATH... (if more than 10 dump files are stored)")

	err = RemoveOldestBackups(envMap["DB_BACKUPS_FOLDER_PATH"], envMap["DRIVE_BACKUP_FOLDER_ID"], service.Files)
	if err != nil {
		log.Println("ERR in RemoveOldestBackups: " + err.Error())
		return
	}

	log.Println("Oldest backup dump in $DB_BACKUPS_FOLDER_PATH successfully removed. (if more than 10 dump files are stored)")
	log.Println("Database backup procedure successfully finished.")
}
