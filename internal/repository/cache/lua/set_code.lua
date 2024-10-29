-- 验证码在 Redis 上的 key
-- phone_code:login:157xxxx
local key = KEYS[1]
-- 验证次数，一个验证码最多重复三次
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key存在，但是没有过期时间
    -- 系统错误
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end