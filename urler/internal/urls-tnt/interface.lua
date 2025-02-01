local log       = require('log')
local buffer = require('buffer')
local fiber = require('fiber')
local ffi = require('ffi')

local exports   = {}

local addUrlAttempts = 5

local interface = {
    'user_add',
    'url_add',
    'url_get',
    'qr_update',
    'url_shorts',
}

local function res_make(code, message)
    return {
        code = code,
        message = message
    }
end

local function res_ok(fields)
    local res = res_make(0, 'ok')
    if fields ~= nil then
        for k, v in pairs(fields) do res[k] = v end
    end
    return res
end

local res_bad_request = res_make(1, 'bad request')
local res_not_found = res_make(2, 'not found')

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

        return res_ok({id = tpl[1]})
    end

    return res_bad_request
end

function url_add(req)
    log.info('add url user=%s long=%s short=%s', req.user_id, req.long, req.short)

    if req.user_id ~= nil and type(req.user_id) == 'number'
        and req.long ~= nil and type(req.long) == 'string'
        and req.short ~= nil and type(req.short) == 'string'
    then

        local usr = box.space.usr:select{req.user_id}
        -- log.info("users %s found %s", #usr, usr)
        if #usr ~= 1 then
            return res_not_found
        end

        for i = 1, addUrlAttempts do -- Делаем попытки вставки
            local status, err = pcall(function()
                box.begin()
                -- TODO: isolation level

                local res = box.space.url.index.long_index:get(req.long)
                -- log.info("long url found: %s", res)
                if res == nil then
                    local tplUrl = box.space.url:insert({ nil, req.long, req.short, nil })
                    -- log.info("url added=%s", tplUrl)
                    req.id = tplUrl[1]
                else
                    -- log.info("url found=%s", res[1])
                    req.id = res[1]
                    req.short = res[3]
                end

                local usrUrl = { req.user_id, req.id }
                if box.space.usr_url:get(usrUrl) == nil then
                    -- log.info("add user=%s url=%s", usrUrl[1], usrUrl[2])
                    box.space.usr_url:insert(usrUrl)
                end

                box.commit()
            end)
    
            if status then
                print(string.format("✅ [SUCCESS] add url long='%s' and short='%s'", req.long, req.short))
                return res_ok({url = req})
            end
    
            print(string.format("🔁 [RETRY %d] add url failed long='%s' and short='%s': %s", i, req.long, req.short, err))
            fiber.sleep(0.1) -- Ждём и пробуем ещё раз
        end

        return res_make(3, "transation rollback")
    end

    return res_bad_request
end

function qr_update(req)
    log.info('update qr code short=%s ', req.short)
    log.info('qr type=%s', type(req.qr))
    if req.short ~= nil and type(req.short) == 'string'
        and req.qr ~= nil and type(req.qr) == 'string'
    then
        local userIDs = {}

        box.begin()

        local res = box.space.url.index.short_index:get(req.short)
        -- log.info("short url found: %s", res)
        if res == nil then
            return res_not_found
        end

        --[[ local bytes = req.qr:gsub('.', function (c)
            return string.byte(c)
        end)

        local tmpbuf = buffer.ibuf()
        local p = tmpbuf:alloc(4 + #bytes)
        p[0] = 0x91 -- MsgPack code for "array-1"
        p[1] = 0xC5 -- MsgPack code for "bin-16" so up to 65536 bytes
        p[2] = #bytes / 256
        p[3] = #bytes % 256
        for i, c in pairs(bytes) do p[i + 4 - 1] = c end
        C insert func
            API_EXPORT int
            box_update(uint32_t space_id, uint32_t index_id, const char *key,
                const char *key_end, const char *ops, const char *ops_end,
                int index_base, box_tuple_t **result)
            {
                mp_tuple_assert(key, key_end);
                mp_tuple_assert(ops, ops_end);
                struct request request;
                memset(&request, 0, sizeof(request));
                request.type = IPROTO_UPDATE;
                request.space_id = space_id;
                request.index_id = index_id;
                request.key = key;
                request.key_end = key_end;
                request.index_base = index_base;
                /** Legacy: in case of update, ops are passed in in request tuple */
                request.tuple = ops;
                request.tuple_end = ops_end;
                return box_process1(&request, result);
            }
        ]]
        --[[ ffi.cdef[[int box_update(uint32_t space_id,
                                uint32_t index_id
                                const char *key,
                                const char *key_end,
                                const char *ops,
                                const char *ops_end,
                                int index_base,
                                box_tuple_t **result);
        ffi.C.box_update(box.space.url.id, res[1][1], k, k_e, tmpbuf.rpos, tmpbuf.wpos, 1, nil)
        ffi.C.box_insert(box.space.url.id, tmpbuf.rpos, tmpbuf.wpos, nil)
        tmpbuf:recycle()
        ]]

        local tpl = box.space.url:update({res[1]},{{'=', 4, req.qr}})
        local usrUrlTpls = {}
        if #tpl > 0 then
            usrUrlTpls = box.space.usr_url.index.url_index:select(res[1])
        end

        box.commit()

        if #tpl == 0 then
            return {
                code = 2,
                message = 'url id not found'
            }
        end

        if #usrUrlTpls > 0 then
            for _,v in pairs(usrUrlTpls) do
                table.insert(userIDs, v[1])
            end
        end

        log.info("url updated=%s", tpl)

        return res_ok({user_ids = userIDs})
    end

    return res_bad_request
end

function url_get(req)
    log.info('get url %s', req) 

    if req ~= nil
        and req.short ~= nil and type(req.short) == 'string'
    then
        local res = box.space.url.index.short_index:get{req.short}
        if res == nil then
            log.info("url not found")
            return res_not_found
        end

        log.info("url found=%s", res)

        return res_ok({url = { id = res[1], long = res[2], short = res[3], qr = res[4] }})
    end

    return res_bad_request
end

function url_shorts(req)
    if req ~= nil
        and req.limit ~= nil and type(req.limit) == 'number'
        and req.offset ~= nil and type(req.offset) == 'number'
    then
        local tpl = box.space.url:select({0},{iterator='GT', limit = req.limit, offset = req.offset})
        local res = {}
        for _,v in ipairs(tpl) do
            table.insert(res, v[3])
        end

        return res_ok({shorts = res})
    end

    return res_bad_request
end

return exports
