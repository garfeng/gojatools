import json from "./js/encoding.json";
import image from "./js/image";
import http from "./js/net.http";

import ioutil from "./js/io.ioutil";
import io from "./js/io";

function byteToString(arr) {
  if (typeof arr === "string") {
    return arr;
  }
  var str = "",
    _arr = arr;
  for (var i = 0; i < _arr.length; i++) {
    var one = _arr[i].toString(2),
      v = one.match(/^1+?(?=0)/);
    if (v && one.length == 8) {
      var bytesLength = v[0].length;
      var store = _arr[i].toString(2).slice(7 - bytesLength);
      for (var st = 1; st < bytesLength; st++) {
        store += _arr[st + i].toString(2).slice(2);
      }
      str += String.fromCharCode(parseInt(store, 2));
      i += bytesLength - 1;
    } else {
      str += String.fromCharCode(_arr[i]);
    }
  }
  return str;
}

let rect = image._newRectangle({
  Min: { X: 0, Y: 0 },
  Max: { X: 100, Y: 100 },
});

let buff = json.MarshalIndent(rect, "", "  ");

console.log(byteToString(buff));

let resp = http.Get("https://api.github.com/users/garfeng");

let buff2 = ioutil.ReadAll(resp.Body);

resp.Body.Close();

console.log(byteToString(buff2));
