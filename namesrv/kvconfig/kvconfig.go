package kvconfig

import (
	"encoding/json"
	"gdiamond/util/fileutil"
	"log"
	"sync"
)

type KVConfig struct {
	sync.RWMutex
	kvConfigPath string
	configTable  map[string] /* Namespace */ map[string] /* Key */ string /* Value */
}

func New(kvConfigPath string) *KVConfig {
	return &KVConfig{kvConfigPath: kvConfigPath, configTable: make(map[string]map[string]string)}
}

//Load get file content of KvConfigPath and load it to configTable
func (kc *KVConfig) Load() {
	content, err := fileutil.GetFileContent(kc.kvConfigPath)
	if err != nil {
		log.Println("Load KV config table failed", err)
		return
	}
	if content != "" {
		kv := make(map[string]map[string]string)
		err = json.Unmarshal([]byte(content), &kv)
		if err != nil {
			log.Println("Unmarshal content to kv config  failed", err)
			return
		}
		for k, v := range kv {
			kc.configTable[k] = v
		}
	}
}

//PutKVConfig put kv to configTable and persist to local file
func (kc *KVConfig) PutKVConfig(namespace, key, value string) {
	kc.Lock()
	kvTable := kc.configTable[namespace]
	if kvTable == nil {
		kvTable = make(map[string]string)
		kc.configTable[namespace] = kvTable
		log.Println("PutKVConfig create new Namespace", namespace)
	}

	prev, _ := kvTable[key]
	if prev != "" {
		log.Println("putKVConfig update config item, Namespace,Key,Value:", namespace, key, value)
	} else {
		log.Println("putKVConfig create new config item, Namespace,Key,Value", namespace, key, value)
	}
	kvTable[key] = value
	kc.Unlock()

	kc.persist()
}

//GetKVListByNamespace  get kvlist by namespace from configTable
//return json.Marshal []byte
func (kc *KVConfig) GetKVListByNamespace(namespace string) []byte {
	kc.RLock()
	defer kc.RUnlock()
	kvTable := kc.configTable[namespace]
	if kvTable != nil {
		b, err := json.Marshal(kvTable)
		if err != nil {
			log.Println("getKVListByNamespace Marshal failed", err)
		}
		return b
	}
	return nil
}

//GetKVConfig  get kvconfig from configTable
func (kc *KVConfig) GetKVConfig(namespace, key string) string {
	kc.Lock()
	defer kc.Unlock()
	kvTable := kc.configTable[namespace]
	if kvTable != nil {
		return kvTable[key]
	}
	return ""
}

//DeleteKVConfig  delete kv from configTable and persist to local file
func (kc *KVConfig) DeleteKVConfig(namespace, key string) {
	kc.Lock()

	kvTable := kc.configTable[namespace]
	if kvTable != nil {
		value := kvTable[key]
		log.Println("deleteKVConfig delete a config item, Namespace: {} Key: {} Value: {}", namespace, key, value)
		delete(kvTable, key)
	}
	kc.Unlock()
	kc.persist()
}

//write configTable to local file
func (kc *KVConfig) persist() {
	kc.RLock()
	defer kc.RUnlock()
	fileName := kc.kvConfigPath
	b, err := json.Marshal(kc.configTable)
	if err != nil {
		log.Println("Marshal kvconfig failed", fileName, err)
		return
	}
	err = fileutil.String2File(string(b), fileName)
	if err != nil {
		log.Println("persist kvconfig failed", fileName, err)
		return
	}
}

func (kc *KVConfig) PrintAllPeriodically() {
	kc.Lock()
	defer kc.Unlock()
	log.Println("--------------------------------------------------------")
	log.Println("configTable SIZE: {}", len(kc.configTable))

	for ns, value := range kc.configTable {
		if value != nil {
			for k, v := range value {
				log.Println("configTable NS: {} Key: {} Value: {}", ns, k, v)
			}
		}
	}
}
