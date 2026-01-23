package database

import "github.com/jmoiron/sqlx"

// RunMigrations runs all database migrations
func RunMigrations(db *sqlx.DB) error {
	migrations := []string{
		// Collections table
		`CREATE TABLE IF NOT EXISTS collections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Folders table
		`CREATE TABLE IF NOT EXISTS folders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			collection_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
		)`,

		// Requests table
		`CREATE TABLE IF NOT EXISTS requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			collection_id INTEGER NOT NULL,
			folder_id INTEGER,
			name TEXT NOT NULL,
			method TEXT NOT NULL DEFAULT 'GET',
			url TEXT NOT NULL DEFAULT '',
			headers TEXT DEFAULT '[]',
			params TEXT DEFAULT '[]',
			body TEXT DEFAULT '',
			body_type TEXT DEFAULT 'none',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
			FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
		)`,

		// History table
		`CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id INTEGER,
			method TEXT NOT NULL,
			url TEXT NOT NULL,
			request_headers TEXT,
			request_body TEXT,
			status_code INTEGER,
			response_headers TEXT,
			response_body TEXT,
			duration_ms INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
		)`,

		// Environments table
		`CREATE TABLE IF NOT EXISTS environments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			variables TEXT DEFAULT '[]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Global variables table
		`CREATE TABLE IF NOT EXISTS global_variables (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			variables TEXT DEFAULT '[]',
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// App state table
		`CREATE TABLE IF NOT EXISTS app_state (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			window_width INTEGER DEFAULT 1200,
			window_height INTEGER DEFAULT 800,
			window_x INTEGER,
			window_y INTEGER,
			window_maximized INTEGER DEFAULT 0,
			sidebar_open INTEGER DEFAULT 1,
			sidebar_width INTEGER DEFAULT 260,
			layout_direction TEXT DEFAULT 'horizontal',
			split_ratio INTEGER DEFAULT 50,
			theme TEXT DEFAULT 'system',
			active_env_id INTEGER,
			request_timeout REAL DEFAULT 30,
			auto_locate_sidebar INTEGER DEFAULT 1,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Sidebar state table
		`CREATE TABLE IF NOT EXISTS sidebar_state (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			item_type TEXT NOT NULL,
			item_id INTEGER NOT NULL,
			expanded INTEGER DEFAULT 0,
			UNIQUE(item_type, item_id)
		)`,

		// Tab sessions table
		`CREATE TABLE IF NOT EXISTS tab_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tab_id TEXT UNIQUE NOT NULL,
			request_id INTEGER,
			title TEXT NOT NULL,
			sort_order INTEGER NOT NULL,
			is_active INTEGER DEFAULT 0,
			is_dirty INTEGER DEFAULT 0,
			method TEXT DEFAULT 'GET',
			url TEXT DEFAULT '',
			headers TEXT DEFAULT '[]',
			params TEXT DEFAULT '[]',
			body TEXT DEFAULT '',
			body_type TEXT DEFAULT 'none',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE SET NULL
		)`,

		// Initialize app_state with default values
		`INSERT OR IGNORE INTO app_state (id) VALUES (1)`,

		// Initialize global_variables with default values
		`INSERT OR IGNORE INTO global_variables (id) VALUES (1)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}
