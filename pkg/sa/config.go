package sa

import (
	"os"
	"strconv"
)

var pdbName = loadEnv("PDB_NAME", "products")
var pdbUser = loadEnv("PDB_USER", "postgres")
var pdbPass = loadEnv("PDB_PASS", "postgres")
var pdbHost = loadEnv("PDB_HOST", "localhost")

var RdbHostPort = loadEnv("RDB_HOST_PORT", "localhost:6379")
var LogLevel = loadEnv("LOG_LEVEL", "debug")
var doTrigger = loadBoolEnv("ENABLE_TRIGGERS", false)

func loadEnv(env, def string) string {
	if val, ok := os.LookupEnv(env); ok {
		return val
	}
	return def
}

func loadBoolEnv(env string, def bool) bool {
	if val, ok := os.LookupEnv(env); ok {
		ret, err := strconv.ParseBool(val)
		if err != nil {
			panic(err)
		}
		return ret == true
	}
	return def
}
