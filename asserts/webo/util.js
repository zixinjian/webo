/**
 * Created by rick on 15/9/5.
 */


var wbSprintf = function (str) {
    var args = arguments,
        i = 1;

    str = str.replace(/%s/g, function () {
        var arg = args[i++];

        if (typeof arg === 'undefined') {
            return '';
        }
        return arg;
    });
    return str;
};