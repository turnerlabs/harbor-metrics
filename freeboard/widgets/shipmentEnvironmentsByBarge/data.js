var result = [];
result[0] = [];
datasources["shipmentEnvironmentsByBarge"].map(function(item) {
  objArray = []
  objArray.push(item.barge);
  objArray.push(item.count);
  result[0].push(objArray);
});
return result;