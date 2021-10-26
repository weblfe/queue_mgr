package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/utils"
	"os"
	"sort"
	"strings"
)

type (
	LocalStorage struct {
		DbFile  string  `json:"db_file" yaml:"dbFile"`
		Options Options `json:"options" yaml:"options"`
		Default bool    `json:"default,default=false" yaml:"default"`
	}

	Options struct {

		// The default value is 8MiB.
		BlockCacheCapacity string `json:"block_cache_capacity" yaml:"block_cache_capacity"`

		// The default if false.
		BlockCacheEvictRemoved bool `json:"block_cache_evict_removed" yaml:"block_cache_evict_removed"`

		// BlockRestartInterval is the number of keys between restart points for
		// delta encoding of keys.
		//
		// The default value is 16.
		BlockRestartInterval int `json:"block_restart_interval" yaml:"block_restart_interval"`

		// BlockSize is the minimum uncompressed size in bytes of each 'sorted table'
		// block.
		//
		// The default value is 4 KiB.
		BlockSize string `json:"block_size" yaml:"block_size"`

		// CompactionExpandLimitFactor limits compaction size after expanded.
		// This will be multiplied by table size limit at compaction target level.
		//
		// The default value is 25.
		CompactionExpandLimitFactor int `json:"compaction_expand_limit_factor" yaml:"compaction_expand_limit_factor"`

		// CompactionGPOverlapsFactor limits overlaps in grandparent (Level + 2) that a
		// single 'sorted table' generates.
		// This will be multiplied by table size limit at grandparent level.
		//
		// The default value is 10.
		CompactionGPOverlapsFactor int `json:"compaction_gp_overlaps_factor" yaml:"compaction_gp_overlaps_factor"`

		// CompactionL0Trigger defines number of 'sorted table' at level-0 that will
		// trigger compaction.
		//
		// The default value is 4.
		CompactionL0Trigger int `json:"compaction_l_0_trigger" yaml:"compaction_l_0_trigger"`

		// CompactionSourceLimitFactor limits compaction source size. This doesn't apply to
		// level-0.
		// This will be multiplied by table size limit at compaction target level.
		//
		// The default value is 1.
		CompactionSourceLimitFactor int `json:"compaction_source_limit_factor" yaml:"compaction_source_limit_factor"`

		CompactionTableSize string `json:"compaction_table_size" yaml:"compaction_table_size"`

		// CompactionTableSizeMultiplier defines multiplier for CompactionTableSize.
		//
		// The default value is 1.
		CompactionTableSizeMultiplier float64 `json:"compaction_table_size_multiplier" yaml:"compaction_table_size_multiplier"`

		// CompactionTableSizeMultiplierPerLevel defines per-level multiplier for
		// CompactionTableSize.
		// Use zero to skip a level.
		//
		// The default value is nil.
		CompactionTableSizeMultiplierPerLevel []float64 `json:"compaction_table_size_multiplier_per_level" yaml:"compaction_table_size_multiplier_per_level"`

		// CompactionTotalSize limits total size of 'sorted table' for each level.
		// The limits for each level will be calculated as:
		//   CompactionTotalSize * (CompactionTotalSizeMultiplier ^ Level)
		// The multiplier for each level can also fine-tuned using
		// CompactionTotalSizeMultiplierPerLevel.
		//
		// The default value is 10 MiB.
		CompactionTotalSize string `json:"compaction_total_size" yaml:"compaction_total_size"`

		// CompactionTotalSizeMultiplier defines multiplier for CompactionTotalSize.
		//
		// The default value is 10.
		CompactionTotalSizeMultiplier float64 `json:"compaction_total_size_multiplier" yaml:"compaction_total_size_multiplier"`

		// CompactionTotalSizeMultiplierPerLevel defines per-level multiplier for
		// CompactionTotalSize.
		// Use zero to skip a level.
		//
		// The default value is nil.
		CompactionTotalSizeMultiplierPerLevel []float64 `json:"compaction_total_size_multiplier_per_level" yaml:"compaction_total_size_multiplier_per_level"`

		// Compression defines the 'sorted table' block compression to use.
		//
		// The default value (DefaultCompression) uses snappy compression.
		Compression string `json:"compression" yaml:"compression"`

		// The default value is false.
		DisableBufferPool bool `json:"disable_buffer_pool" yaml:"disable_buffer_pool"`

		// The default value is false.
		DisableBlockCache bool `json:"disable_block_cache" yaml:"disable_block_cache"`

		// The default value is false.
		DisableCompactionBackoff bool `json:"disable_compaction_backoff" yaml:"disable_compaction_backoff"`

		// The default is false.
		DisableLargeBatchTransaction bool `json:"disable_large_batch_transaction" yaml:"disable_large_batch_transaction"`

		// The default value is false.
		ErrorIfExist bool `json:"error_if_exist" yaml:"error_if_exist"`

		// The default value is false.
		ErrorIfMissing bool `json:"error_if_missing" yaml:"error_if_missing"`

		// Filter defines an 'effective filter' to use. An 'effective filter'
		// if defined will be used to generate per-table filter block.
		// The filter name will be stored on disk.
		// During reads LevelDB will try to find matching filter from
		// 'effective filter' and 'alternative filters'.
		//
		// Filter can be changed after a DB has been created. It is recommended
		// to put old filter to the 'alternative filters' to mitigate lack of
		// filter during transition period.
		//
		// A filter is used to reduce disk reads when looking for a specific key.
		//
		// IteratorSamplingRate defines approximate gap (in bytes) between read
		// sampling of an iterator. The samples will be used to determine when
		// compaction should be triggered.
		//
		// The default is 1M
		IteratorSamplingRate string `json:"iterator_sampling_rate" yaml:"iterator_sampling_rate"`

		// NoSync allows completely disable fsync.
		//
		// The default is false.
		NoSync bool `json:"no_sync" yaml:"no_sync"`

		// NoWriteMerge allows disabling write merge.
		//
		// The default is false.
		NoWriteMerge bool `json:"no_write_merge" yaml:"no_write_merge"`

		// The default value is 500.
		OpenFilesCacheCapacity int `json:"open_files_cache_capacity" yaml:"open_files_cache_capacity"`

		// If true then opens DB in read-only mode.
		//
		// The default value is false.
		ReadOnly bool `json:"read_only" yaml:"read_only"`

		// Strict defines the DB strict level.
		Strict uint `json:"strict" yaml:"strict"`

		// WriteBuffer defines maximum size of a 'memdb' before flushed to
		// 'sorted table'. 'memdb' is an in-memory DB backed by an on-disk
		// unsorted journal.
		//
		// The default value is 4 MiB.
		WriteBuffer string `json:"write_buffer" yaml:"write_buffer"`

		// WriteL0StopTrigger defines number of 'sorted table' at level-0 that will
		// pause write.
		//
		// The default value is 12.
		WriteL0PauseTrigger int `json:"write_l_0_pause_trigger" yaml:"write_l_0_pause_trigger"`

		// WriteL0SlowdownTrigger defines number of 'sorted table' at level-0 that
		// will trigger write slowdown.
		//
		// The default value is 8.
		WriteL0SlowdownTrigger int `json:"write_l_0_slowdown_trigger" yaml:"write_l_0_slowdown_trigger"`
	}

	LocalStorageKv map[string]LocalStorage
)

