package software

import "dtc-service/core/process"

func InstallSoftware(filepath string) {
	process.InstallProgram(filepath, "/S")
}
