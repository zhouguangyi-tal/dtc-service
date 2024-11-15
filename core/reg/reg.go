package reg

import (
	"golang.org/x/sys/windows/registry"
	"log"
)

type Registry struct {
	path   string
	talReg map[string]map[string]string
}

func (r *Registry) Init(path string) {
	r.path = path
	r.talReg = make(map[string]map[string]string)
	r.readRegistry()
}

func (r *Registry) readSubKey() { //读取 注册表的子项
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, r.path, registry.READ)
	if err != nil {
		log.Fatalln("Error opening registry key:", err)
		return
	}
	defer key.Close()
	subKeyNames, err := key.ReadSubKeyNames(-1)
	for _, subKey := range subKeyNames {
		r.talReg[subKey] = make(map[string]string)
		r.readKeyValue(subKey)
	}
}
func (r *Registry) readKeyValue(keyName string) { //读取注册表的值
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, r.path+"\\"+keyName, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalln("Error opening registry value:", err)
		return
	}
	defer key.Close()

	names, err := key.ReadValueNames(-1)
	for _, name := range names {
		val, _, err := key.GetStringValue(name)
		if err != nil {
			log.Fatalln("Error reading value:", err)
		}
		r.talReg[keyName][name] = val
	}

}

func (r *Registry) readRegistry() {
	r.readSubKey()
}

func (r *Registry) GetRegistryInfo(app AppName) RegeditInfo {
	info := RegeditInfo{}
	info.ExePath = r.talReg[app]["exe path"]
	info.InstallPath = r.talReg[app]["install path"]
	info.InstallDate = r.talReg[app]["install date"]
	info.Version = r.talReg[app]["version"]

	return info
}

var Reg = Registry{}
