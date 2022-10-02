var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var __values = (this && this.__values) || function(o) {
    var s = typeof Symbol === "function" && Symbol.iterator, m = s && o[s], i = 0;
    if (m) return m.call(o);
    if (o && typeof o.length === "number") return {
        next: function () {
            if (o && i >= o.length) o = void 0;
            return { value: o && o[i++], done: !o };
        }
    };
    throw new TypeError(s ? "Object is not iterable." : "Symbol.iterator is not defined.");
};
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
var CanvasObject = /** @class */ (function () {
    function CanvasObject(x, y, width, height) {
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.hover = false;
    }
    CanvasObject.prototype.pointInObject = function (p) {
        if (p.x < this.x) {
            return false;
        }
        if (p.y < this.y) {
            return false;
        }
        if (p.x > this.x + this.width) {
            return false;
        }
        if (p.y > this.y + this.height) {
            return false;
        }
        return true;
    };
    return CanvasObject;
}());
var Button = /** @class */ (function (_super) {
    __extends(Button, _super);
    function Button(x, y, label, ctx, callback, node) {
        var _this = _super.call(this, x, y, 0, 0) || this;
        _this.callback = function () { };
        _this.label = label;
        _this.callback = callback;
        _this.node = node;
        _this.resize(ctx);
        return _this;
    }
    Button.prototype.update = function (ms) {
        this.hover = this.pointInObject(new Point(ms.world.x + ms.offset.x, ms.world.y + ms.offset.y));
        if (ms.click && this.hover) {
            this.callback(this.node);
            ms.click = false;
        }
    };
    Button.prototype.draw = function (ctx, ms) {
        ctx.fillStyle = this.hover ? "black" : "#6B6B6B";
        ctx.font = "15px Helvetica";
        ctx.fillText(this.label, this.x + 3, this.y + this.labelHeight + 3);
    };
    Button.prototype.resize = function (ctx) {
        ctx.font = "15px Helvetica";
        var labelSize = ctx.measureText(this.label);
        this.labelWidth = labelSize.width;
        this.width = this.labelWidth + 6;
        this.labelHeight = labelSize.actualBoundingBoxAscent + labelSize.actualBoundingBoxDescent;
        this.height = this.labelHeight + 6;
    };
    return Button;
}(CanvasObject));
var circleTopRadians = Math.PI / 2;
var circleRightRadians = (Math.PI * 3) / 2;
var circleBottomRadians = Math.PI + (Math.PI * 3);
var circleLeftRadians = Math.PI;
var NodeIO = /** @class */ (function (_super) {
    __extends(NodeIO, _super);
    function NodeIO(node, input) {
        var _this = _super.call(this, 0, 0, 0, 0) || this;
        _this.input = false;
        _this.radius = 15;
        _this.input = input;
        _this.node = node;
        _this.reposition();
        return _this;
    }
    NodeIO.prototype.update = function (ms) {
        if (!ms.draggingConnection && !this.input && this.pointInObject(ms.world) && ms.leftDown) {
            ms.draggingConnection = true;
            _diagram.newConnection = new NewConnection(this.node);
        }
    };
    NodeIO.prototype.draw = function (ctx, ms) {
        ctx.fillStyle = this.input ? "red" : "blue";
        ctx.beginPath();
        ctx.arc(ms.offset.x + this.x, ms.offset.y + this.y, this.radius, circleRightRadians, circleTopRadians, this.input);
        ctx.fill();
    };
    NodeIO.prototype.reposition = function () {
        if (this.input) {
            this.x = this.node.x;
            this.y = this.node.y + this.node.height / 2;
        }
        else {
            this.x = this.node.x + this.node.width;
            this.y = this.node.y + this.node.height / 2;
        }
    };
    NodeIO.prototype.pointInObject = function (p) {
        var inCircle = Math.pow(p.x - this.x, 2) + Math.pow(p.y - this.y, 2) <= this.radius * this.radius;
        if (!inCircle) {
            this.hover = false;
        }
        else {
            this.hover = this.input ? p.x < this.x : p.x > this.x;
        }
        return this.hover;
    };
    return NodeIO;
}(CanvasObject));
var NodeConnection = /** @class */ (function (_super) {
    __extends(NodeConnection, _super);
    function NodeConnection(output, input) {
        var _this = _super.call(this, 0, 0, 0, 0) || this;
        _this.controlPoints = {
            dX: 0,
            outputX: 0,
            outputY: 0,
            inputX: 0,
            inputY: 0,
            cp1x: 0,
            cp1y: 0,
            cp2x: 0,
            cp2y: 0
        };
        _this.halfWayPoint = new Point();
        _this.output = output;
        _this.input = input;
        return _this;
    }
    NodeConnection.prototype.update = function (ms) {
        this.controlPoints.outputX = ms.offset.x + this.output.output.x;
        this.controlPoints.outputY = ms.offset.y + this.output.output.y;
        this.controlPoints.inputX = ms.offset.x + this.input.input.x;
        this.controlPoints.inputY = ms.offset.y + this.input.input.y;
        this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        this.controlPoints.cp1x = (this.controlPoints.outputX + this.controlPoints.dX);
        this.controlPoints.cp1y = this.controlPoints.outputY;
        this.controlPoints.cp2x = (this.controlPoints.inputX - this.controlPoints.dX);
        this.controlPoints.cp2y = this.controlPoints.inputY;
        this.halfWayPoint = getBezierXY(0.5, this.controlPoints.outputX, this.controlPoints.outputY, this.controlPoints.cp1x, this.controlPoints.cp1y, this.controlPoints.cp2x, this.controlPoints.cp2y, this.controlPoints.inputX, this.controlPoints.inputY);
        this.hover = Math.pow(this.halfWayPoint.x - ms.canvas.x, 2) + Math.pow(this.halfWayPoint.y - ms.canvas.y, 2) <= 15 * 15;
        if (this.hover && ms.click) {
            _diagram.removeConnection(this.output, this.input);
            ms.click = false;
        }
    };
    NodeConnection.prototype.draw = function (ctx, ms) {
        ctx.beginPath();
        ctx.moveTo(this.controlPoints.outputX, this.controlPoints.outputY);
        ctx.strokeStyle = "#757575";
        ctx.lineWidth = 5;
        ctx.bezierCurveTo(this.controlPoints.cp1x, this.controlPoints.cp1y, this.controlPoints.cp2x, this.controlPoints.cp2y, this.controlPoints.inputX, this.controlPoints.inputY);
        ctx.stroke();
        ctx.closePath();
        ctx.beginPath();
        ctx.strokeStyle = this.hover ? "red" : "rgba(200, 200, 200, 0.8)";
        ctx.moveTo(this.halfWayPoint.x - 10, this.halfWayPoint.y - 10);
        ctx.lineTo(this.halfWayPoint.x + 10, this.halfWayPoint.y + 10);
        ctx.moveTo(this.halfWayPoint.x + 10, this.halfWayPoint.y - 10);
        ctx.lineTo(this.halfWayPoint.x - 10, this.halfWayPoint.y + 10);
        ctx.stroke();
        ctx.closePath();
    };
    NodeConnection.prototype.reposition = function () {
    };
    return NodeConnection;
}(CanvasObject));
var NewConnection = /** @class */ (function (_super) {
    __extends(NewConnection, _super);
    function NewConnection(output) {
        var _this = _super.call(this, 0, 0, 0, 0) || this;
        _this.controlPoints = {
            dX: 0,
            outputX: 0,
            outputY: 0,
            inputX: 0,
            inputY: 0,
            cp1x: 0,
            cp1y: 0,
            cp2x: 0,
            cp2y: 0
        };
        _this.output = output;
        return _this;
    }
    NewConnection.prototype.update = function (ms) {
        var e_1, _a;
        this.input = null;
        try {
            for (var _b = __values(_diagram.nodes.values()), _c = _b.next(); !_c.done; _c = _b.next()) {
                var node = _c.value;
                if (this.output.id != node.id && node.pointNearNode(ms.world)) {
                    this.input = node;
                }
            }
        }
        catch (e_1_1) { e_1 = { error: e_1_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_1) throw e_1.error; }
        }
        if (this.input == null) {
            this.controlPoints.outputX = ms.offset.x + this.output.output.x;
            this.controlPoints.outputY = ms.offset.y + this.output.output.y;
            this.controlPoints.inputX = ms.offset.x + ms.world.x;
            this.controlPoints.inputY = ms.offset.y + ms.world.y;
            this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        }
        else {
            this.controlPoints.outputX = ms.offset.x + this.output.output.x;
            this.controlPoints.outputY = ms.offset.y + this.output.output.y;
            this.controlPoints.inputX = ms.offset.x + this.input.input.x;
            this.controlPoints.inputY = ms.offset.y + this.input.input.y;
            this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        }
        this.controlPoints.cp1x = (this.controlPoints.outputX + this.controlPoints.dX);
        this.controlPoints.cp1y = this.controlPoints.outputY;
        this.controlPoints.cp2x = (this.controlPoints.inputX - this.controlPoints.dX);
        this.controlPoints.cp2y = this.controlPoints.inputY;
    };
    NewConnection.prototype.draw = function (ctx, ms) {
        ctx.beginPath();
        ctx.moveTo(this.controlPoints.outputX, this.controlPoints.outputY);
        ctx.strokeStyle = "#7575A5";
        ctx.lineWidth = 5;
        ctx.bezierCurveTo(this.controlPoints.cp1x, this.controlPoints.cp1y, this.controlPoints.cp2x, this.controlPoints.cp2y, this.controlPoints.inputX, this.controlPoints.inputY);
        ctx.stroke();
        ctx.closePath();
    };
    NewConnection.prototype.reposition = function () {
    };
    return NewConnection;
}(CanvasObject));
var DiagramNode = /** @class */ (function (_super) {
    __extends(DiagramNode, _super);
    function DiagramNode(id, x, y, label, meta, ctx, results, logs) {
        if (meta === void 0) { meta = {}; }
        if (results === void 0) { results = new Array(); }
        if (logs === void 0) { logs = new Array(); }
        var _this = _super.call(this, x, y, 0, 0) || this;
        _this.dragging = false;
        _this.dragOrigin = new Point();
        _this.meta = {};
        _this.id = id;
        _this.label = label;
        _this.meta = meta;
        _this.fixType();
        _this.resize(ctx);
        _this.deleteButton = new Button(0, 0, "Del", ctx, _diagram.deleteNodeCallback, _this);
        _this.editButton = new Button(0, 0, "Edit", ctx, _diagram.editNodeCallback, _this);
        _this.logButton = new Button(0, 0, "Log", ctx, _diagram.logNodeCallback, _this);
        _this.input = new NodeIO(_this, true);
        _this.output = new NodeIO(_this, false);
        _this.results = results;
        _this.logs = logs;
        return _this;
    }
    DiagramNode.prototype.update = function (ms) {
        if (this.pointNearNode(ms.world)) {
            this.input.update(ms);
            this.output.update(ms);
        }
        this.hover = (!ms.draggingNode || this.dragging) && _super.prototype.pointInObject.call(this, ms.world);
        if (this.hover) {
            this.deleteButton.update(ms);
            this.editButton.update(ms);
            this.logButton.update(ms);
            var onButtons = this.deleteButton.hover || this.editButton.hover || this.logButton.hover;
            if (!this.dragging && ms.leftDown && !ms.draggingNode && !ms.draggingConnection && !onButtons) {
                this.dragging = true;
                ms.draggingNode = true;
                this.dragOrigin.x = this.x - ms.world.x;
                this.dragOrigin.y = this.y - ms.world.y;
            }
        }
        else {
            this.deleteButton.hover = false;
            this.editButton.hover = false;
            this.logButton.hover = false;
        }
        if (!ms.leftDown) {
            this.dragging = false;
            ms.draggingNode = false;
        }
        if (this.dragging) {
            this.x = ms.world.x + this.dragOrigin.x;
            this.y = ms.world.y + this.dragOrigin.y;
            this.input.reposition();
            this.output.reposition();
        }
        this.input.update(ms);
        this.output.update(ms);
    };
    DiagramNode.prototype.draw = function (ctx, ms) {
        ctx.fillStyle = this.hover ? "#DDDDDD" : "#BFBFBF";
        ctx.fillRect(ms.offset.x + this.x, ms.offset.y + this.y, this.width, this.height);
        ctx.font = "20px Helvetica";
        ctx.fillStyle = "black";
        var labelX = ms.offset.x + this.x + this.width / 2 - this.labelWidth / 2;
        var labelY = ms.offset.y + this.y + 3 * 2 + this.labelHeight;
        ctx.fillText(this.label, labelX, labelY);
        ctx.font = "15px Helvetica";
        var typeX = ms.offset.x + this.x + this.width / 2 - this.typeWidth / 2;
        var typeY = ms.offset.y + this.y + 3 * 3 + this.typeHeight + this.labelHeight;
        ctx.fillText(this.type, typeX, typeY);
        var resultCount = "" + this.results.length;
        var resultCountSize = ctx.measureText(resultCount);
        var resultCountWidth = resultCountSize.width;
        var resultCountHeight = resultCountSize.actualBoundingBoxAscent + resultCountSize.actualBoundingBoxDescent;
        var resultCountX = ms.offset.x + this.x + this.width - resultCountWidth - 3 * 3;
        var resultCountY = ms.offset.y + this.y + resultCountHeight + 3 * 3;
        ctx.fillText(resultCount, resultCountX, resultCountY);
        this.deleteButton.x = ms.offset.x + this.x;
        this.deleteButton.y = ms.offset.y + this.y + this.height - this.deleteButton.height;
        this.deleteButton.draw(ctx, ms);
        this.editButton.x = ms.offset.x + this.x + this.width - this.editButton.width;
        this.editButton.y = ms.offset.y + this.y + this.height - this.editButton.height;
        this.editButton.draw(ctx, ms);
        this.logButton.x = ms.offset.x + this.x + this.width / 2 - this.logButton.width / 2;
        this.logButton.y = ms.offset.y + this.y + this.height - this.logButton.height;
        this.logButton.draw(ctx, ms);
        this.input.draw(ctx, ms);
        this.output.draw(ctx, ms);
        if (this.logs.length > 0) {
            ctx.moveTo(this.x + 21, this.y + 6);
            ctx.fillStyle = "orange";
            ctx.beginPath();
            ctx.lineTo(this.x + 23, this.y + 21);
            ctx.lineTo(this.x + 6, this.y + 21);
            ctx.lineTo(this.x + 14, this.y + 6);
            ctx.fill();
        }
        ctx.strokeStyle = "#8E8E8E";
        ctx.lineWidth = 3;
        ctx.strokeRect(ms.offset.x + this.x, ms.offset.y + this.y, this.width, this.height);
    };
    DiagramNode.prototype.fixType = function () {
        // @ts-ignore
        this.type = this.meta.type;
        if (["math", "condition"].indexOf(this.type) >= 0) {
            // @ts-ignore
            this.type = this.meta.var1;
        }
    };
    DiagramNode.prototype.resize = function (ctx) {
        ctx.font = "20px Helvetica";
        var labelSize = ctx.measureText(this.label);
        this.labelWidth = labelSize.width;
        this.labelHeight = labelSize.actualBoundingBoxAscent + labelSize.actualBoundingBoxDescent;
        this.height = 70;
        ctx.font = "15px Helvetica";
        var typeSize = ctx.measureText(this.type);
        this.typeWidth = typeSize.width;
        this.typeHeight = typeSize.actualBoundingBoxAscent + typeSize.actualBoundingBoxDescent;
        this.width = Math.max(150, this.labelWidth, this.typeWidth);
    };
    DiagramNode.prototype.pointInObject = function (p) {
        return this.pointNearNode(p) && (_super.prototype.pointInObject.call(this, p) || this.input.pointInObject(p) || this.output.pointInObject(p));
    };
    DiagramNode.prototype.pointNearNode = function (p) {
        // including the input/output circles
        if (p.x < this.x - this.input.radius) {
            return false;
        }
        if (p.y < this.y) {
            return false;
        }
        if (p.x > this.x + this.width + this.output.radius) {
            return false;
        }
        if (p.y > this.y + this.height) {
            return false;
        }
        return true;
    };
    DiagramNode.prototype.getInputCircleXY = function () {
        return [this.x, this.y + this.height / 2];
    };
    DiagramNode.prototype.getOutputCircleXY = function () {
        return [this.x + this.width, this.y + this.height / 2];
    };
    DiagramNode.prototype.pointInInputCircle = function (x, y) {
        var _a = __read(this.getInputCircleXY(), 2), circleX = _a[0], circleY = _a[1];
        var radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    };
    DiagramNode.prototype.pointInOutputCircle = function (x, y) {
        var _a = __read(this.getOutputCircleXY(), 2), circleX = _a[0], circleY = _a[1];
        var radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    };
    return DiagramNode;
}(CanvasObject));
var _diagram;
function tick() {
    _diagram.tick();
    setTimeout(function () {
        tick(), 1000 / 60;
    });
}
function diagramOnResize() {
    _diagram.onresize();
}
function diagramOnMouseDown(ev) {
    _diagram.onmousedown(ev);
}
function diagramOnMouseUp(ev) {
    _diagram.onmouseup(ev);
}
function diagramOnMouseMove(ev) {
    _diagram.onmousemove(ev);
}
function diagramOnContext(ev) {
    ev.preventDefault();
}
var Point = /** @class */ (function () {
    function Point(x, y) {
        if (x === void 0) { x = 0; }
        if (y === void 0) { y = 0; }
        this.x = 0;
        this.y = 0;
        this.x = x;
        this.y = y;
    }
    return Point;
}());
var MouseState = /** @class */ (function () {
    function MouseState() {
        this.canvas = new Point();
        this.world = new Point();
        this.offset = new Point();
        this.delta = new Point();
        this.leftDown = false;
        this.leftUp = false;
        this.panning = false;
        this.draggingNode = false;
        this.draggingConnection = false;
        this.click = true;
    }
    return MouseState;
}());
var Diagrams = /** @class */ (function () {
    function Diagrams(canvasId, editNodeCallback, deleteNodeCallback, logNodeCallback) {
        if (editNodeCallback === void 0) { editNodeCallback = function () { }; }
        if (deleteNodeCallback === void 0) { deleteNodeCallback = function () { }; }
        if (logNodeCallback === void 0) { logNodeCallback = function () { }; }
        this.shouldTick = true;
        this.nodes = new Map();
        this.connections = new Array();
        this.mouseState = new MouseState();
        this.panning = false;
        this.nodeDragging = null;
        this.nodeHover = null;
        this.newConnection = null;
        this.scale = 1.0;
        this.editNodeCallback = function () { };
        this.logNodeCallback = function () { };
        this.deleteNodeCallback = function () { };
        this.canvas = document.getElementById(canvasId);
        if (this.canvas === null) {
            throw "Could not getElementById " + canvasId;
        }
        var ctx = this.canvas.getContext("2d");
        if (ctx === null) {
            throw "Could not get 2d rendering context";
        }
        _diagram = this;
        this.ctx = ctx;
        this.editNodeCallback = editNodeCallback;
        this.logNodeCallback = logNodeCallback;
        this.deleteNodeCallback = deleteNodeCallback;
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        window.onresize = diagramOnResize;
        tick();
    }
    Diagrams.prototype.tick = function () {
        var e_2, _a, e_3, _b, e_4, _c, e_5, _d;
        this.drawBackground();
        if (this.mouseState.leftUp && !this.mouseState.panning && !this.mouseState.draggingNode && !this.mouseState.draggingConnection) {
            this.mouseState.click = true;
        }
        try {
            for (var _e = __values(this.nodes.values()), _f = _e.next(); !_f.done; _f = _e.next()) {
                var node = _f.value;
                node.update(this.mouseState);
            }
        }
        catch (e_2_1) { e_2 = { error: e_2_1 }; }
        finally {
            try {
                if (_f && !_f.done && (_a = _e["return"])) _a.call(_e);
            }
            finally { if (e_2) throw e_2.error; }
        }
        try {
            for (var _g = __values(this.connections), _h = _g.next(); !_h.done; _h = _g.next()) {
                var connection = _h.value;
                connection.update(this.mouseState);
            }
        }
        catch (e_3_1) { e_3 = { error: e_3_1 }; }
        finally {
            try {
                if (_h && !_h.done && (_b = _g["return"])) _b.call(_g);
            }
            finally { if (e_3) throw e_3.error; }
        }
        if (this.newConnection != null) {
            this.newConnection.update(this.mouseState);
        }
        try {
            for (var _j = __values(this.connections), _k = _j.next(); !_k.done; _k = _j.next()) {
                var connection = _k.value;
                connection.draw(this.ctx, this.mouseState);
            }
        }
        catch (e_4_1) { e_4 = { error: e_4_1 }; }
        finally {
            try {
                if (_k && !_k.done && (_c = _j["return"])) _c.call(_j);
            }
            finally { if (e_4) throw e_4.error; }
        }
        if (this.newConnection != null) {
            this.newConnection.draw(this.ctx, this.mouseState);
        }
        try {
            for (var _l = __values(this.nodes.values()), _m = _l.next(); !_m.done; _m = _l.next()) {
                var node = _m.value;
                node.draw(this.ctx, this.mouseState);
            }
        }
        catch (e_5_1) { e_5 = { error: e_5_1 }; }
        finally {
            try {
                if (_m && !_m.done && (_d = _l["return"])) _d.call(_l);
            }
            finally { if (e_5) throw e_5.error; }
        }
        this.mouseState.leftUp = false;
        this.mouseState.click = false;
    };
    Diagrams.prototype.onmousemove = function (ev) {
        var canvasRect = this.canvas.getBoundingClientRect();
        this.mouseState.canvas.x = ev.x - canvasRect.left;
        this.mouseState.canvas.y = ev.y - canvasRect.top;
        this.mouseState.delta.x = ev.movementX;
        this.mouseState.delta.y = ev.movementY;
        if (this.mouseState.panning) {
            this.mouseState.offset.x += this.mouseState.delta.x;
            this.mouseState.offset.y += this.mouseState.delta.y;
        }
        this.mouseState.world.x = this.mouseState.canvas.x - this.mouseState.offset.x;
        this.mouseState.world.y = this.mouseState.canvas.y - this.mouseState.offset.y;
    };
    Diagrams.prototype.onmousedown = function (ev) {
        var e_6, _a;
        if (ev.button != 0) {
            return;
        }
        this.mouseState.leftDown = true;
        try {
            for (var _b = __values(this.nodes.values()), _c = _b.next(); !_c.done; _c = _b.next()) {
                var object = _c.value;
                if (object.pointInObject(this.mouseState.world)) {
                    return;
                }
            }
        }
        catch (e_6_1) { e_6 = { error: e_6_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_6) throw e_6.error; }
        }
        this.mouseState.panning = true;
    };
    Diagrams.prototype.onmouseup = function (ev) {
        this.mouseState.leftDown = false;
        this.mouseState.panning = false;
        this.mouseState.leftUp = true;
        if (this.newConnection != null) {
            if (this.newConnection.input != null) {
                this.addConnection(this.newConnection.output, this.newConnection.input);
                this.mouseState.draggingConnection = false;
            }
        }
        this.newConnection = null;
    };
    Diagrams.prototype.drawBackground = function () {
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    };
    Diagrams.prototype.draw = function () {
    };
    Diagrams.prototype.addNode = function (id, x, y, label, meta, results, logs) {
        if (meta === void 0) { meta = {}; }
        if (results === void 0) { results = new Array(); }
        if (logs === void 0) { logs = new Array(); }
        var node = new DiagramNode(id, x, y, label, meta, this.ctx, results, logs);
        this.nodes.set(id, node);
    };
    Diagrams.prototype.addConnection = function (A, B) {
        this.connections.push(new NodeConnection(A, B));
    };
    Diagrams.prototype.addConnectionById = function (a, b) {
        var A = this.nodes.get(a);
        if (A === undefined) {
            console.error("No node with ID: " + a);
            return;
        }
        var B = this.nodes.get(b);
        if (B === undefined) {
            console.error("No node with ID: " + b);
            return;
        }
        this.connections.push(new NodeConnection(A, B));
    };
    Diagrams.prototype.removeConnection = function (A, B) {
        var e_7, _a;
        var index = 0;
        try {
            for (var _b = __values(this.connections), _c = _b.next(); !_c.done; _c = _b.next()) {
                var connection = _c.value;
                var output = connection.output;
                var input = connection.input;
                if (output.id == A.id && input.id == B.id) {
                    this.connections.splice(index, 1);
                }
                index++;
            }
        }
        catch (e_7_1) { e_7 = { error: e_7_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_7) throw e_7.error; }
        }
    };
    Diagrams.prototype.onresize = function () {
        this.fillParent();
    };
    Diagrams.prototype.fillParent = function () {
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        //this.draw();
    };
    return Diagrams;
}());
// http://www.independent-software.com/determining-coordinates-on-a-html-canvas-bezier-curve.html
function getBezierXY(t, sx, sy, cp1x, cp1y, cp2x, cp2y, ex, ey) {
    return new Point(Math.pow(1 - t, 3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x
        + 3 * t * t * (1 - t) * cp2x + t * t * t * ex, Math.pow(1 - t, 3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y
        + 3 * t * t * (1 - t) * cp2y + t * t * t * ey);
}
