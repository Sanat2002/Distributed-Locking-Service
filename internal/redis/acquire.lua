if redis.call("EXISTS", KEYS[1]) == 1 then
  local expiry = redis.call("HGET", KEYS[1], "expiry_ts")
  if tonumber(expiry) > tonumber(ARGV[3]) then
    return {0}
  end
end

local token = redis.call("INCR", KEYS[2])

redis.call("HSET", KEYS[1],
  "owner_id", ARGV[1],
  "expiry_ts", tonumber(ARGV[3]) + tonumber(ARGV[2]),
  "fencing_token", token
)

redis.call("PEXPIRE", KEYS[1], ARGV[2])

return {1, token}
