package handler

// RegisterOrUpdate registers an account if it doesn't exist and
// updates the account information if it exists
func RegisterOrUpdate(w http.ResponseWriter, r *http.Request) {
	var account model.Account

	// get post data from request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	account.Address = strings.ToLower(account.Address)
	account.DeviceType = strings.ToLower(account.DeviceType)

	// check device type
	if account.DeviceType != model.Android && account.DeviceType != model.IOS {
		errors.ErrInvalidDeviceType(w, http.StatusBadRequest)
		return
	}

	// check chain id
	if account.ChainID != model.CosmosHub && account.ChainID != model.IrisHub && account.ChainID != model.Kava {
		errors.ErrInvalidChainID(w, http.StatusBadRequest)
		return
	}

	// [TODO]: check validity of an address depending on which network
	if !strings.Contains(account.Address, sdk.Bech32PrefixAccAddr) || len(account.Address) != 45 {
		errors.ErrInvalidFormat(w, http.StatusBadRequest)
		return
	}

	// insert account information if it doesn't exist
	exist, _ := databases.QueryExistsAccount(w, db, account)
	if !exist {
		result, _, _ := databases.InsertAccount(w, db, account)
		if result != 1 {
			return
		}
		model.Respond(w, true, "successfully inserted")
		return
	}

	result, err := databases.UpdateAccount(w, db, account)
	if !result {
		return
	}

	model.Respond(w, true, "successfully updated")
	return
}

// Delete delete the account information
// CURRENTLY NOT USED
func Delete(w http.ResponseWriter, r *http.Request) {
	var account model.Account

	// get post data from request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		errors.ErrBadRequest(w, http.StatusBadRequest)
		return
	}

	// check if there is the same account
	exist, _ := databases.QueryExistsAccount(w, db, account)
	if !exist {
		errors.ErrNotFound(w, http.StatusNotFound)
		return
	}

	// delete the account
	if exist {
		databases.DeleteAccount(w, db, account)
	}

	model.Respond(w, true, "successfully deleted")
	return
}
