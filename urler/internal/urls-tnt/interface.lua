local log       = require('log')

local exports   = {}

local interface = {
    'user_add',
    'url_add',
    'url_get'
}

exports.get = function()
    return interface
end

function user_add(req)
    log.info('add user %s', req) 

    if req ~= nil
        and req.name ~= nil and type(req.name) == 'string'
        and req.email ~= nil and type(req.email) == 'string'
        and string.find(req.email, '@')
    then
        local tpl = box.space.usr:insert({ nil, req.name, req.email })

        log.info("user added=%s", tpl)

        return {
            id = tpl[1]
        }
    end

    error('bad_request')
end

function url_add(req)
    log.info('add url user=%s long=%s short=%s', req.user_id, req.long, req.short)

    if req ~= nil
        and req.user_id ~= nil and type(req.user_id) == 'number'
        and req.long ~= nil and type(req.long) == 'string'
        and req.short ~= nil and type(req.short) == 'string'
    then

        local usr = box.space.usr:select{req.user_id}
        -- log.info("users %s found %s", #usr, usr)
        if #usr ~= 1 then
            error("user not found")
        end

        box.begin()
        -- TODO: isolation level

        local res = box.space.url.index.long_index:select(req.long)
        -- log.info("long url found: %s", res)
        if #res == 0 then
            local tpl = box.space.url:insert({ nil, req.user_id, req.long, req.short })
            -- log.info("url added=%s", tpl)
            req.id = tpl[1]
        else
            -- log.info("url found=%s", res[1])
            req.id = res[1][1]
            req.short = res[1][4]
        end

        box.commit()

        return req
    end

    error('bad_request')
end


function url_get(req)
    log.info('get url %s', req) 

    if req ~= nil
        and req.short ~= nil and type(req.short) == 'string'
    then
        local res = box.space.url.index.short_index:select{req.short}
        if #res == 0 then
            log.info("url not found")
            return { id = 0 }
        end

        log.info("url found=%s", res[1])

        return res[1]
    end

    error('bad_request')
end

return exports
