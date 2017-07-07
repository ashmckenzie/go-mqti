local json = require("json")

function match(i)
	local jsonObj = json.decode(i)

	if jsonObj._type == "location" then
		return true
	else
		return false
	end
end
