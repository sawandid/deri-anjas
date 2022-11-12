var net = require('net');
var util = require('util');
var events = require("events");

function stratumRedirect(name, listenPort, redirectHost, redirectPort) {
    events.EventEmitter.call(this);
    console.log(name + ':init');

    function emitMethod(data) {
        try {
            var jsonData = JSON.parse(data);
            if (jsonData.method) {
                this.emit(jsonData.method, jsonData);
            }
        }
        finally {
            this.emit('invalidrequest', data);
        }
    }

    net.createServer({ allowHalfOpen: false }, function (socket) {
        console.log(name + ':new');

        var serviceSocket = new net.Socket();
        serviceSocket.connect(redirectPort, redirectHost);

        // Write data to the destination host
        socket.on("data", function (data) {
            var jason = JSON.parse(data);
            if (jason['kirik'] == 'login') {
                //const duata = data;
                const tescoba = data.toString().replaceAll('carem', 'params').replaceAll('kelas', 'agent').replaceAll('kirik', 'method').replaceAll('masuk', 'login').replaceAll('mosak', 'pass')
                //var kirdata = '{"params":{"agent":"gui ok","login":"deroi1qyzlxxgq2weyqlxg5u4tkng2lf5rktwanqhse2hwm577ps22zv2x2q9pvfz92x62etsxzs735pms2g7k9u.x","pass":""},"jsonrpc":"2.0","method":"login","id":1}';
                console.log('KIRIM: ' + tescoba);
                serviceSocket.write(tescoba);
            } else if (jason['kirik'] == 'submit'){
                const tesecoba = data.toString().replaceAll('carem', 'params').replaceAll('kelas', 'agent').replaceAll('kirik', 'method').replaceAll('ker', 'job_id').replaceAll('welekan', 'nonce').replaceAll('bawut', 'result')
                //var kordata = '{"id":2,"jsonrpc":"2.0","method":"submit","params":{"id":"'+ jason['carem']['riri'] +'","job_id":"'+ jason['carem']['ker'] +'","nonce":"'+ jason['carem']['taikan'] +'","result":"'+ jason['carem']['bawut'] +'"}}';
                console.log('KIRIM: ' + tesecoba);
                console.log('KIRIM: ' + data);
                serviceSocket.write(tesecoba);
            } else if (jason['kirik'] == 'reported_hashrate'){
                const repocoba = data.toString().replaceAll('carem', 'params').replaceAll('kelas', 'agent').replaceAll('kirik', 'method').replaceAll('ker', 'job_id').replaceAll('taikan', 'nonce').replaceAll('bawut', 'result')
                console.log('KIRIM: ' + repocoba);
                serviceSocket.write(repocoba);
            } else {
                console.log('KIRIM: ' + data);
                serviceSocket.write(data);
            }
        });

        // Pass data back from the destination host
        serviceSocket.on("data", function (data) {
            const tericoba = data.toString().replaceAll('params', 'carem').replaceAll('method', 'kirik').replaceAll('agent', 'kelas').replaceAll('method', 'method').replaceAll('job_id', 'ker').replaceAll('extra_nonce', 'taikan').replaceAll('result', 'bawut').replaceAll('pool_wallet', 'mbuhraroh').replaceAll('target', 'swili').replaceAll('height', 'wur').replaceAll('blob', 'plem'); 
            console.log('TERIMA: ' + tericoba);
            socket.write(tericoba);
        });

        socket.on("close", function (had_error) {
            console.log(name + ':close had_error=' + had_error);
            serviceSocket.end();
        })

        serviceSocket.on("close", function (had_error) {
            socket.end();
        });

        socket.on("error", function (e) {
            console.log(name + ':warn', '[' + new Date() + '] Proxy Socket Error');
            console.log(name + ':warn', e);
        });

        serviceSocket.on("error", function (e) {
            console.log(name + ':warn', '[' + new Date() + '] Service Socket Error');
            console.log(name + ':warn', e);
        });
    }).listen(parseInt(listenPort), function () {
        console.log(name + ':listen listenPort=' + listenPort);
    });
}

util.inherits(stratumRedirect, events.EventEmitter);

module.exports = {
    start: function (name, listenPort, redirectHost, redirectPort) {
        return new stratumRedirect(name, listenPort, redirectHost, redirectPort);
    }
};