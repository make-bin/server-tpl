package service

// InitServiceBean convert service interface to bean type
func InitServiceBean() []interface{} {
	return []interface{}{
		NewApplicationService(),
		NewVariablesService(),
		NewUserService(),
	}
}
