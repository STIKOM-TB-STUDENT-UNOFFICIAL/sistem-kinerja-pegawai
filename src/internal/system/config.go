package system

import (
	"kinerja-pegawai/internal/utils"
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	godotenv.Load()

	maxprocs, err := strconv.Atoi(os.Getenv("MAXPROCS"))

	if err != nil {
		utils.PanicLog(err.Error())
	}

	runtime.GOMAXPROCS(maxprocs)
}
