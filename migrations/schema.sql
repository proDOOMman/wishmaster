CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "version_idx" ON "schema_migration" (version);
CREATE TABLE IF NOT EXISTS "channels_packages" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"active" bool NOT NULL,
"name" TEXT NOT NULL,
"description" text,
"m3u_url" TEXT,
"xmltv_url" TEXT,
"google_cx" TEXT,
"google_key" TEXT,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "channels" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"name" TEXT NOT NULL,
"search_name" TEXT NOT NULL,
"description" text,
"url" TEXT NOT NULL,
"num" integer,
"crypted" bool NOT NULL,
"erotic" bool NOT NULL,
"stream_aspect_ratio" integer NOT NULL,
"zoom_ratio" decimal NOT NULL,
"epg_offset" decimal NOT NULL,
"epg_id" TEXT,
"channels_package_id" int NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
FOREIGN KEY (channels_package_id) REFERENCES channels_packages (id) ON DELETE cascade
);
CREATE TABLE IF NOT EXISTS "epg_channels" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"url_hash" uint32 NOT NULL,
"epg_id" TEXT NOT NULL,
"display_name" TEXT NOT NULL,
"search_name" TEXT NOT NULL,
"icon_src" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "epg_programmes" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"url_hash" uint32 NOT NULL,
"channel_epg_id" TEXT NOT NULL,
"programme_id" TEXT,
"start" integer NOT NULL,
"stop" integer NOT NULL,
"title" TEXT NOT NULL,
"desc" TEXT,
"categories" TEXT,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "service_accounts" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"uuid" TEXT NOT NULL,
"channels_sort" TEXT,
"reminders" TEXT,
"favorite_channels" TEXT,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
