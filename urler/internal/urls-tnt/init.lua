#!/usr/bin/env tarantool
local log = require('log')
local schema = require('schema')
local grants = require('grants')
local interface = require('interface')
local repl_pswd = os.getenv('TARANTOOL_REPLICATION_PASSWORD')

box.cfg {
    listen          = os.getenv('TARANTOOL_LISTEN'),
    wal_dir         = os.getenv('TARANTOOL_WAL_DIR'),
    memtx_dir       = os.getenv('TARANTOOL_MEMTX_DIR'),
    vinyl_dir       = os.getenv('TARANTOOL_VINYL_DIR'),
    memtx_memory    = tonumber(os.getenv('TARANTOOL_MEMTX_MEM') .. '', 10),
    -- memtx_use_mvcc_engine = true,
    replication     = {'replicator:' .. repl_pswd .. '@urls-tnt-m:3301',  -- URI мастера
                        'replicator:' .. repl_pswd .. '@urls-tnt-r1:3301', -- URI реплики 1
                        'replicator:' .. repl_pswd .. '@urls-tnt-r2:3301' -- URI реплики 2
                    },
    read_only = string.lower(os.getenv('TARANTOOL_IS_REPLICA') .. '') == "true",
    -- sync replication
    -- replication_synchro_quorum = 2, -- "N / 2 + 1"
}

box.once("schema", function()
    box.schema.user.create('replicator', {password = repl_pswd})
    box.schema.user.grant('replicator', 'replication') -- настроить роль для репликации
    print('box.once executed on replica')
end)

local role_name = 'urls_role'
local user_name = 'urls_user'


log.info('init begins')

schema.init()

schema.init_functions(interface.get())

grants.init_role(role_name, { interface.get() })
grants.makegrants(user_name, role_name, os.getenv('USER_PASS'))

log.info('init completed')

if os.getenv('TARANTOOL_ADMIN_ADDR') then
    require('console').listen(os.getenv('TARANTOOL_ADMIN_ADDR'))
    log.info('admin console addr: %s', os.getenv('TARANTOOL_ADMIN_ADDR'))
end

local metrics = require('metrics')
local prometheus = require('metrics.plugins.prometheus')

metrics.cfg{}
metrics.enable_default_metrics()
metrics.set_global_labels{alias = 'my-tnt-app'}

local httpd = require('http.server').new('0.0.0.0', 3380)
httpd:route( { path = '/metrics' }, prometheus.collect_http)
httpd:start()
