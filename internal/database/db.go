package database

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const (
	appDataDirName   = "postme"
	dataDirName      = "data"
	dbFileName       = "postme.db"
	portableFlagName = "portable.flag"
)

// DB is the global database connection
var DB *sqlx.DB

// Init initializes the database connection
func Init() error {
	// Get data directory path
	dataDir, portable := getDataDirInfo()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	if !portable {
		if err := migrateLegacyDatabase(dataDir); err != nil {
			return err
		}
	}

	dbPath := filepath.Join(dataDir, dbFileName)
	var err error
	DB, err = sqlx.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return err
	}

	// Run migrations
	if err := RunMigrations(DB); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func getDataDirInfo() (string, bool) {
	exeDir, hasExeDir := getExecutableDir()
	if hasExeDir {
		if dataDir, ok := getPortableDataDir(exeDir); ok {
			return dataDir, true
		}
	}

	configDir, err := os.UserConfigDir()
	if err == nil && configDir != "" {
		return filepath.Join(configDir, appDataDirName), false
	}

	if hasExeDir {
		return filepath.Join(exeDir, dataDirName), true
	}

	return dataDirName, true
}

func getExecutableDir() (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	return filepath.Dir(exe), true
}

func getPortableDataDir(exeDir string) (string, bool) {
	flagPath, ok := findPortableFlag(exeDir)
	if !ok {
		return "", false
	}

	flagDir := filepath.Dir(flagPath)
	content, err := os.ReadFile(flagPath)
	if err != nil {
		return filepath.Join(flagDir, dataDirName), true
	}

	dataDir := strings.TrimSpace(string(content))
	if dataDir == "" {
		return filepath.Join(flagDir, dataDirName), true
	}

	return resolveDataDir(flagDir, dataDir), true
}

func findPortableFlag(exeDir string) (string, bool) {
	for _, dir := range portableFlagDirs(exeDir) {
		flagPath := filepath.Join(dir, portableFlagName)
		if info, err := os.Stat(flagPath); err == nil && !info.IsDir() {
			return flagPath, true
		}
	}

	return "", false
}

func portableFlagDirs(exeDir string) []string {
	dirs := []string{exeDir}

	if appRoot, ok := macAppRoot(exeDir); ok {
		appParent := filepath.Dir(appRoot)
		dirs = append(dirs, appRoot, appParent, filepath.Dir(appParent))
	} else {
		dirs = append(dirs, filepath.Dir(exeDir))
	}

	return uniqueDirs(dirs)
}

func macAppRoot(exeDir string) (string, bool) {
	if filepath.Base(exeDir) != "MacOS" {
		return "", false
	}

	contentsDir := filepath.Dir(exeDir)
	if filepath.Base(contentsDir) != "Contents" {
		return "", false
	}

	appRoot := filepath.Dir(contentsDir)
	if filepath.Ext(appRoot) != ".app" {
		return "", false
	}

	return appRoot, true
}

func uniqueDirs(dirs []string) []string {
	seen := make(map[string]bool, len(dirs))
	result := make([]string, 0, len(dirs))

	for _, dir := range dirs {
		if dir == "." || dir == "" {
			continue
		}

		cleanDir := filepath.Clean(dir)
		if seen[cleanDir] {
			continue
		}

		seen[cleanDir] = true
		result = append(result, cleanDir)
	}

	return result
}

func resolveDataDir(baseDir string, dataDir string) string {
	if expanded, ok := expandHomeDir(dataDir); ok {
		return filepath.Clean(expanded)
	}

	if filepath.IsAbs(dataDir) {
		return filepath.Clean(dataDir)
	}

	return filepath.Clean(filepath.Join(baseDir, dataDir))
}

func expandHomeDir(path string) (string, bool) {
	if path != "~" && !strings.HasPrefix(path, "~/") && !strings.HasPrefix(path, `~\`) {
		return "", false
	}

	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return "", false
	}

	if path == "~" {
		return homeDir, true
	}

	return filepath.Join(homeDir, path[2:]), true
}

func migrateLegacyDatabase(dataDir string) error {
	exeDir, ok := getExecutableDir()
	if !ok {
		return nil
	}

	legacyDataDir := filepath.Join(exeDir, dataDirName)
	if dataDir == legacyDataDir {
		return nil
	}

	legacyDBPath := filepath.Join(legacyDataDir, dbFileName)
	newDBPath := filepath.Join(dataDir, dbFileName)

	if _, err := os.Stat(newDBPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(legacyDBPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := copyFile(legacyDBPath, newDBPath); err != nil {
		return err
	}

	for _, suffix := range []string{"-wal", "-shm"} {
		if err := copyFileIfExists(legacyDBPath+suffix, newDBPath+suffix); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}

	return out.Close()
}

func copyFileIfExists(src string, dst string) error {
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return copyFile(src, dst)
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return DB
}
