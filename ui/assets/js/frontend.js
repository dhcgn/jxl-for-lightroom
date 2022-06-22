// @ts-check
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var _this = this;
var stopGetProgress = false;
function getProgress() {
    return __awaiter(this, void 0, void 0, function () {
        var response, result, error_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (stopGetProgress)
                        return [2 /*return*/];
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 4, , 5]);
                    return [4 /*yield*/, fetch('/progress', {
                            method: 'GET',
                            headers: {
                                Accept: 'application/text'
                            }
                        })];
                case 2:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error("Error! status: ".concat(response.status));
                    }
                    return [4 /*yield*/, response.text()];
                case 3:
                    result = _a.sent();
                    console.log('result is: ', result);
                    return [2 /*return*/, result];
                case 4:
                    error_1 = _a.sent();
                    stopGetProgress = true;
                    if (error_1 instanceof Error) {
                        console.log('error message: ', error_1.message);
                        return [2 /*return*/, error_1.message];
                    }
                    else {
                        console.log('unexpected error: ', error_1);
                        return [2 /*return*/, 'An unexpected error occurred'];
                    }
                    return [3 /*break*/, 5];
                case 5: return [2 /*return*/];
            }
        });
    });
}
var sleep = function (ms) {
    return new Promise(function (resolve) { return setTimeout(resolve, ms); });
};
var app = document.getElementById("progressbarts");
var lastProgress = "-1";
function updateProgress() {
    return __awaiter(this, void 0, void 0, function () {
        var progress, p;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    progress = sleep(100).then(function (n) { return getProgress(); });
                    return [4 /*yield*/, progress];
                case 1:
                    p = _a.sent();
                    if (lastProgress == p)
                        return [2 /*return*/];
                    lastProgress = p;
                    if (app) {
                        app.ariaValueNow = p;
                        app.innerText = p;
                        app.style.width = "".concat(p, "%");
                        console.log('set dom ');
                    }
                    return [2 /*return*/];
            }
        });
    });
}
// for (; ;) {
//     updateProgress()
// }
function updateLog() {
    return __awaiter(this, void 0, void 0, function () {
        var log;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, sleep(100)];
                case 1:
                    _a.sent();
                    log = document.getElementById("logts");
                    if (log) {
                        log.src = log.src;
                    }
                    return [2 /*return*/];
            }
        });
    });
}
(function () { return __awaiter(_this, void 0, void 0, function () {
    var e_1;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                _a.trys.push([0, 6, , 7]);
                _a.label = 1;
            case 1: 
            // await sleep(100);
            return [4 /*yield*/, updateProgress()];
            case 2:
                // await sleep(100);
                _a.sent();
                return [4 /*yield*/, updateLog()];
            case 3:
                _a.sent();
                _a.label = 4;
            case 4: return [3 /*break*/, 1];
            case 5: return [3 /*break*/, 7];
            case 6:
                e_1 = _a.sent();
                return [3 /*break*/, 7];
            case 7: return [2 /*return*/];
        }
    });
}); })();
