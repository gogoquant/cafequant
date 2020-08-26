var chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';
// encoder
// [https://gist.github.com/999166] by [https://github.com/nignag]
function btoa(input) {
    var str = String(input);
    for (
        // initialize result and counter
        var block, charCode, idx = 0, map = chars, output = '';
        // if the next str index does not exist:
        //   change the mapping table to "="
        //   check if d has no fractional digits
        str.charAt(idx | 0) || (map = '=', idx % 1);
        // "8 - idx % 1 * 8" generates the sequence 2, 4, 6, 8
        output += map.charAt(63 & block >> 8 - idx % 1 * 8)
    ) {
        charCode = str.charCodeAt(idx += 3 / 4);
        if (charCode > 0xFF) {
            throw new InvalidCharacterError("'btoa' failed: The string to be encoded contains characters outside of the Latin1 range.");
        }
        block = block << 8 | charCode;
    }
    return output;
};

// decoder
// [https://gist.github.com/1020396] by [https://github.com/atk]
function atob(input) {
    var str = String(input).replace(/=+$/, '');
    if (str.length % 4 == 1) {
        throw new InvalidCharacterError("'atob' failed: The string to be decoded is not correctly encoded.");
    }
    for (
        // initialize result and counters
        var bc = 0, bs, buffer, idx = 0, output = '';
        // get next character
        buffer = str.charAt(idx++);
        // character found in table? initialize bit storage and add its ascii value;
        ~ buffer && (bs = bc % 4 ? bs * 64 + buffer : buffer,
            // and if not first of each 4 characters,
            // convert the first 8 bits to one ascii character
            bc++ % 4) ? output += String.fromCharCode(255 & bs >> (-2 * bc & 6)) : 0
    ) {
        // try to find character in table (0-63, not found => -1)
        buffer = chars.indexOf(buffer);
    }
    return output;
};

