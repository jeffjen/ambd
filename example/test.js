"user strict"

var Redis = require("ioredis")

var inst = new Redis({host: "ambassador", enableOfflineQueue: false});

function role() {
    console.log("role")

    inst.role().then(function(d) {
        console.log(d);
        setTimeout(role, 5000);
    })
    .catch(function(e) {
        console.log(e);
        setTimeout(role, 200);
    });
}

role();
