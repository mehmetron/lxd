//go:build linux && cgo && !agent
// +build linux,cgo,!agent

package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"database/sql"
	"fmt"

	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

var _ = api.ServerEnvironment{}

var profileObjects = cluster.RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles JOIN projects ON profiles.project_id = projects.id
  ORDER BY projects.id, profiles.name
`)

var profileObjectsByID = cluster.RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles JOIN projects ON profiles.project_id = projects.id
  WHERE profiles.id = ? ORDER BY projects.id, profiles.name
`)

var profileObjectsByProject = cluster.RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles JOIN projects ON profiles.project_id = projects.id
  WHERE project = ? ORDER BY projects.id, profiles.name
`)

var profileObjectsByProjectAndName = cluster.RegisterStmt(`
SELECT profiles.id, profiles.project_id, projects.name AS project, profiles.name, coalesce(profiles.description, '')
  FROM profiles JOIN projects ON profiles.project_id = projects.id
  WHERE project = ? AND profiles.name = ? ORDER BY projects.id, profiles.name
`)

var profileID = cluster.RegisterStmt(`
SELECT profiles.id FROM profiles JOIN projects ON profiles.project_id = projects.id
  WHERE projects.name = ? AND profiles.name = ?
`)

var profileCreate = cluster.RegisterStmt(`
INSERT INTO profiles (project_id, name, description)
  VALUES ((SELECT projects.id FROM projects WHERE projects.name = ?), ?, ?)
`)

var profileRename = cluster.RegisterStmt(`
UPDATE profiles SET name = ? WHERE project_id = (SELECT projects.id FROM projects WHERE projects.name = ?) AND name = ?
`)

var profileDeleteByProjectAndName = cluster.RegisterStmt(`
DELETE FROM profiles WHERE project_id = (SELECT projects.id FROM projects WHERE projects.name = ?) AND name = ?
`)

var profileUpdate = cluster.RegisterStmt(`
UPDATE profiles
  SET project_id = (SELECT id FROM projects WHERE name = ?), name = ?, description = ?
 WHERE id = ?