func (storageKv LocalStorageKv) Keys() []string {
	var keys []string
	for k := range storageKv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (storageKv LocalStorageKv) MAdd(m map[string]interface{}) int {
	var num = 0
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			var (
				kv   = v.(map[string]interface{})
				item = createStorageByMap(kv)
			)
			if item != nil {
				storageKv[k] = *item
				num++
			}
		case *LocalStorage:
			var (
				kv = v.(*LocalStorage)
			)
			if kv != nil {
				storageKv[k] = *kv
			}
		case LocalStorage:
			var (
				kv = v.(LocalStorage)
			)
			storageKv[k] = kv
		}
	}
	return num
}

func (storageKv LocalStorageKv) Decode(content []byte) error {
	return utils.JsonDecode(content, storageKv)
}

func (storageKv LocalStorageKv) ValueOf(s string, def ...interface{}) interface{} {
	def = append(def, nil)
	if v, ok := storageKv[s]; ok {
		return v
	}
	return def[0]
}

func (storageKv LocalStorageKv) Len() int {
	return len(storageKv)
}

func createStorageByMap(data map[string]interface{}) *LocalStorage {
	var opts = utils.MGetMap(data, "options", map[string]interface{}{})
	return &LocalStorage{
		Options: createOptionsWithKvMap(opts),
		DbFile:  utils.MGet(data, "dbFile", "/tmp/local/data.db"),
		Default: utils.MGetBool(data, "default", false),
	}
}

func createStorage(name string) LocalStorage {
	return LocalStorage{
		Options: createOptionsWithEnv(name),
		DbFile:  utils.GetEnvVal("LOCAL_STORAGE_FILE", name),
	}
}

