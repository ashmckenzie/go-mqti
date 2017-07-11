local json = require("json")

function process(i)
  local jsonObj = json.decode(i)
  jsonObj.Message = '{"value":' .. tonumber(jsonObj.Message) .. "}"

  return jsonObj, {}
end

