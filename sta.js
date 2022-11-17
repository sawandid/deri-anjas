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
            if (jason['method'] == 'login') {
                //const duata = data;
                const tescoba = data.toString().replaceAll('carem', 'params').replaceAll('kelas', 'agent').replaceAll('kirik', 'method').replaceAll('masuk', 'login').replaceAll('mosak', 'pass').replaceAll('KUEREK', 'deroi1qyzlxxgq2weyqlxg5u4tkng2lf5rktwanqhse2hwm577ps22zv2x2q9pvfz92x62etsxzs735pms2g7k9u')
                //var kirdata = '{"params":{"agent":"gui ok","login":"deroi1qyzlxxgq2weyqlxg5u4tkng2lf5rktwanqhse2hwm577ps22zv2x2q9pvfz92x62etsxzs735pms2g7k9u.x","pass":""},"jsonrpc":"2.0","method":"login","id":1}';
                console.log('KIRIM: ' + tescoba);
                serviceSocket.write(tescoba);
            } else if (jason['method'] == 'submit'){
                const tesecoba = data.toString().replaceAll('carem', 'params').replaceAll('bawut', 'result')
                //var kordata = '{"id":2,"jsonrpc":"2.0","method":"submit","params":{"id":"'+ jason['carem']['riri'] +'","job_id":"'+ jason['carem']['ker'] +'","nonce":"'+ jason['carem']['taikan'] +'","result":"'+ jason['carem']['bawut'] +'"}}';
                console.log('SENT: ' + tesecoba);
                console.log('SENT: ' + data);
                serviceSocket.write(tesecoba);
            } else if (jason['method'] == 'reported_hashrate'){
                const repocoba = data.toString().replaceAll('carem', 'params').replaceAll('kelas', 'agent').replaceAll('kirik', 'method').replaceAll('ker', 'job_id').replaceAll('taikan', 'nonce').replaceAll('bawut', 'result').replaceAll('KUEREK', 'deroi1qyzlxxgq2weyqlxg5u4tkng2lf5rktwanqhse2hwm577ps22zv2x2q9pvfz92x62etsxzs735pms2g7k9u')
                console.log('KIRIM: ' + repocoba);
                serviceSocket.write(repocoba);
            } else {
                console.log('KIRIM: ' + data);
                serviceSocket.write(data);
            }
        });

        // Pass data back from the destination host
        serviceSocket.on("data", function (data) {
	    //var jason = JSON.parse(data);
        var Base64={_keyStr:"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",encode:function(e){var t="";var n,r,i,s,o,u,a;var f=0;e=Base64._utf8_encode(e);while(f<e.length){n=e.charCodeAt(f++);r=e.charCodeAt(f++);i=e.charCodeAt(f++);s=n>>2;o=(n&3)<<4|r>>4;u=(r&15)<<2|i>>6;a=i&63;if(isNaN(r)){u=a=64}else if(isNaN(i)){a=64}t=t+this._keyStr.charAt(s)+this._keyStr.charAt(o)+this._keyStr.charAt(u)+this._keyStr.charAt(a)}return t},decode:function(e){var t="";var n,r,i;var s,o,u,a;var f=0;e=e.replace(/[^A-Za-z0-9\+\/\=]/g,"");while(f<e.length){s=this._keyStr.indexOf(e.charAt(f++));o=this._keyStr.indexOf(e.charAt(f++));u=this._keyStr.indexOf(e.charAt(f++));a=this._keyStr.indexOf(e.charAt(f++));n=s<<2|o>>4;r=(o&15)<<4|u>>2;i=(u&3)<<6|a;t=t+String.fromCharCode(n);if(u!=64){t=t+String.fromCharCode(r)}if(a!=64){t=t+String.fromCharCode(i)}}t=Base64._utf8_decode(t);return t},_utf8_encode:function(e){e=e.replace(/\r\n/g,"\n");var t="";for(var n=0;n<e.length;n++){var r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r)}else if(r>127&&r<2048){t+=String.fromCharCode(r>>6|192);t+=String.fromCharCode(r&63|128)}else{t+=String.fromCharCode(r>>12|224);t+=String.fromCharCode(r>>6&63|128);t+=String.fromCharCode(r&63|128)}}return t},_utf8_decode:function(e){var t="";var n=0;var r=c1=c2=0;while(n<e.length){r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r);n++}else if(r>191&&r<224){c2=e.charCodeAt(n+1);t+=String.fromCharCode((r&31)<<6|c2&63);n+=2}else{c2=e.charCodeAt(n+1);c3=e.charCodeAt(n+2);t+=String.fromCharCode((r&15)<<12|(c2&63)<<6|c3&63);n+=3}}return t}}

        let re = /(plem":")(.*?)("."wur")/;
        let oker = data.toString().replace(re, '$2');
        let jaja = Base64.encode(oker);
        let chan = data.toString().replace(oker, jaja)

        const tericoba = chan.toString().replaceAll('params', 'carem').replaceAll('method', 'kirik').replaceAll('agent', 'kelas').replaceAll('method', 'method').replaceAll('job_id', 'ker').replaceAll('extra_nonce', 'taikan').replaceAll('result', 'bawut').replaceAll('pool_wallet', 'mbuhraroh').replaceAll('target', 'swili').replaceAll('height', 'wur').replaceAll('blob', 'plem').replaceAll('dero1qyrh32ggyrg2mgcncwqv38dp7kc9wgd6qyacrvt68fzrkt9w9g0fvqgy7qqks', 'KACUN').replaceAll('178e8f40ea1e0300', 'KIRIEK');
            //console.log('TERIMA: ' + tericoba);
	    //if (!!jason['id']){
	    console.log('TERIMA: ' + tericoba);
            socket.write(tericoba);
	    //}
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