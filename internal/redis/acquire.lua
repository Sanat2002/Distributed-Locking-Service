-- KEYS[1] = lock key (lock:{resource})
-- KEYS[2] = fencing key (fencing:{resource})

-- ARGV[1] = owner_id
-- ARGV[2] = ttl_ms
-- ARGV[3] = current_time_ms

-- If lock exists and not expired → reject
if redis.call("EXISTS", KEYS[1]) == 1 then
  local expiry = redis.call("HGET", KEYS[1], "expiry_ts")
  if tonumber(expiry) > tonumber(ARGV[3]) then
    return {0}
  end
end

-- Generate fencing token
local token = redis.call("INCR", KEYS[2])

-- Set lock metadata
redis.call("HSET", KEYS[1],
  "owner_id", ARGV[1],
  "expiry_ts", tonumber(ARGV[3]) + tonumber(ARGV[2]),
  "fencing_token", token
)

-- Set TTL
redis.call("PEXPIRE", KEYS[1], ARGV[2])

return {1, token}
