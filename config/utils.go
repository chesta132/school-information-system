package config

import "school-information-system/internal/libs/errorlib"

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

func EnvCheck() error {
	if IsEnvDev() || IsEnvProd() {
		return nil
	}
	return errorlib.ErrInvalidEnv
}
