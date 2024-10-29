local key = KEYS[1]
local cntKey = key..":cnt"
local expectedCode = ARGV[1]

-- 验证次数
local cnt = tonumber(redis.call("get"), cntKey)
-- 验证码
local code = redis.call("get", key)

-- 验证次数耗尽
if cnt <=0 then
    return -1
end
-- 验证码相等
-- 不能删除验证码，因为如果你删除了验证码
-- 就可以再次发送验证码消耗短信额度
if code == expectedCode then
    -- 把次数标记为 -1，认为验证码不可用
    redis.call("set", cntKey, -1)
    return 0
else
    -- 可能用户偶尔输错了
    redis.call("decr", cntKey, -1)
    return -2
end