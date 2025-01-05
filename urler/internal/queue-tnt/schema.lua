local log = require('log')
local queue = require('queue')

local exports = {}

exports.init_functions = function(interface)
    for _, v in pairs(interface) do
        box.schema.func.create(v, { setuid = true, if_not_exists = true })
    end
end

exports.init = function()
end

exports.init_queue = function(cfg)
    log.info('initializing queue')

    for _, tube in ipairs(cfg.tubes) do
        queue.create_tube(tube.name, tube.driver, tube.opts)

        log.info('created tube %s', tube.name)
    end
end

return exports
