package handler

// GetVersion returns version number of an app
func GetVersion(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceType := vars["deviceType"]

	// lower case
	deviceType = strings.ToLower(deviceType)

	var version schema.MobileVersion

	switch deviceType {
	case model.Android:
		version, _ = s.db.QueryMobileVersion(w, db, deviceType)
	case model.IOS:
		version, _ = s.db.QueryMobileVersion(w, db, deviceType)
	default:
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	// in case when data is empty
	if version.Version == 0 {
		return
	}

	version.IdfVersion = 0

	u.Respond(w, version)
	return
}

// SetVersion sets version number of an app
func SetVersion(db *pg.DB, w http.ResponseWriter, r *http.Request) {
	var version model.AppVersion

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&version)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	// lower case
	version.AppName = strings.ToLower(version.AppName)
	version.DeviceType = strings.ToLower(version.DeviceType)

	exist, _ := databases.QueryExistsAppVersion(w, db, version)

	// update version info if it exists. Otherwise, insert version info
	if exist {
		databases.UpdateAppVersion(w, db, version)
	} else {
		databases.InsertAppVersion(w, db, version)
	}

	model.Respond(w, version)
	return
}