`)

// GetProfileURIs returns all available profile URIs.
// generator: profile URIs
func (c *ClusterTx) GetProfileURIs(filter ProfileFilter) ([]string, error) {
	var err error

	// Result slice.
	objects := make([]Profile, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if filter.Project != nil && filter.Name != nil && filter.ID == nil {
		stmt = c.stmt(profileObjectsByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if filter.Project != nil && filter.ID == nil && filter.Name == nil {
		stmt = c.stmt(profileObjectsByProject)
		args = []interface{}{
			filter.Project,
		}
	} else if filter.ID != nil && filter.Project == nil && filter.Name == nil {
		stmt = c.stmt(profileObjectsByID)
		args = []interface{}{
			filter.ID,
		}
	} else if filter.ID == nil && filter.Project == nil && filter.Name == nil {
		stmt = c.stmt(profileObjects)
		args = []interface{}{}
	} else {
		return nil, fmt.Errorf("No statement exists for the given Filter")
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Profile{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].ProjectID,
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Description,
		}
	}

	// Select.
	err = query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	uris := make([]string, len(objects))
	for i := range objects {
		uri := api.NewURL().Path(version.APIVersion, "profiles", objects[i].Name)
		uri.Project(objects[i].Project)

		uris[i] = uri.String()
	}

	return uris, nil
}

// GetProfiles returns all available profiles.
// generator: profile GetMany
func (c *ClusterTx) GetProfiles(filter ProfileFilter) ([]Profile, error) {
	var err error

	// Result slice.
	objects := make([]Profile, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if filter.Project != nil && filter.Name != nil && filter.ID == nil {
		stmt = c.stmt(profileObjectsByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if filter.Project != nil && filter.ID == nil && filter.Name == nil {
		stmt = c.stmt(profileObjectsByProject)
		args = []interface{}{
			filter.Project,
		}
	} else if filter.ID != nil && filter.Project == nil && filter.Name == nil {
		stmt = c.stmt(profileObjectsByID)
		args = []interface{}{
			filter.ID,
		}
	} else if filter.ID == nil && filter.Project == nil && filter.Name == nil {
		stmt = c.stmt(profileObjects)
		args = []interface{}{}
	} else {
		return nil, fmt.Errorf("No statement exists for the given Filter")
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Profile{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].ProjectID,
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Description,
		}
	}

	// Select.
	err = query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	config, err := c.GetConfig("profile")
	if err != nil {
		return nil, err
	}

	for i := range objects {
		if _, ok := config[objects[i].ID]; !ok {
			objects[i].Config = map[string]string{}
		} else {
			objects[i].Config = config[objects[i].ID]
		}
	}

	devices, err := c.GetDevices("profile")
	if err != nil {
		return nil, err
	}

	for i := range objects {
		objects[i].Devices = map[string]Device{}
		for _, obj := range devices[objects[i].ID] {
			if _, ok := objects[i].Devices[obj.Name]; !ok {
				objects[i].Devices[obj.Name] = obj
			} else {
				return nil, fmt.Errorf("Found duplicate Device with name %q", obj.Name)
			}
		}
	}

	// Use non-generated custom method for UsedBy fields.
	for i := range objects {
		usedBy, err := c.GetProfileUsedBy(objects[i])
		if err != nil {
			return nil, err
		}

		objects[i].UsedBy = usedBy
	}

	return objects, nil
}

// GetProfile returns the profile with the given key.
// generator: profile GetOne
func (c *ClusterTx) GetProfile(project string, name string) (*Profile, error) {
	filter := ProfileFilter{}
	filter.Project = &project
	filter.Name = &name

	objects, err := c.GetProfiles(filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"profiles\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, ErrNoSuchObject
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"profiles\" entry matches")
	}
}

// ProfileExists checks if a profile with the given key exists.
// generator: profile Exists
func (c *ClusterTx) ProfileExists(project string, name string) (bool, error) {
	_, err := c.GetProfileID(project, name)
	if err != nil {
		if err == ErrNoSuchObject {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetProfileID return the ID of the profile with the given key.
// generator: profile ID
func (c *ClusterTx) GetProfileID(project string, name string) (int64, error) {
	stmt := c.stmt(profileID)
	rows, err := stmt.Query(project, name)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"profiles\" ID: %w", err)
	}

	defer rows.Close()

	// Ensure we read one and only one row.
	if !rows.Next() {
		return -1, ErrNoSuchObject
	}
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("Failed to scan ID: %w", err)
	}

	if rows.Next() {
		return -1, fmt.Errorf("More than one row returned")
	}
	err = rows.Err()
	if err != nil {
		return -1, fmt.Errorf("Result set failure: %w", err)
	}

	return id, nil
}

// CreateProfile adds a new profile to the database.
// generator: profile Create
func (c *ClusterTx) CreateProfile(object Profile) (int64, error) {
	// Check if a profile with the same key exists.
	exists, err := c.ProfileExists(object.Project, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, fmt.Errorf("This \"profiles\" entry already exists")
	}

	args := make([]interface{}, 3)

	// Populate the statement arguments.
	args[0] = object.Project
	args[1] = object.Name
	args[2] = object.Description

	// Prepared statement to use.
	stmt := c.stmt(profileCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"profiles\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"profiles\" entry ID: %w", err)
	}

	referenceID := int(id)
	for key, value := range object.Config {
		insert := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err = c.CreateConfig("profile", insert)
		if err != nil {
			return -1, fmt.Errorf("Insert Config failed for Profile: %w", err)
		}

	}
	for _, insert := range object.Devices {
		insert.ReferenceID = int(id)
		err = c.CreateDevice("profile", insert)
		if err != nil {
			return -1, fmt.Errorf("Insert Devices failed for Profile: %w", err)
		}

	}
	return id, nil
}

// RenameProfile renames the profile matching the given key parameters.
// generator: profile Rename
func (c *ClusterTx) RenameProfile(project string, name string, to string) error {
	stmt := c.stmt(profileRename)
	result, err := stmt.Exec(to, project, name)
	if err != nil {
		return fmt.Errorf("Rename Profile failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}
	return nil
}

// DeleteProfile deletes the profile matching the given key parameters.
// generator: profile DeleteOne-by-Project-and-Name
func (c *ClusterTx) DeleteProfile(project string, name string) error {
	stmt := c.stmt(profileDeleteByProjectAndName)
	result, err := stmt.Exec(project, name)
	if err != nil {
		return fmt.Errorf("Delete \"profiles\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query deleted %d rows instead of 1", n)
	}

	return nil
}

// UpdateProfile updates the profile matching the given key parameters.
// generator: profile Update
func (c *ClusterTx) UpdateProfile(project string, name string, object Profile) error {
	id, err := c.GetProfileID(project, name)
	if err != nil {
		return err
	}

	stmt := c.stmt(profileUpdate)
	result, err := stmt.Exec(object.Project, object.Name, object.Description, id)
	if err != nil {
		return fmt.Errorf("Update \"profiles\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	err = c.UpdateConfig("profile", int(id), object.Config)
	if err != nil {
		return fmt.Errorf("Replace Config for Profile failed: %w", err)
	}

	err = c.UpdateDevice("profile", int(id), object.Devices)
	if err != nil {
		return fmt.Errorf("Replace Devices for Profile failed: %w", err)
	}

	return nil
}