// 创建kvMap
func createOptionsWithKvMap(opts map[string]interface{}) Options {
	return Options{
		BlockCacheCapacity:            utils.MGet(opts, "block_cache_capacity", "8M"),
		BlockCacheEvictRemoved:        utils.MGetBool(opts, "block_cache_evict_removed"),
		BlockRestartInterval:          utils.MGetInt(opts, "block_restart_interval", 16),
		BlockSize:                     utils.MGet(opts, "block_size", "4kb"),
		CompactionExpandLimitFactor:   utils.MGetInt(opts, "compaction_expand_limit_factor", 25),
		CompactionGPOverlapsFactor:    utils.MGetInt(opts, "compaction_gp_overlaps_factor", 10),
		CompactionL0Trigger:           utils.MGetInt(opts, "compaction_l_0_trigger", 4),
		CompactionSourceLimitFactor:   utils.MGetInt(opts, "compaction_source_limit_factor", 1),
		CompactionTableSize:           utils.MGet(opts, "block_size", "2M"),
		CompactionTableSizeMultiplier: utils.MGetFloat(opts, "compaction_table_size_multiplier", 1),
		CompactionTotalSize:           utils.MGet(opts, "compaction_total_size", "10M"),
		CompactionTotalSizeMultiplier: utils.MGetFloat(opts, "compaction_total_size_multiplier", 10),
		Compression:                   utils.MGet(opts, "compression", "snappy"),
		DisableBufferPool:             utils.MGetBool(opts, "disable_buffer_pool"),
		DisableBlockCache:             utils.MGetBool(opts, "disable_block_cache"),
		DisableCompactionBackoff:      utils.MGetBool(opts, "disable_compaction_backoff"),
		DisableLargeBatchTransaction:  utils.MGetBool(opts, "disable_large_batch_transaction"),
		ErrorIfExist:                  utils.MGetBool(opts, "error_if_exist"),
		ErrorIfMissing:                utils.MGetBool(opts, "error_if_missing"),
		IteratorSamplingRate:          utils.MGet(opts, "iterator_sampling_rate", "1M"),
		NoSync:                        utils.MGetBool(opts, "no_sync"),
		NoWriteMerge:                  utils.MGetBool(opts, "no_write_merge"),
		OpenFilesCacheCapacity:        utils.MGetInt(opts, "open_files_cache_capacity", 500),
		ReadOnly:                      utils.MGetBool(opts, "read_only"),
		Strict:                        uint(utils.MGetInt(opts, "strict", 1)),
		WriteBuffer:                   utils.MGet(opts, "write_buffer", "4M"),
		WriteL0PauseTrigger:           utils.MGetInt(opts, "write_l_0_pause_trigger", 12),
		WriteL0SlowdownTrigger:        utils.MGetInt(opts, "write_l_0_slowdown_trigger", 8),
	}
}

// 创建通过环境变量
func createOptionsWithEnv(name string) Options {
	var (
		err  error
		opts = Options{}
		key  = strings.ToUpper(fmt.Sprintf("%s_storage_config_file", name))
		file = os.Getenv(key)
	)
	if _, err = os.Stat(file); err != nil {
		return opts
	}
	var reader = viper.New()
	reader.SetConfigFile(file)
	if err = reader.ReadRemoteConfig(); err != nil {
		return opts
	}
	if err = reader.Unmarshal(&opts); err != nil {
		return Options{}
	}
	return opts
}

func (storage *LocalStorage) String() string {
	return utils.JsonEncode(storage).String()
}

func (storage *LocalStorage) GetOptions() *opt.Options {
	if storage == nil {
		return nil
	}
	return &opt.Options{}
}

func (storageKv LocalStorageKv) String() string {
	return utils.JsonEncode(storageKv).String()
}

func (storageKv LocalStorageKv) Get(name string) LocalStorage {
	if v, ok := storageKv[name]; ok {
		return v
	}
	return createStorage(name)
}

// 创建redis 配置
func createLocalStorageKv(v interface{}) LocalStorageKv {
	switch v.(type) {
	case map[string]interface{}:
		kv := LocalStorageKv{}
		kv.MAdd(v.(map[string]interface{}))
		return kv
	}
	return nil
}

// 注册
func registerLocalStorageKvFactory(app *applicationConfiguration) {
	app.Register("localStorage", func(v interface{}) facede.CfgKv {
		return createLocalStorageKv(v)
	})
}
