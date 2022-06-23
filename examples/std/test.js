import json from "./js/encoding.json"
import image from "./js/image";
import http from "./js/net.http"

import ioutil from "./js/io.ioutil"
import io from "./js/io";


let rect = image._newRectangle()
rect.Min = image._newPoint()
rect.Max = image.Pt(100, 100)

let buff = json.MarshalIndent(rect, "", "  ")

console.log(String.fromCharCode(...buff));


let resp = http.Get("https://www.baidu.com/")

let buff2 = ioutil.ReadAll(resp.Body)


console.log(resp);

resp.Body.Close()

