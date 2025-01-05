#!/usr/bin/env tarantool
local log = require('log')
local schema = require('schema')
local grants = require('grants')
local interface = require('interface')

box.cfg {
    listen    = os.getenv('TARANTOOL_LISTEN'),
    wal_dir   = os.getenv('TARANTOOL_WAL_DIR'),
    memtx_dir = os.getenv('TARANTOOL_MEMTX_DIR'),
    vinyl_dir = os.getenv('TARANTOOL_VINYL_DIR'),
    memtx_use_mvcc_engine = true,
}

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
