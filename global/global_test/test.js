import global from "./js/global";

// initial value
var data = {
    name: "Alice",
    id: 1,
    age: 10,
};

// Load from go engine if existed
data = global.Load("alice") || data;

// update
data.age ++

global.Save("alice", data)

console.log(JSON.stringify(data, null, 2))

console.log(data.age);
console.log(data.name);
