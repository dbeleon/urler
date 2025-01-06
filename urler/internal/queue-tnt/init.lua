#!/usr/bin/env tarantool
local log = require('log')
local schema = require('schema')
local grants = require('grants')
local interface = require('interface')
local queue_interface = require('queue_interface')
local queue = require('queue')

queue.cfg = {
    in_replicaset = true,
    ttr = 4 * 365 * 24 * 60 * 60
}

box.cfg {
    listen    = os.getenv('TARANTOOL_LISTEN'),
    wal_dir   = os.getenv('TARANTOOL_WAL_DIR'),
    memtx_dir = os.getenv('TARANTOOL_MEMTX_DIR'),
    vinyl_dir = os.getenv('TARANTOOL_VINYL_DIR'),
    memtx_use_mvcc_engine = true,
}

queue_cfg = {
    default_consume_timeout = 5, -- 5 s
    default_pri = 0,
    default_ttr = 600, -- 10min

    tubes = {
        {
            name = 'qr_queue',
            driver = 'limfifottl',
            opts = {
                temporary = false,
                if_not_exists = true,
            }
        }
    }
}

local role_name = 'qr_role'
local user_name = 'qr_user'


log.info('init begins')

schema.init()
schema.init_queue(queue_cfg)

schema.init_functions(interface.get())
schema.init_functions(queue_interface.get())

grants.init_role(role_name, { interface.get(), queue_interface.get() })
grants.makegrants(user_name, role_name, os.getenv('USER_PASS'))

log.info('init completed')

if os.getenv('TARANTOOL_ADMIN_ADDR') then
    require('console').listen(os.getenv('TARANTOOL_ADMIN_ADDR'))
    log.info('admin console addr: %s', os.getenv('TARANTOOL_ADMIN_ADDR'))
end
