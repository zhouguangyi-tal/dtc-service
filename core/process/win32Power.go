package process

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	modAdvapi32                  = syscall.NewLazyDLL("advapi32.dll")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procOpenProcess              = modKernel32.NewProc("OpenProcess")
	procOpenProcessToken         = modAdvapi32.NewProc("OpenProcessToken")
	procDuplicateTokenEx         = modAdvapi32.NewProc("DuplicateTokenEx")
	procCreateProcessAsUser      = modAdvapi32.NewProc("CreateProcessAsUserW")
)

const (
	TH32CS_SNAPPROCESS        = 0x00000002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_CREATE_PROCESS    = 0x0080
	TOKEN_DUPLICATE           = 0x0002
	TOKEN_QUERY               = 0x0008
	TOKEN_ASSIGN_PRIMARY      = 0x0001
	TOKEN_ADJUST_DEFAULT      = 0x0080
	TOKEN_ADJUST_SESSIONID    = 0x0100
	STARTF_USESHOWWINDOW      = 0x00000001
	CREATE_NEW_CONSOLE        = 0x00000010
)

type PROCESSENTRY32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   uintptr
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [260]uint16
}

func getExplorerProcessId() (uint32, error) {
	hSnapshot, _, err := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if hSnapshot == uintptr(syscall.InvalidHandle) {
		return 0, err
	}
	defer syscall.CloseHandle(syscall.Handle(hSnapshot))

	var pe32 PROCESSENTRY32
	pe32.dwSize = uint32(unsafe.Sizeof(pe32))
	ret, _, err := procProcess32First.Call(hSnapshot, uintptr(unsafe.Pointer(&pe32)))
	if ret == 0 {
		return 0, err
	}

	for {
		if syscall.UTF16ToString(pe32.szExeFile[:]) == "explorer.exe" {
			return pe32.th32ProcessID, nil
		}
		ret, _, err = procProcess32Next.Call(hSnapshot, uintptr(unsafe.Pointer(&pe32)))
		if ret == 0 {
			break
		}
	}

	return 0, fmt.Errorf("explorer.exe process not found")
}

func CreateForegroundProcess(applicationName string, args ...string) error {
	explorerPID, err := getExplorerProcessId()
	if err != nil {
		return err
	}

	hProcess, _, err := procOpenProcess.Call(PROCESS_QUERY_INFORMATION, 0, uintptr(explorerPID))
	if hProcess == 0 {
		return err
	}
	defer syscall.CloseHandle(syscall.Handle(hProcess))

	var hToken syscall.Token
	ret, _, err := procOpenProcessToken.Call(hProcess, TOKEN_DUPLICATE|TOKEN_QUERY, uintptr(unsafe.Pointer(&hToken)))
	if ret == 0 {
		return err
	}
	defer syscall.CloseHandle(syscall.Handle(hToken))

	var hNewToken syscall.Token
	ret, _, err = procDuplicateTokenEx.Call(uintptr(hToken), TOKEN_ASSIGN_PRIMARY|TOKEN_DUPLICATE|TOKEN_QUERY|TOKEN_ADJUST_DEFAULT|TOKEN_ADJUST_SESSIONID, 0, 2, 1, uintptr(unsafe.Pointer(&hNewToken)))
	if ret == 0 {
		return err
	}
	defer syscall.CloseHandle(syscall.Handle(hNewToken))

	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))
	si.Flags = STARTF_USESHOWWINDOW
	si.ShowWindow = syscall.SW_SHOW

	commandLine := applicationName
	for _, arg := range args {
		commandLine += " " + arg
	}

	commandLinePtr, err := syscall.UTF16PtrFromString(commandLine)
	if err != nil {
		return err
	}

	ret, _, err = procCreateProcessAsUser.Call(uintptr(hNewToken), 0, uintptr(unsafe.Pointer(commandLinePtr)), 0, 0, 0, CREATE_NEW_CONSOLE, 0, 0, uintptr(unsafe.Pointer(&si)), uintptr(unsafe.Pointer(&pi)))
	if ret == 0 {
		return err
	}

	syscall.CloseHandle(pi.Thread)
	syscall.CloseHandle(pi.Process)

	return nil
}
