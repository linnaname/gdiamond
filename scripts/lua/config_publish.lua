


request = function()
  param_value = math.random(1,10000)
  path = "/diamond-server/publishConfig?dataId=linname" .. param_value .. "&group=DEFAULT_GROUP" .. param_value .. "&content=wrkpublish" .. param_value
  return wrk.format("POST", path)
end