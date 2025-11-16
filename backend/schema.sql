-- RecipeBook Database Schema

-- Users table (stores user profiles from Firebase Auth)
CREATE TABLE IF NOT EXISTS users (
	firebase_uid TEXT PRIMARY KEY,
	email TEXT NOT NULL,
	display_name TEXT,
	role TEXT DEFAULT 'viewer' CHECK(role IN ('viewer', 'editor', 'admin')),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	last_login_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Icons table (shared collection of recipe icons)
CREATE TABLE IF NOT EXISTS icons (
	id TEXT PRIMARY KEY,
	filename TEXT NOT NULL,
	icon_url TEXT NOT NULL,
	uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recipes (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	recipe_type TEXT,
	cuisine TEXT,
	ingredients TEXT,
	method TEXT,
	notes TEXT,
	sources TEXT,
	icon_id TEXT,
	created_by_user_id TEXT,
	created_by_name TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (icon_id) REFERENCES icons(id)
);

CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type);
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
CREATE INDEX IF NOT EXISTS idx_recipes_icon ON recipes(icon_id);

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
	id TEXT PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

-- Many-to-many relationship between recipes and tags
CREATE TABLE IF NOT EXISTS recipe_tags (
	recipe_id TEXT NOT NULL,
	tag_id TEXT NOT NULL,
	PRIMARY KEY (recipe_id, tag_id),
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe ON recipe_tags(recipe_id);
CREATE INDEX IF NOT EXISTS idx_recipe_tags_tag ON recipe_tags(tag_id);

-- Recipe images table
CREATE TABLE IF NOT EXISTS recipe_images (
	id TEXT PRIMARY KEY,
	recipe_id TEXT NOT NULL,
	image_url TEXT NOT NULL,
	display_order INTEGER DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_recipe_images_recipe ON recipe_images(recipe_id);

-- Make logs table (tracks when recipes were made)
CREATE TABLE IF NOT EXISTS make_logs (
	id TEXT PRIMARY KEY,
	recipe_id TEXT NOT NULL,
	made_at DATE NOT NULL,
	notes TEXT,
	created_by_user_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
	FOREIGN KEY (created_by_user_id) REFERENCES users(firebase_uid)
);

CREATE INDEX IF NOT EXISTS idx_make_logs_recipe ON make_logs(recipe_id);
CREATE INDEX IF NOT EXISTS idx_make_logs_made_at ON make_logs(made_at);

-- Full-text search table (now includes description, cuisine, and sources)
-- Note: We include recipe_id as an unindexed column to enable joining back to recipes table
CREATE VIRTUAL TABLE IF NOT EXISTS recipes_fts USING fts5(
	recipe_id UNINDEXED,
	title, description, cuisine, ingredients, method, notes, sources
);

-- Triggers to keep FTS in sync
CREATE TRIGGER IF NOT EXISTS recipes_fts_insert AFTER INSERT ON recipes BEGIN
	INSERT INTO recipes_fts(recipe_id, title, description, cuisine, ingredients, method, notes, sources)
	VALUES (new.id, new.title, new.description, new.cuisine, new.ingredients, new.method, new.notes, new.sources);
END;

CREATE TRIGGER IF NOT EXISTS recipes_fts_update AFTER UPDATE ON recipes BEGIN
	UPDATE recipes_fts
	SET title=new.title, description=new.description, cuisine=new.cuisine,
		ingredients=new.ingredients, method=new.method, notes=new.notes, sources=new.sources
	WHERE recipe_id=new.id;
END;

CREATE TRIGGER IF NOT EXISTS recipes_fts_delete AFTER DELETE ON recipes BEGIN
	DELETE FROM recipes_fts WHERE recipe_id=old.id;
END;
