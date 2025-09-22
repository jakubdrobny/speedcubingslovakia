package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"syscall"
	"time"

	"github.com/araddon/dateparse"
	"github.com/joho/godotenv"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var sourceDirs = []string{"/app/grafana_data", "/app/logs", "/app/loki_data", "/app/mimir_data"}

// return filename, stdout, stderr, err
func CompressMonitoringData(envMap map[string]string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("%w: when creating outputFile=%s", err, outputFile)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, sourceDir := range sourceDirs {
		err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return fmt.Errorf("%w: when creating tar header for %s", err, path)
			}

			relPath, err := filepath.Rel(filepath.Dir(sourceDir), path)
			if err != nil {
				return fmt.Errorf("%w: when getting relative path for %s", err, path)
			}
			header.Name = relPath

			if err := tarWriter.WriteHeader(header); err != nil {
				return fmt.Errorf("%w: when writing tar header for %s", err, path)
			}

			if !info.IsDir() {
				fileToTar, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("%w: when opening file %s", err, path)
				}

				if _, err := io.Copy(tarWriter, fileToTar); err != nil {
					return fmt.Errorf("%w: when copying file content for %s", err, path)
				}

				if err = fileToTar.Close(); err != nil {
					return fmt.Errorf("%w: when closing file=%s", err, path)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("%w: when walking directory %s", err, sourceDir)
		}
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
		return "", fmt.Errorf("%w: when uploading file with name=%s to drive", err, file.Name())
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
		return fmt.Errorf("%w: when reading dir=%s", err, folderPath)
	}

	// remove from local
	backupsFileEntries := make([]BackupFileEntry, 0)

	for _, e := range entries {
		filePath := filepath.Join(folderPath, e.Name())
		fi, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("%w: when stating file with path=%s", err, filePath)
		}

		stat := fi.Sys().(*syscall.Stat_t)
		ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))

		backupsFileEntries = append(backupsFileEntries, BackupFileEntry{Path: filePath, CTime: ctime})
	}

	if len(backupsFileEntries) <= 10 {
		log.Println("Not enough file entries to erase an extra one.")
		return nil
	}

	sort.Slice(backupsFileEntries, func(i, j int) bool { return backupsFileEntries[i].CTime.Before(backupsFileEntries[j].CTime) })

	for i := range len(backupsFileEntries) - 10 {
		err = os.Remove(backupsFileEntries[i].Path)
		if err != nil {
			return fmt.Errorf("%w: when removing file with path=%s", err, backupsFileEntries[i].Path)
		}
	}

	// remove from gdrive
	backupsFileEntries = make([]BackupFileEntry, 0)

	f, err := fileService.List().Fields("files(id,name,createdTime,parents)").Do()
	if err != nil {
		return fmt.Errorf("%w: when listing files in drive backup folder", err)
	}

	for _, e := range f.Files {
		isInFolder := slices.Contains(e.Parents, driveBackupFolderId)

		if !isInFolder {
			continue
		}

		parsedTime, err := dateparse.ParseAny(e.CreatedTime)
		if err != nil {
			return fmt.Errorf("%w: when parsing create time of backup", err)
		}

		backupsFileEntries = append(backupsFileEntries, BackupFileEntry{Id: e.Id, CTime: parsedTime})
	}

	sort.Slice(backupsFileEntries, func(i, j int) bool { return backupsFileEntries[i].CTime.Before(backupsFileEntries[j].CTime) })
	for i := range len(backupsFileEntries) - 10 {
		err := fileService.Delete(backupsFileEntries[i].Id).SupportsAllDrives(true).Do()
		if err != nil {
			return fmt.Errorf("%w: when deleting file with id=%s from drive", err, backupsFileEntries[i].Id)
		}
	}

	return nil
}

func main() {
	log.Println("Starting monitoring backup procedure...")
	log.Println("Loading environment variables...")
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatalf("Unable to load environmental variables from file: %v\n", err)
		return
	}

	log.Println("Environment variables successfully loaded.")

	log.Println("Creating new drive service...")
	credsFilePath := "/app/configs/drive-credentials.json"
	service, err := drive.NewService(context.Background(), option.WithCredentialsFile(credsFilePath))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return
	}

	log.Println("New drive service successfully created.")
	log.Println("Creating compressed backup of monitoring...")

	filename := envMap["MONITORING_BACKUPS_FOLDER_PATH"] + "/" + time.Now().Format("2006-01-02_15-04-05") + ".tar.gz"
	err = CompressMonitoringData(envMap, filename)
	if err != nil {
		log.Println("ERR in DumpDatabase: " + err.Error())
		return
	}

	log.Println("Compressed monitoring backup successfully created.")

	log.Println("Opening backup...")
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open backup.")
		return
	}

	log.Println("Backup successfully opened.")
	log.Println("Uploading backup to Google Drive...")

	_, err = UploadFile(service, f, envMap["DRIVE_MONITORING_BACKUP_FOLDER_ID"])
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when uploading file", err))
		return
	}

	log.Println("Backup uploaded successfully to Google Drive.")
	log.Println("Closing backup in os...")

	err = f.Close()
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when closing backup file", err))
		return
	}

	log.Println("Backup in os successfully closed.")
	log.Printf("Removing oldest backup in %s... (if more than 10 dump files are stored)\n", envMap["MONITORING_BACKUPS_FOLDER_PATH"])

	err = RemoveOldestBackups(envMap["MONITORING_BACKUPS_FOLDER_PATH"], envMap["DRIVE_MONITORING_BACKUP_FOLDER_ID"], service.Files)
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when removing oldest backup from drive folder", err))
		return
	}

	log.Printf("Oldest backup in %s successfully removed. (if more than 10 dump files are stored)\n", envMap["MONITORING_BACKUPS_FOLDER_PATH"])
	log.Println("Monitoring backup procedure successfully finished.")
}