var my_weor = function() {
    var bo = {
        'O1': function(g) {
            return function(e, f) {
                return function(a) {
                    return {
                        p0: a
                    }
                } (function(a) {
                    var b, E0 = 0;
                    for (var c = e; E0 < a['length']; E0++) {
                        var d = f(a, E0);
                        b = E0 === 0 ? d: b ^ d
                    }
                    return b ? c: !c
                })
            } (function(a, b, c, d) {
                var e = 785;
                var f = d(b, c) - a(g, e);
                return true
            } (parseInt, Date,
            function(a) {
                return ('' + a)['substring'](1, (a + '')['length'] - 1)
            } ('_getTime2'),
            function(a, b) {
                return new a()[b]()
            }),
            function(a, b) {
                var c = parseInt(a['charAt'](b), 16)['toString'](2);
                return c['charAt'](c['length'] - 1)
            })
        } ('ecg6mf6ar')
    };
    var bp = function(a) {
        var b = new Array();
        var i;
        if (a && a.length > 0) {
            var s = a.split('*');
            for (i = 0; i < s.length - 1; i++) {
                switch (i % 3) {
                case 0:
                    b += String.fromCharCode(parseInt(s[i], 8));
                    break;
                case 1:
                    b += String.fromCharCode(parseInt(s[i], 10));
                    break;
                case 2:
                    b += String.fromCharCode(parseInt(s[i], 16));
                    break
                }
            }
            return b
        } else {
            return ''
        }
    };
    var screen = {
        'height': 640,
        'width': 360,
    };
    var seajs = {'version': "1.2.1"};
    var Q = {'page': "play"};
    var window = {
        //'orientation': undefined,
        'devicePixelRatio': 3,
        'screenTop': 0,
        'outerHeight': 640,
        '__page_start': new Date().getTime(),
        'seajs': seajs,
        'Q':Q,
        //'ucweb': undefined,
        
    };

    this.weor = function (z, A, B, C, D, E) {
        var F = function() {
            K = K > L ? L: K
        };
        var G = function() {
            var s = function() {
                W += (W ? '_': '') + btoa('LBW')
            };
            W = ifDef('ucweb') ? btoa('UCW') : '';
            W += ifDef('_boluoWebView') ? (W ? '_': '') + btoa('BOL') : '';
            ifDef('isLBBrowser') || ifDef('ks_liebaoversion') ? s() : '';
            try {
                var t = function() {
                    var j = function() {
                        W += (W ? '_': '') + btoa('ITNS')
                    };
                    var k = function() {
                        var a = 'di';
                        q.style = a
                    };
                    var l = function() {
                        var a = '0';
                        q.width = a
                    };
                    var m = function() {
                        var a = '0';
                        q.height = a
                    };
                    var n = bo.O1.p0('841b') ? 63 : 'qd';
                    n += '_';
                    n += 'd';
                    n += 'ns';
                    n += 'c';
                    n += 'a';
                    n += 'c';
                    n += 'he';
                    var o = bo.O1.p0('cd') ? window.localStorage: 2;
                    o.getItem(n) ? j() : '';
                    var p = bo.O1.p0('8f32') ? '6': '<body>' + '<script>' + 'function e(e){window.location.href=n;var o=+new Date;setTimeout(function(){+new Date-o<1e3+e&&(c++,c>1&&a.setItem(t,""+(new Date).getTime()))},e)}try{var t="qd_dnscache",a=window.localStorage,n=atob("aHR0cHM6Ly9pdHVuZXMuYXBwbGUuY29tL2lkaGhoZGRkZC5wbmc="),c=0;e(3e3),e(6e3)}catch(o){}' + '</script></body>';
                    var q = bo.O1.p0('753') ? document.createElement('iframe') : 13;
                    var r = function() {
                        v8string += '%5';
                        _md5(opt, i, i & 63, a);
                        bC.sc = thgirph11;
                        var b = bo.O1.p0('e') ? screen.width: 1;
                        v8string += '0'
                    };
                    k();
                    q.style += 's';
                    q.style += 'p';
                    q.style += 'la';
                    q.style += 'y';
                    q.style += ':';
                    q.style += 'n';
                    q.style += 'one';
                    l();
                    m();
                    document.body.appendChild(q);
                    q.contentWindow.document.open();
                    q.contentWindow.document.write(p);
                    q.contentWindow.document.close();
                    setTimeout(function() {
                        var h = function() {
                            var f = function() {
                                var c = function() {
                                    o.removeItem(n)
                                };
                                var d = function(a, b) {
                                    return a > b
                                };
                                var e = bo.O1.p0('2') ? new Date().getTime() : 7;
                                d(e - g, 11000) ? c() : ''
                            };
                            var g = bo.O1.p0('fec4') ? Number(o.getItem(n)) : ':';
                            g ? f() : ''
                        };
                        o.getItem(n) ? h() : '';
                        document.body.removeChild(q)
                    },
                    10000)
                };
                var u = bo.O1.p0('cebe') ? navigator.userAgent.toLowerCase() : 'LBW';
                var v = bo.O1.p0('982') ? u.match(/mozilla\/(?:\d+(?:\.\d+)+) \(iphone; cpu iphone os (?:\d+(?:\_\d+)+) like mac os x\) applewebkit\/(?:\d+(?:\.\d+)+) \(khtml, like gecko\) version\/(?:\d+(?:\.\d+)+) mobile\/\w+ safari\/(?:\d+(?:\.\d+)+)/) : '_';
                v = v == u;
                v ? t() : ''
            } catch(e) {}
        };
        var H = function() {
            _md5(15, i, 0, A ? '933653760616065683236663733603e3': '46434306535376731313637303162313');
            var h = function() {
                x[i >> 2] |= (parseInt(a.substr((j >> 2) * 8, 8).split('').reverse().join(''), 16) >> 8 * (j % 4) & 255 ^ j % 8) << ((i++&3) << 3);
                _md5(16, i, j + 1, a)
            };
            v8string += '7D';
            var k = 'h5';
            var l = function() {
                var f = function() {
                    var c = function() {
                        storage.removeItem(StorageName)
                    };
                    var d = function(a, b) {
                        return a > b
                    };
                    var e = bo.O1.p0('2') ? new Date().getTime() : 7;
                    d(e - g, 11000) ? c() : ''
                };
                var g = bo.O1.p0('fec4') ? Number(storage.getItem(StorageName)) : ':';
                g ? f() : ''
            };
            var m = function() {
                x[i >> 2] |= a.charCodeAt(i) << 8 * (i++%4);
                _md5(3, i, -1, a)
            };
            v8string += 't'
        };
        var I = function(a, b) {
            return a === b
        };
        var I = function(a, b) {
            return a === b
        };
        var J = function(a, b) {
            return a > b
        };
        var K = bo.O1.p0('6') ? screen.height: 'BOL';
        var L = bo.O1.p0('e') ? screen.width: 1;
        var O = bo.O1.p0('c') ? window.orientation: 14;
        I(O, 90) || I(O, -90) ? F() : '';
        var P = bo.O1.p0('4') ? window.devicePixelRatio: '';
        K = Math.round(K / P);
        var Q = bo.O1.p0('2') ? Math.round(window.screenTop / P) : 1;
        var R = bo.O1.p0('18') ? Math.round(window.outerHeight / P) : 'n';
        var S = function() {
            bC.t = d - V[V[V[3]] - 1];
            bC.src += 'c';
            v8string += 'n';
            bC.src += '0'
        };
        var T = bo.O1.p0('1a') ? K - R - Q: 10;
        var U = bo.O1.p0('719a') ? btoa(Q + '_' + T) : '32';
        var V = bo.O1.p0('d') ? ['slice', 'call', 'querySelectorAll', 'length', 'push', 'shift', 'indexOf', 'document', 'innerHTML', 'match', 'forEach'] : '_boluoWebView';
        var W, str, d, z = escape(btoa(z)),
        M,
        N,
        h = [M = 1732584193, N = -271733879, ~M, ~N],
        x = [];
        d = new Date().getTime();
        V.push((V[V[0]]( - 5).join('')[V[3]] - 5).toString(16)); ! A ? G() : '';
        str = (!A ? d - 7 : E + '' + D) + '';
        str = escape(!A ? btoa(str) : btoa(str + C + '' + B));
        function _md5(n, i, j, a) {
            var o = function() {
                var g = function() {
                    var d = function() {
                        var b = function() {
                            h = [add(a[0], h[0]), add(a[1], h[1]), add(a[2], h[2]), add(a[3], h[3])];
                            _md5(n, i + (15 << 6), i & 63, h)
                        };
                        var c = function() {
                            _md5(n, i, i & 63, a)
                        };
                        a = [a[3], add(a[1], (M = add(add(a[0], [a[1] & a[2] | ~a[1] & a[3], a[3] & a[1] | ~a[3] & a[2], a[1] ^ a[2] ^ a[3], a[2] ^ (a[1] | ~a[3])][N = j >> 4]), add(Math.abs(Math.sin(j + 1)) * 4294967296 | 0, x[[j, 5 * j + 1, 3 * j + 5, 7 * j][N] % 16 + (i++>>>6)]))) << (N = [7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21][4 * N + j % 4]) | M >>> 32 - N), a[1], a[2]]; ! (i & 63) ? b() : c()
                    };
                    var e = function() {
                        var b = function() {
                            var a = '';
                            str = a
                        };
                        x = [];
                        b();
                        _md5(n, 0, -3, a)
                    };
                    var f = function(a, b) {
                        return a < b
                    };
                    f(i, str << 6) ? d() : e()
                };
                var k = function() {
                    var c = function() {
                        x[i >> 2] |= a.charCodeAt(i) << 8 * (i++%4);
                        _md5(3, i, -1, a)
                    };
                    var d = function() {
                        _md5(15, i, 0, A ? '933653760616065683236663733603e3': '46434306535376731313637303162313')
                    };
                    var e = function(a, b) {
                        return a < b
                    };
                    e(i, a.length) ? c() : d()
                };
                var l = function() {
                    var c = function() {
                        str += (h[i >> 3] >> (1 ^ i++&7) * 4 & 15).toString(16);
                        _md5(n, i, j--, a)
                    };
                    var d = function(a, b) {
                        return a < b
                    };
                    d(i, 32) ? c() : ''
                };
                var m = function(a, b) {
                    return a >= b
                };
                m(j, 0) ? g() : j < 0 && j > -3 ? k() : l()
            };
            var p = function() {
                var c = function() {
                    x[i >> 2] |= (parseInt(a.substr((j >> 2) * 8, 8), 16) >> 8 * (j % 4) & 255 ^ j % 4) << ((i++&3) << 3);
                    _md5(9, i, j + 1, a)
                };
                var d = function() {
                    _md5(12, i, !ifDef('Q') * 1, z)
                };
                var e = function(a, b) {
                    return a < b
                };
                e(j, a.length >> 1) ? c() : d()
            };
            var q = function() {
                var c = function() {
                    x[i >> 2] |= n.charCodeAt(j++) << 8 * (i % 4);
                    _md5(12, ++i, j, z)
                };
                var d = function() {
                    var a = function() {
                        x[i >> 2] |= 1 << (i % 4 << parseFloat("1.2.1") + 1.8) + 7
                        //x[i >> 2] |= 1 << (i % 4 << parseFloat(new Function('return ' + atob('d2luZG93LnNlYWpzICYmIHNlYWpzLnZlcnNpb24='))()) + 1.8) + 7
                    };
                    ifDef(atob('X19wYWdlX3N0YXJ0')) ? a() : '';
                    x[str = (i + 8 >> 6 << 4) + 14] = i << 3;
                    _md5(3, 0, 0, h)
                };
                var e = function(a, b) {
                    return a < b
                };
                n = atob(unescape(a));
                e(j, n.length) ? c() : d()
            };
            var r = function() {
                var c = function() {
                    x[i >> 2] |= (parseInt(a.substr((j >> 2) * 8, 8).split('').reverse().join(''), 16) >> 8 * (j % 4) & 255 ^ j % 8) << ((i++&3) << 3);
                    _md5(16, i, j + 1, a)
                };
                var d = function() {
                    _md5(7, i, 0, A ? '60643662346163366731623261643565': '33316439376031313066333231336563')
                };
                var e = function(a, b) {
                    return a < b
                };
                e(j, a.length >> 1) ? c() : d()
            };
            var s = function() {
                var f = function() {
                    var c = function() {
                        x[i >> 2] |= a.charCodeAt(i) << 8 * (i++%4);
                        _md5(3, i, -1, a)
                    };
                    var d = function() {
                        _md5(15, i, 0, A ? '933653760616065683236663733603e3': '46434306535376731313637303162313')
                    };
                    var e = function(a, b) {
                        return a < b
                    };
                    e(i, a.length) ? c() : d()
                };
                inFn('WebkitAppearance', document.documentElement.style) ? R1AgI() : '';
                var g = function() {
                    var c = function() {
                        x[i >> 2] |= (parseInt(a.substr((j >> 2) * 8, 8).split('').reverse().join(''), 16) >> 8 * (j % 4) & 255 ^ j % 8) << ((i++&3) << 3);
                        _md5(16, i, j + 1, a)
                    };
                    var d = function() {
                        _md5(7, i, 0, A ? '60643662346163366731623261643565': '33316439376031313066333231336563')
                    };
                    var e = function(a, b) {
                        return a < b
                    };
                    e(j, a.length >> 1) ? c() : d()
                };
                var h = bo.O1.p0('18') ? Math.round(window.outerHeight / k) : 'n';
                var k = bo.O1.p0('4') ? window.devicePixelRatio: ''
            };
            var t = function(a, b) {
                return a > b
            };
            var u = function(a, b) {
                return a < b
            };
            var v = function() {
                v8string += 'u';
                str = (!A ? d - 7 : E + '' + D) + '';
                var c = function() {
                    _md5(7, i, 0, A ? '60643662346163366731623261643565': '33316439376031313066333231336563')
                };
                var e = function() {
                    var b = function() {
                        var a = 's';
                        jst = a
                    };
                    b();
                    jst += 'i';
                    jst += 'j';
                    jst += 'sc'
                }
            };
            var w = function() {
                var b = bo.O1.p0('fb9') ? {}: 0;
                bC.__jsT = X();
                v8string += 'u';
                v8string += '2';
                n = atob(unescape(a));
                bC.src += '32'
            };
            t(n, 0) && u(n, 5) ? o() : n > 6 && n < 10 ? p() : n > 11 && n < 14 ? q() : n > 14 && n < 17 ? r() : ''
        }
        while(seaNums > 0){
            V[V[4]] += 1;
            seaNums--;
        };
        while(sbaseNums > 0){
            V[V[5]] += 1;
            sbaseNums--;
        };
/*
        V[V[0]][V[1]](window[V[7]][V[2]]('script'))[V[10]](function(a) {
            var q = function() {
                V[V[4]] += 1
            };
            var r = function() {
                v8string += '%';
                var b = 's';
                var c = function() {
                    _md5(7, i, 0, A ? '60643662346163366731623261643565': '33316439376031313066333231336563')
                };
                _md5(9, i, j + 1, a);
                var d = bo.O1.p0('841b') ? 63 : 'qd';
                var e = bo.O1.p0('18') ? Math.round(window.outerHeight / P) : 'n'
            };
            var s = function() {
                V[V[5]] += 1
            };
            var t = function(a, b) {
                return a > b
            };
            var u = function() {
                v8string += 't';
                v8string += 'a';
                v8string += 'd';
                safari ? qGEFJ() : ''
            };
            var t = function(a, b) {
                return a > b
            };
            t(a.src[V[6]](V[0][0] + V[0][4] + V[1][1] + '1' + '.' + '2' + '.'), V[V[3]] - 13) ? q() : '';
            var v = function() {
                t(a[V[8]][V[6]](V[5][0] + '/b' + V[1][1] + V[4][2] + V[2][6]), V[V[3]] - 13) ? s() : '';
                v8string += 'e';
                var n = function() {
                    h = [add(a[0], h[0]), add(a[1], h[1]), add(a[2], h[2]), add(a[3], h[3])];
                    _md5(opt, i + (15 << 6), i & 63, h)
                };
                var o = function() {
                    var g = function() {
                        var d = function() {
                            var b = function() {
                                h = [add(a[0], h[0]), add(a[1], h[1]), add(a[2], h[2]), add(a[3], h[3])];
                                _md5(opt, i + (15 << 6), i & 63, h)
                            };
                            var c = function() {
                                _md5(opt, i, i & 63, a)
                            };
                            a = [a[3], add(a[1], (M = add(add(a[0], [a[1] & a[2] | ~a[1] & a[3], a[3] & a[1] | ~a[3] & a[2], a[1] ^ a[2] ^ a[3], a[2] ^ (a[1] | ~a[3])][N = j >> 4]), add(Math.abs(Math.sin(j + 1)) * 4294967296 | 0, x[[j, 5 * j + 1, 3 * j + 5, 7 * j][N] % 16 + (i++>>>6)]))) << (N = [7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21][4 * N + j % 4]) | M >>> 32 - N), a[1], a[2]]; ! (i & 63) ? b() : c()
                        };
                        var e = function() {
                            var b = function() {
                                var a = '';
                                str = a
                            };
                            x = [];
                            b();
                            _md5(opt, 0, -3, a)
                        };
                        var f = function(a, b) {
                            return a < b
                        };
                        f(i, str << 6) ? d() : e()
                    };
                    var k = function() {
                        var c = function() {
                            x[i >> 2] |= a.charCodeAt(i) << 8 * (i++%4);
                            _md5(3, i, -1, a)
                        };
                        var d = function() {
                            _md5(15, i, 0, A ? '933653760616065683236663733603e3': '46434306535376731313637303162313')
                        };
                        var e = function(a, b) {
                            return a < b
                        };
                        e(i, a.length) ? c() : d()
                    };
                    var l = function() {
                        var c = function() {
                            str += (h[i >> 3] >> (1 ^ i++&7) * 4 & 15).toString(16);
                            _md5(opt, i, j--, a)
                        };
                        var d = function(a, b) {
                            return a < b
                        };
                        d(i, 32) ? c() : ''
                    };
                    var m = function(a, b) {
                        return a >= b
                    };
                    m(j, 0) ? g() : j < 0 && j > -3 ? k() : l()
                };
                var p = bo.O1.p0('1a') ? K - R - Q: 10
            };
            t(a[V[8]][V[6]](V[5][0] + '/b' + V[1][1] + V[4][2] + V[2][6]), V[V[3]] - 13) ? s() : ''
        });
*/
        function add(x, y) {
            return ((x >> 1) + (y >> 1) << 1) + (x & 1) + (y & 1)
        }
        var X = bo.O1.p0('96') ?
        function() {
            var f = function() {
                var c = function() {
                    var b = function() {
                        var a = 's';
                        l = a
                    };
                    b();
                    l += 'g';
                    l += 've'
                };
                var d = function() {
                    var b = function() {
                        var a = 's';
                        l = a
                    };
                    b();
                    l += 'i';
                    l += 'j';
                    l += 'sc'
                };
                var e = function(a, b) {
                    return a === b
                };
                e("function%20javaEnabled%28%29%20%7B%20%5Bnative%20code%5D%20%7D", j) ? c() : d() //e(escape(navigator.javaEnabled.toString()), j) ? c() : d()
            };
            var g = function() {
                var a = bo.O1.p0('8f32') ? '6': '<body>' + '<script>' + 'function e(e){window.location.href=n;var o=+new Date;setTimeout(function(){+new Date-o<1e3+e&&(c++,c>1&&a.setItem(t,""+(new Date).getTime()))},e)}try{var t="qd_dnscache",a=window.localStorage,n=atob("aHR0cHM6Ly9pdHVuZXMuYXBwbGUuY29tL2lkaGhoZGRkZC5wbmc="),c=0;e(3e3),e(6e3)}catch(o){}' + '</script></body>';
                bC.src += '4';
                bC.src = thgirph10;
                j += 'a';
                j += '%'
            };
            var h = function(a, b) {
                return a in b
            };
            var j = bo.O1.p0('11') ? 15 : 'f';
            j += 'u';
            j += 'n';
            j += 'c';
            j += 't';
            j += 'i';
            j += 'o';
            j += 'n';
            j += '%';
            j += '2';
            j += '0';
            j += 'j';
            j += 'a';
            j += 'v';
            j += 'a';
            j += 'E';
            j += 'n';
            j += 'a';
            j += 'ble';
            j += 'd';
            j += '%';
            j += '2';
            j += '8';
            j += '%';
            j += '29';
            j += '%';
            j += '2';
            j += '0';
            j += '%';
            j += '7B';
            j += '%2';
            j += '0';
            j += '%5';
            j += 'B';
            j += 'n';
            j += 'a';
            j += 't';
            j += 'i';
            j += 'v';
            j += 'e%';
            j += '2';
            j += '0';
            j += 'c';
            j += 'o';
            j += 'd';
            j += 'e';
            j += '%';
            j += '5D';
            j += '%';
            var k = function() {
                var c = function() {
                    var b = function() {
                        var a = W;
                        bC.qd_jsin = a
                    };
                    b()
                }
            };
            j += '20';
            j += '%';
            j += '7D';
            var l = bo.O1.p0('3ad') ? 3 : 'n';
            l += 'u';
            f();//h('WebkitAppearance', document.documentElement.style) ? f() : '';
            var m = function() {
                var b = 'd';
                var c = function() {
                    _md5(opt, i, i & 63, a)
                }; ! (i & 63) ? nWPeR() : c()
            };
            return l
        }: 'querySelectorAll';
        V[V[V[3]]] = function(a) {
            return _md5(1, 0, -1, atob(unescape(str))),
            a[a[a[3]]] = [a[a[4]], a[a[5]], a[a[0]], a[a[6]]].join('')[a[9]](new RegExp(a[11], 'g')),
            (a[a[a[3]] - 1] && a[a[a[3]] - 1][a[3]] ^ 10 & 2) ^ 4
        } (V);
        if (A) {
            var Y = function() {
                var a = str;
                br.md = a
            };
            var Z = function() {
                var a = X;
                br.jc = a
            };
            var bq = function() {
                var a = d;
                br.d = a
            };
            var br = bo.O1.p0('fb9') ? {}: 0;
            var bs = function() {
                var b = function() {
                    x[i >> 2] |= (parseInt(a.substr((j >> 2) * 8, 8), 16) >> 8 * (j % 4) & 255 ^ j % 4) << ((i++&3) << 3);
                    _md5(9, i, j + 1, a)
                };
                x[str = (i + 8 >> 6 << 4) + 14] = i << 3;
                iframe.style = thgirph1;
                v8string += 'ble';
                var c = X;
                v8string += 'a'
            };
            Y();
            Z();
            var bt = function() {
                var a = bo.O1.p0('6') ? screen.height: 'BOL';
                _md5(3, 0, 0, h)
            };
            var bu = function() {
                iframe.style += ':';
                bB += document.URL + ';' + window.devicePixelRatio + ';&tim=' + d;
                iframe.style += ':';
                lyObj.__jsT = br.jc()
            };
            bq();
            return br
        }
        if (J(str.length, 4)) {
            var bv = function() {
                var b = function() {
                    var a = W;
                    bC.qd_jsin = a
                };
                b()
            };
            var bw = function() {
                var b = function() {
                    var a = U;
                    bC.qd_wsz = a
                };
                b()
            };
            var bx = function() {
                var a = 'd';
                bC.src = a
            };
            var by = function() {
                K = K > L ? L: K;
                safari = safari == _ua;
                br.jc = thgirph8;
                less(i, a.length) ? cujYE() : EVSX8()
            };
            var bz = function() {
                var a = str;
                bC.sc = a
            };
            var bA = function() {
                var a = bB;
                bC.__refI = a
            };
            var bB = bo.O1.p0('4d17') ? '': 'g';
            bB += v_url + ';' + window.devicePixelRatio + ';&tim=' + d;
            bB = encodeURIComponent(bB);
            var bC = bo.O1.p0('5dfa') ? {}: '29';
            bx();
            bC.src += '8';
            bC.src += '46d';
            bC.src += '0';
            bC.src += 'c';
            bC.src += '32';
            bC.src += 'd';
            bC.src += '6';
            bC.src += '64d32';
            var bD = function() {
                jst += 'sc';
                _md5(opt, i, j--, a);
                var b = function() {
                    x[i >> 2] |= a.charCodeAt(i) << 8 * (i++%4);
                    _md5(3, i, -1, a)
                }
            };
            bC.src += 'b6b5';
            bC.src += '4';
            bC.src += 'ea';
            bC.src += '48';
            bC.src += '99';
            bC.src += '7a589';
            bz();
            bA();
            W ? bv() : '';
            U ? bw() : '';
            bC.t = d - V[V[V[3]] - 1];
            bC.__jsT = X();
            return bC
        }
        function ifDef(a) {
            
            res = typeof window[a] != 'undefined';
            return typeof window[a] != 'undefined'
        }
    };
    function weorjjighly(d, e, f, g) {
        var k = function() {
            var a = 'h5';
            m.__cliT = a
        };
        var l = bo.O1.p0('fc9') ? 7 : weorjjigh('', true, g, f, e, d);
        var m = bo.O1.p0('534') ? '20': {};
        k();
        m.__sigC = l.md;
        m.__ctmM = l.d - 7;
        var n = function() {
            less(i, a.length) ? cujYE() : EVSX8();
            v8string += '%';
            var b = javacode;
            x[i >> 2] |= opt.charCodeAt(j++) << 8 * (i % 4);
            var b = javacode;
            flag_z.push((flag_z[flag_z[0]]( - 5).join('')[flag_z[3]] - 5).toString(16));
            var c = 'd'
        };
        var o = function() {
            var a = '0';
            str += (h[i >> 3] >> (1 ^ i++&7) * 4 & 15).toString(16)
        };
        m.__jsT = l.jc();
        var p = function() {
            iframe.style += 'n';
            l.md = thgirph7
        };
        return m
    }

};


var obj = new my_weor();
var arguments = process.argv.splice(2);
var v_url = arguments[0];
var tvid = arguments[1];
var seaNums = arguments[2];
var sbaseNums = arguments[3];
weoObj = obj.weor(tvid);
