package config

// ---- DO NOT CHANGE ----

func IsEnvProd() bool {
	return GO_ENV == "production"
}

func IsEnvDev() bool {
	return GO_ENV == "development"
}

func SplitByEnv[T any](prodValue, devValue T) T {
	if IsEnvProd() {
		return prodValue
	} else {
		return devValue
	}
}

func init() {
	if IsEnvDev() || IsEnvProd() {
		return
	}
	panic("[ENVIRONMENT] invalid environment, must be 'development' or 'production'")
}
