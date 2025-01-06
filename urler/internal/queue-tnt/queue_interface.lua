local queue     = require('queue')
local log       = require('log')

local exports   = {}

local interface = {
    'qr_publish',
    'qr_consume',
    'qr_ack',
    'qr_prune'
}

local function isempty(s)
    return s == nil or s == ''
end

function qr_publish(request)
    log.info(request)

    if request.url == nil or type(request.url) ~= 'string' then
        error("bad request")
    end

    if request.pri == nil then
        request.pri = queue_cfg.default_pri
    end

    if request.ttr == nil or request.ttr == 0 then
        request.ttr = queue_cfg.default_ttr
    end

    log.info("pri=%d ttr=%d", request.pri, request.ttr)

    local created_task = queue.tube.qr_queue:put({ request.url },
        { pri = request.pri, ttr = request.ttr, delay = request.delay })

    return {
        code = 0,
        message = 'ok',
        id = created_task[1]
    }
end

function qr_consume(request)
    log.info("consume %s", request)

    if request.timeout == nil or request.timeout < 0 then
        request.timeout = queue_cfg.default_consume_timeout
    end

    local task = queue.tube.qr_queue:take(request.timeout)

    log.info('taken=%s', task)

    if task == nil then
        return {
            code = 1,
            message = 'empty_queue'
        }
    end

    return {
        code = 0,
        message = 'ok',
        id = task[1],
        url = task[3][1]
    }
end

function qr_ack(request)
    log.info('acking task=%s', request.id)

    local res = queue.tube.qr_queue:ack(request.id)

    log.info('acked task=%s result=%s', request, res)

    return {
        code = 0,
        message = 'ok'
    }
end

function qr_prune()
    queue.tube.qr_queue:truncate()

    log.info('qr queue truncated')
end

exports.get = function()
    return interface
end

return exports
