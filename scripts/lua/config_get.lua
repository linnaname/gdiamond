request = function()
  param_value = math.random(1,10000)
  path = "/diamond-server/config?dataId=linname" .. param_value .. "&group=DEFAULT_GROUP" .. param_value
  return wrk.format("GET", path)
end