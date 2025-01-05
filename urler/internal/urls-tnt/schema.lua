local log = require('log')

local exports = {}

exports.init_functions = function(interface)
    for _, v in pairs(interface) do
        box.schema.func.create(v, { setuid = true, if_not_exists = true })
    end
end

exports.init = function()
    box.schema.sequence.create('user_id_seq', { if_not_exists = true })
    -- usr space
    local usr_space = box.schema.space.create('usr', {
        if_not_exists = true,
        engine = 'memtx',
        format = {
            { name = 'id',      type = 'unsigned',  is_nullable = false },
            { name = 'name',    type = 'string',    is_nullable = false },
            { name = 'email',   type = 'string',    is_nullable = false },
            -- { name = 'qr',       type = 'string', is_nullable = false },
        },
    })

    usr_space:create_index('primary', {
        sequence = 'user_id_seq',
        if_not_exists = true,
    })

    box.schema.sequence.create('url_id_seq', { if_not_exists = true })
    -- url space
    local url_space = box.schema.space.create('url', {
        if_not_exists = true,
        engine = 'memtx',
        format = {
            { name = 'id',      type = 'unsigned',  is_nullable = false },
            { name = 'user_id', type = 'number',    is_nullable = false },
            { name = 'long',    type = 'string',    is_nullable = false },
            { name = 'short',   type = 'string',    is_nullable = false },
            -- { name = 'qr',       type = 'string', is_nullable = false },
        },
    })

    url_space:create_index('primary', {
        sequence = 'url_id_seq',
        if_not_exists = true,
    })
    
    url_space:create_index('long_index', {
        type = 'HASH',
        unique = true,
        if_not_exists = true,
        parts = {
            { 'long', is_nullable = false, type = 'string' }
        }
    })
    
    url_space:create_index('short_index', {
        type = 'HASH',
        unique = true,
        if_not_exists = true,
        parts = {
            { 'short', is_nullable = false, type = 'string' }
        }
    })
end

return exports
