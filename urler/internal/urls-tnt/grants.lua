local exports = {}

exports.init_role = function(role_name, interfaces)
    box.schema.role.create(role_name, { if_not_exists = true })

    for _, interface in ipairs(interfaces) do
        for _, v in pairs(interface) do
            box.schema.role.grant(role_name, 'execute', 'function', v, { if_not_exists = true })
        end
    end
end

exports.makegrants = function(user, role, password)
    box.session.su('admin')
    box.schema.user.create(user, { password = password, if_not_exists = true })
    box.schema.user.grant(user, 'execute', 'role', role, { if_not_exists = true })
end

exports.devrole = function(interface)
    local role_name = 'devrole'

    exports.init_role(role_name, interface)

    return role_name
end

return exports
