var result = [];
result[0] = [];
datasources["shipmentEnvironmentsByGroup"].map(function(item) {
  objArray = []
  objArray.push(item.group);
  objArray.push(item.count);
  result[0].push(objArray);
});
return result;