local queue     = require('queue')
local log       = require('log')

local exports   = {}

local interface = {
    'qr_put',
    'qr_consume',
    'qr_ack',
    'qr_prune',
    'notif_put',
    'notif_consume',
    'notif_ack',
    'notif_prune'
}

local function isempty(s)
    return s == nil or s == ''
end

function qr_put(request)
    -- log.info(request)

    if request.url == nil or type(request.url) ~= 'string' then
        error("bad request")
    end

    if request.pri == nil then
        request.pri = queue_cfg.default_pri
    end

    if request.ttr == nil or request.ttr == 0 then
        request.ttr = queue_cfg.default_ttr
    end

    -- log.info("pri=%d ttr=%d", request.pri, request.ttr)

    local created_task = queue.tube.qr_queue:put({ request.url },
        { pri = request.pri, ttr = request.ttr, delay = request.delay })

    return {
        code = 0,
        message = 'ok',
        id = created_task[1]
    }
end

function qr_consume(request)
    -- log.info("consume %s", request)

    if request.timeout == nil or request.timeout < 0 then
        request.timeout = queue_cfg.default_consume_timeout
    end

    local task = queue.tube.qr_queue:take(request.timeout)


    if task == nil then
        return {
            code = 1,
            message = 'empty_queue'
        }
    end

    log.info('taken qr=%s', task)

    return {
        code = 0,
        message = 'ok',
        id = task[1],
        url = task[3][1]
    }
end

function qr_ack(request)
    log.info('acking qr task=%s', request.id)

    local res = queue.tube.qr_queue:ack(request.id)

    -- log.info('acked task=%s result=%s', request, res)

    return {
        code = 0,
        message = 'ok'
    }
end

function qr_prune()
    queue.tube.qr_queue:truncate()

    log.info('qr queue truncated')
end

function notif_put(request)
    log.info("put notification %s", request)

    if request.url == nil or type(request.url) ~= 'string' or
        request.user_ids == nil or type(request.user_ids) ~= 'table' or
        #request.user_ids < 1 or
        request.qr == nil or type(request.qr) ~= 'string' then
        print("bad request")
        error("bad request")
    end

    if request.pri == nil then
        request.pri = queue_cfg.default_pri
    end

    if request.ttr == nil or request.ttr == 0 then
        request.ttr = queue_cfg.default_ttr
    end

    -- log.info("pri=%d ttr=%d", request.pri, request.ttr)

    local created_task = queue.tube.notif_queue:put({ request.url, request.user_ids, request.qr },
        { pri = request.pri, ttr = request.ttr, delay = request.delay })

    return {
        code = 0,
        message = 'ok',
        id = created_task[1]
    }
end

function notif_consume(request)
    -- log.info("consume %s", request)

    if request.timeout == nil or request.timeout < 0 then
        request.timeout = queue_cfg.default_consume_timeout
    end

    local task = queue.tube.notif_queue:take(request.timeout)

    if task == nil then
        return {
            code = 1,
            message = 'empty_queue'
        }
    end

    log.info('taken notif=%s', task)

    return {
        code = 0,
        message = 'ok',
        id = task[1],
        url_d = task[3][1],
        user_ids = task[3][2],
        qr = task[3][3]
    }
end

function notif_ack(request)
    log.info('acking task notif=%s', request.id)

    local res = queue.tube.notif_queue:ack(request.id)

    -- log.info('acked task notif=%s result=%s', request, res)

    return {
        code = 0,
        message = 'ok'
    }
end

function notif_prune()
    queue.tube.notif_queue:truncate()

    log.info('notif queue truncated')
end

exports.get = function()
    return interface
end

return exports
