package cluster

// DO NOT EDIT BY HAND
//
// This code was generated by the schema.DotGo function. If you need to
// modify the database schema, please add a new schema update to update.go
// and the run 'make update-schema'.
const freshSchema = `
CREATE TABLE certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    fingerprint TEXT NOT NULL,
    type INTEGER NOT NULL,
    name TEXT NOT NULL,
    certificate TEXT NOT NULL,
    UNIQUE (fingerprint)
);
CREATE TABLE config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (key)
);
CREATE TABLE containers (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    node_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    architecture INTEGER NOT NULL,
    type INTEGER NOT NULL,
    ephemeral INTEGER NOT NULL DEFAULT 0,
    creation_date DATETIME NOT NULL DEFAULT 0,
    stateful INTEGER NOT NULL DEFAULT 0,
    last_use_date DATETIME,
    description TEXT,
    project_id INTEGER NOT NULL DEFAULT 0 REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE (name),
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE containers_backups (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    container_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    creation_date DATETIME,
    expiry_date DATETIME,
    container_only INTEGER NOT NULL default 0,
    optimized_storage INTEGER NOT NULL default 0,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE,
    UNIQUE (container_id, name)
);
CREATE TABLE containers_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    container_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE,
    UNIQUE (container_id, key)
);
CREATE TABLE containers_devices (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    container_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL default 0,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE,
    UNIQUE (container_id, name)
);
CREATE TABLE containers_devices_config (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    container_device_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    FOREIGN KEY (container_device_id) REFERENCES containers_devices (id) ON DELETE CASCADE,
    UNIQUE (container_device_id, key)
);
CREATE TABLE containers_profiles (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    container_id INTEGER NOT NULL,
    profile_id INTEGER NOT NULL,
    apply_order INTEGER NOT NULL default 0,
    UNIQUE (container_id, profile_id),
    FOREIGN KEY (container_id) REFERENCES containers(id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
CREATE TABLE images (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    fingerprint TEXT NOT NULL,
    filename TEXT NOT NULL,
    size INTEGER NOT NULL,
    public INTEGER NOT NULL DEFAULT 0,
    architecture INTEGER NOT NULL,
    creation_date DATETIME,
    expiry_date DATETIME,
    upload_date DATETIME NOT NULL,
    cached INTEGER NOT NULL DEFAULT 0,
    last_use_date DATETIME,
    auto_update INTEGER NOT NULL DEFAULT 0,
    project_id INTEGER NOT NULL DEFAULT 0 REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE (fingerprint)
);
CREATE TABLE images_aliases (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    image_id INTEGER NOT NULL,
    description TEXT,
    project_id INTEGER NOT NULL DEFAULT 0 REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE,
    UNIQUE (name)
);
CREATE TABLE images_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    UNIQUE (image_id, node_id),
    FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE images_properties (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE
);
CREATE TABLE images_source (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    server TEXT NOT NULL,
    protocol INTEGER NOT NULL,
    certificate TEXT NOT NULL,
    alias TEXT NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE
);
CREATE TABLE networks (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    state INTEGER NOT NULL DEFAULT 0,
    UNIQUE (name)
);
CREATE TABLE networks_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    network_id INTEGER NOT NULL,
    node_id INTEGER,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (network_id, node_id, key),
    FOREIGN KEY (network_id) REFERENCES networks (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE networks_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    network_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    UNIQUE (network_id, node_id),
    FOREIGN KEY (network_id) REFERENCES networks (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE nodes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    address TEXT NOT NULL,
    schema INTEGER NOT NULL,
    api_extensions INTEGER NOT NULL,
    heartbeat DATETIME DEFAULT CURRENT_TIMESTAMP,
    pending INTEGER NOT NULL DEFAULT 0,
    UNIQUE (name),
    UNIQUE (address)
);
CREATE TABLE operations (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    uuid TEXT NOT NULL,
    node_id TEXT NOT NULL,
    type INTEGER NOT NULL DEFAULT 0,
    UNIQUE (uuid),
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    project_id INTEGER NOT NULL DEFAULT 0 REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE (name)
);
CREATE TABLE profiles_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (profile_id, key),
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
CREATE TABLE profiles_devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL default 0,
    UNIQUE (profile_id, name),
    FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);
CREATE TABLE profiles_devices_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_device_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (profile_device_id, key),
    FOREIGN KEY (profile_device_id) REFERENCES profiles_devices (id) ON DELETE CASCADE
);
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name VARCHAR(255) NOT NULL,
    has_images INTEGER NOT NULL default 0,
    has_profiles INTEGER NOT NULL default 0,
    description TEXT,
    UNIQUE (name)
);
CREATE TABLE storage_pools (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    driver TEXT NOT NULL,
    description TEXT,
    state INTEGER NOT NULL DEFAULT 0,
    UNIQUE (name)
);
CREATE TABLE storage_pools_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (storage_pool_id, node_id, key),
    FOREIGN KEY (storage_pool_id) REFERENCES storage_pools (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE storage_pools_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    UNIQUE (storage_pool_id, node_id),
    FOREIGN KEY (storage_pool_id) REFERENCES storage_pools (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE storage_volumes (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    description TEXT,
    snapshot INTEGER NOT NULL DEFAULT 0,
    UNIQUE (storage_pool_id, node_id, name, type),
    FOREIGN KEY (storage_pool_id) REFERENCES storage_pools (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
);
CREATE TABLE storage_volumes_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_volume_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (storage_volume_id, key),
    FOREIGN KEY (storage_volume_id) REFERENCES storage_volumes (id) ON DELETE CASCADE
);

INSERT INTO schema (version, updated_at) VALUES (11, strftime("%s"))
`
