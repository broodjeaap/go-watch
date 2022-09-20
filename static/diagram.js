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
var DiagramNode = /** @class */ (function () {
    function DiagramNode(id, x, y, label, meta) {
        if (meta === void 0) { meta = {}; }
        this.hover = false;
        this.inputHover = false;
        this.outputHover = false;
        this.meta = {};
        this.id = id;
        this.x = x;
        this.y = y;
        this.label = label;
        this.meta = meta;
        this.resize();
    }
    DiagramNode.prototype.resize = function () {
        var textSize = _diagram.ctx.measureText(this.label);
        var height = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.width = textSize.width + height;
        this.height = height;
    };
    DiagramNode.prototype.pointInNode = function (x, y) {
        if (x < this.x) {
            return false;
        }
        if (y < this.y) {
            return false;
        }
        if (x > this.x + this.width) {
            return false;
        }
        if (y > this.y + this.height) {
            return false;
        }
        return true;
    };
    DiagramNode.prototype.pointNearNode = function (x, y) {
        // including the input/output circles
        if (x < this.x - this.height / 3) {
            return false;
        }
        if (y < this.y) {
            return false;
        }
        if (x > this.x + this.width + this.height / 3) {
            return false;
        }
        if (y > this.y + this.height) {
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
}());
var ContextMenuItem = /** @class */ (function () {
    function ContextMenuItem(label, callback) {
        this.x = 0;
        this.y = 0;
        this.hover = false;
        this.callback = function (node) { };
        this.label = label;
        this.callback = callback;
    }
    return ContextMenuItem;
}());
var ContextMenu = /** @class */ (function () {
    function ContextMenu(ctx) {
        this.x = 0;
        this.y = 0;
        this.active = false;
        this.mouseOver = false;
        this.textWidth = 0;
        this.textHeight = 0;
        this.textMargin = 0;
        this.width = 0;
        this.height = 0;
        this.items = new Array();
        this.contextNode = null;
        this.ctx = ctx;
        this.ctx.font = "20px Helvetica";
        var textSize = this.ctx.measureText("SomeLongerText");
        this.textWidth = textSize.width;
        this.textHeight = textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent;
        this.textMargin = this.textWidth / 8;
    }
    ContextMenu.prototype.fitContextMenu = function () {
        var e_1, _a;
        this.width = this.textWidth + this.textMargin * 2;
        this.height = this.textHeight + this.textMargin * (this.items.length + 2);
        var index = 0;
        try {
            for (var _b = __values(this.items), _c = _b.next(); !_c.done; _c = _b.next()) {
                var item = _c.value;
                item.x = this.textMargin;
                item.y = this.textHeight * index + this.textMargin * (index + 2);
                index++;
            }
        }
        catch (e_1_1) { e_1 = { error: e_1_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_1) throw e_1.error; }
        }
    };
    ContextMenu.prototype.pointIn = function (x, y) {
        var e_2, _a;
        if (x < this.x) {
            this.mouseOver = false;
            return false;
        }
        if (y < this.y) {
            this.mouseOver = false;
            return false;
        }
        if (x > this.x + this.width) {
            this.mouseOver = false;
            return false;
        }
        if (y > this.y + this.height) {
            this.mouseOver = false;
            return false;
        }
        try {
            for (var _b = __values(this.items), _c = _b.next(); !_c.done; _c = _b.next()) {
                var item = _c.value;
                if (y >= this.y + item.y - this.textHeight && y <= this.y + item.y + this.textHeight) {
                    item.hover = true;
                }
                else {
                    item.hover = false;
                }
            }
        }
        catch (e_2_1) { e_2 = { error: e_2_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_2) throw e_2.error; }
        }
        this.mouseOver = true;
        return true;
    };
    ContextMenu.prototype.clickOn = function () {
        var e_3, _a;
        if (this.contextNode == null) {
            console.warn("No contextNode");
            return;
        }
        try {
            for (var _b = __values(this.items), _c = _b.next(); !_c.done; _c = _b.next()) {
                var item = _c.value;
                if (item.hover) {
                    item.callback(this.contextNode);
                }
            }
        }
        catch (e_3_1) { e_3 = { error: e_3_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_3) throw e_3.error; }
        }
    };
    ContextMenu.prototype.draw = function () {
        var e_4, _a;
        var cameraX = _diagram.cameraX;
        var cameraY = _diagram.cameraY;
        this.ctx.fillStyle = "lightblue";
        this.ctx.fillRect(this.x + cameraX, this.y + cameraY, this.width, this.height);
        try {
            for (var _b = __values(this.items), _c = _b.next(); !_c.done; _c = _b.next()) {
                var item = _c.value;
                this.ctx.fillStyle = this.mouseOver && item.hover ? "red" : "black";
                this.ctx.fillText(item.label, this.x + item.x + cameraX, this.y + item.y + cameraY);
            }
        }
        catch (e_4_1) { e_4 = { error: e_4_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_4) throw e_4.error; }
        }
    };
    return ContextMenu;
}());
var _diagram;
function diargramOnResize() {
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
function diagramOnWheel(ev) {
    //_diagram.onwheel(ev);
}
function diagramOnContext(ev) {
    ev.preventDefault();
}
var Diagrams = /** @class */ (function () {
    function Diagrams(canvasId, editNodeCallback, deleteNodeCallback) {
        if (editNodeCallback === void 0) { editNodeCallback = function () { }; }
        if (deleteNodeCallback === void 0) { deleteNodeCallback = function () { }; }
        this.nodes = new Map();
        this.connections = new Array();
        this.cameraX = 0; // camera position
        this.cameraY = 0;
        this.mouseX = 0; // mouse position on the canvas
        this.mouseY = 0;
        this.worldX = 0; // relative mouse position
        this.worldY = 0;
        this.panning = false;
        this.nodeDragging = null;
        this.nodeHover = null;
        this.makingConnectionNode = null;
        this.scale = 1.0;
        this.editNodeCallback = function () { };
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
        this.contextMenu = new ContextMenu(this.ctx);
        this.contextMenu.items.push(new ContextMenuItem("Edit", editNodeCallback));
        this.contextMenu.items.push(new ContextMenuItem("Delete", deleteNodeCallback));
        this.contextMenu.fitContextMenu();
        this.ctx.font = "20px Helvetica";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        this.canvas.onwheel = diagramOnWheel;
        this.canvas.oncontextmenu = diagramOnContext;
        window.onresize = diargramOnResize;
        this.editNodeCallback = editNodeCallback;
        this.deleteNodeCallback = deleteNodeCallback;
    }
    Diagrams.prototype.onmousemove = function (ev) {
        var e_5, _a;
        var canvasRect = this.canvas.getBoundingClientRect();
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;
        if (this.contextMenu.active) {
            this.contextMenu.pointIn(this.worldX, this.worldY);
        }
        else if (this.nodeHover != null) {
            this.nodeHover.hover = false;
            this.nodeHover.inputHover = false;
            this.nodeHover.outputHover = false;
            this.nodeHover = null;
        }
        if (this.panning) {
            this.cameraX += ev.movementX;
            this.cameraY += ev.movementY;
        }
        else if (this.nodeDragging != null) {
            // this.nodeDragging.x = worldX - this.nodeDragging.width / 2;
            // this.nodeDragging.y = worldY - this.nodeDragging.height / 2;
            this.nodeDragging.x += ev.movementX;
            this.nodeDragging.y += ev.movementY;
        }
        else {
            try {
                for (var _b = __values(this.nodes), _c = _b.next(); !_c.done; _c = _b.next()) {
                    var _d = __read(_c.value, 2), _ = _d[0], node = _d[1];
                    if (node.pointNearNode(this.worldX, this.worldY)) {
                        if (node.pointInInputCircle(this.worldX, this.worldY)) {
                            node.inputHover = true;
                            this.nodeHover = node;
                            break;
                        }
                        else if (this.makingConnectionNode == null && node.pointInOutputCircle(this.worldX, this.worldY)) {
                            node.outputHover = true;
                            this.nodeHover = node;
                            break;
                        }
                        else if (node.pointInNode(this.worldX, this.worldY)) {
                            node.hover = true;
                            this.nodeHover = node;
                            break;
                        }
                    }
                }
            }
            catch (e_5_1) { e_5 = { error: e_5_1 }; }
            finally {
                try {
                    if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
                }
                finally { if (e_5) throw e_5.error; }
            }
        }
        this.draw();
    };
    Diagrams.prototype.onmousedown = function (ev) {
        var e_6, _a;
        if (ev.button != 0) {
            return;
        }
        //this.contextMenu.active = false;
        var canvasRect = this.canvas.getBoundingClientRect();
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;
        try {
            for (var _b = __values(this.nodes), _c = _b.next(); !_c.done; _c = _b.next()) {
                var _d = __read(_c.value, 2), _ = _d[0], node = _d[1];
                if (node.pointNearNode(this.worldX, this.worldY)) {
                    if (node.pointInInputCircle(this.worldX, this.worldY)) {
                        // no dragging from inputs ?
                    }
                    else if (node.pointInOutputCircle(this.worldX, this.worldY)) {
                        this.makingConnectionNode = node;
                        return;
                    }
                }
                if (node.pointInNode(this.worldX, this.worldY)) {
                    this.nodeDragging = node;
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
        this.panning = true;
    };
    Diagrams.prototype.onmouseup = function (ev) {
        var e_7, _a, e_8, _b, e_9, _c, e_10, _d;
        if (ev.button == 2) {
            try {
                for (var _e = __values(this.nodes), _f = _e.next(); !_f.done; _f = _e.next()) {
                    var _g = __read(_f.value, 2), _ = _g[0], node = _g[1];
                    if (node.pointInNode(this.worldX, this.worldY)) {
                        this.contextMenu.x = this.worldX;
                        this.contextMenu.y = this.worldY;
                        this.contextMenu.active = true;
                        this.contextMenu.contextNode = node;
                        this.draw();
                    }
                }
            }
            catch (e_7_1) { e_7 = { error: e_7_1 }; }
            finally {
                try {
                    if (_f && !_f.done && (_a = _e["return"])) _a.call(_e);
                }
                finally { if (e_7) throw e_7.error; }
            }
        }
        if (ev.button != 0) {
            return;
        }
        this.panning = false;
        this.nodeDragging = null;
        if (this.makingConnectionNode !== null) {
            try {
                for (var _h = __values(this.nodes), _j = _h.next(); !_j.done; _j = _h.next()) {
                    var _k = __read(_j.value, 2), _ = _k[0], node = _k[1];
                    if (node == this.makingConnectionNode) {
                        continue;
                    }
                    if (node.pointInInputCircle(this.worldX, this.worldY)) {
                        var connectionExists = false;
                        try {
                            for (var _l = (e_9 = void 0, __values(this.connections)), _m = _l.next(); !_m.done; _m = _l.next()) {
                                var _o = __read(_m.value, 2), output = _o[0], input = _o[1];
                                if (this.makingConnectionNode.id == output.id && node.id == input.id) {
                                    connectionExists = true;
                                }
                            }
                        }
                        catch (e_9_1) { e_9 = { error: e_9_1 }; }
                        finally {
                            try {
                                if (_m && !_m.done && (_c = _l["return"])) _c.call(_l);
                            }
                            finally { if (e_9) throw e_9.error; }
                        }
                        if (!connectionExists) {
                            this.addConnection(this.makingConnectionNode, node);
                        }
                    }
                }
            }
            catch (e_8_1) { e_8 = { error: e_8_1 }; }
            finally {
                try {
                    if (_j && !_j.done && (_b = _h["return"])) _b.call(_h);
                }
                finally { if (e_8) throw e_8.error; }
            }
            this.makingConnectionNode = null;
        }
        if (this.contextMenu.active) {
            if (this.contextMenu.pointIn(this.worldX, this.worldY)) {
                this.contextMenu.clickOn();
                this.draw();
            }
            this.contextMenu.active = false;
        }
        try {
            for (var _p = __values(this.connections), _q = _p.next(); !_q.done; _q = _p.next()) {
                var _r = __read(_q.value, 2), output = _r[0], input = _r[1];
                var _s = __read(output.getOutputCircleXY(), 2), outputX = _s[0], outputY = _s[1];
                outputX += this.cameraX;
                outputY += this.cameraY;
                var _t = __read(input.getInputCircleXY(), 2), inputX = _t[0], inputY = _t[1];
                inputX += this.cameraX;
                inputY += this.cameraY;
                var dX = Math.abs(outputX - inputX);
                this.ctx.beginPath();
                this.ctx.moveTo(outputX, outputY);
                this.ctx.strokeStyle = "black";
                var cp1x = (outputX + dX / 2);
                var cp1y = outputY;
                var cp2x = (inputX - dX / 2);
                var cp2y = inputY;
                this.ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, inputX, inputY);
                this.ctx.stroke();
                this.ctx.closePath();
                var halfway = getBezierXY(0.5, outputX, outputY, cp1x, cp1y, cp2x, cp2y, inputX, inputY);
                var mouseOnHalfway = Math.pow(this.mouseX - halfway.x, 2) + Math.pow(this.mouseY - halfway.y, 2) <= 10 * 10;
                if (mouseOnHalfway) {
                    this.removeConnection(output, input);
                    break;
                }
            }
        }
        catch (e_10_1) { e_10 = { error: e_10_1 }; }
        finally {
            try {
                if (_q && !_q.done && (_d = _p["return"])) _d.call(_p);
            }
            finally { if (e_10) throw e_10.error; }
        }
        this.draw();
    };
    Diagrams.prototype.onwheel = function (ev) {
        if (ev.deltaY > 0) {
            return;
        }
        this.scale = Math.min(Math.max(this.scale - 0.1, 0.1), 1.0);
        this.ctx.scale(this.scale, this.scale);
    };
    Diagrams.prototype.drawBackground = function () {
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    };
    Diagrams.prototype.draw = function () {
        var e_11, _a, e_12, _b;
        var scale = 1 / this.scale;
        this.ctx.clearRect(0, 0, this.canvas.width * scale, this.canvas.height * scale);
        this.drawBackground();
        var fullCircleRadians = Math.PI + (Math.PI * 3);
        if (this.makingConnectionNode != null) {
            var _c = __read(this.makingConnectionNode.getOutputCircleXY(), 2), circleX = _c[0], circleY = _c[1];
            var dX = Math.abs((circleX + this.cameraX) - this.mouseX);
            this.ctx.beginPath();
            this.ctx.moveTo(circleX + this.cameraX, circleY + this.cameraY);
            this.ctx.strokeStyle = "black";
            var cp1x = (circleX + dX / 2) + this.cameraX;
            var cp1y = circleY + this.cameraY;
            var cp2x = (this.mouseX - dX / 2);
            var cp2y = this.mouseY;
            this.ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, this.mouseX, this.mouseY);
            this.ctx.stroke();
            this.ctx.closePath();
        }
        try {
            for (var _d = __values(this.connections), _e = _d.next(); !_e.done; _e = _d.next()) {
                var _f = __read(_e.value, 2), output = _f[0], input = _f[1];
                var _g = __read(output.getOutputCircleXY(), 2), outputX = _g[0], outputY = _g[1];
                outputX += this.cameraX;
                outputY += this.cameraY;
                var _h = __read(input.getInputCircleXY(), 2), inputX = _h[0], inputY = _h[1];
                inputX += this.cameraX;
                inputY += this.cameraY;
                var dX = Math.abs(outputX - inputX);
                this.ctx.beginPath();
                this.ctx.moveTo(outputX, outputY);
                this.ctx.strokeStyle = "black";
                var cp1x = (outputX + dX / 2);
                var cp1y = outputY;
                var cp2x = (inputX - dX / 2);
                var cp2y = inputY;
                this.ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, inputX, inputY);
                this.ctx.stroke();
                this.ctx.closePath();
                var halfway = getBezierXY(0.5, outputX, outputY, cp1x, cp1y, cp2x, cp2y, inputX, inputY);
                var mouseOnHalfway = Math.pow(this.mouseX - halfway.x, 2) + Math.pow(this.mouseY - halfway.y, 2) <= 10 * 10;
                this.ctx.beginPath();
                this.ctx.strokeStyle = mouseOnHalfway ? "red" : "rgba(200, 200, 200, 0.8)";
                this.ctx.moveTo(halfway.x - 10, halfway.y - 10);
                this.ctx.lineTo(halfway.x + 10, halfway.y + 10);
                this.ctx.moveTo(halfway.x + 10, halfway.y - 10);
                this.ctx.lineTo(halfway.x - 10, halfway.y + 10);
                this.ctx.stroke();
                this.ctx.closePath();
            }
        }
        catch (e_11_1) { e_11 = { error: e_11_1 }; }
        finally {
            try {
                if (_e && !_e.done && (_a = _d["return"])) _a.call(_d);
            }
            finally { if (e_11) throw e_11.error; }
        }
        try {
            for (var _j = __values(this.nodes), _k = _j.next(); !_k.done; _k = _j.next()) {
                var _l = __read(_k.value, 2), _ = _l[0], node = _l[1];
                this.ctx.fillStyle = node.hover ? "#303030" : "#161616";
                this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
                this.ctx.fillStyle = "#D3D3D3";
                this.ctx.font = "20px Helvetica";
                this.ctx.fillText(node.label, node.x + this.cameraX + node.height / 2, node.y + this.cameraY + node.height / 1.5);
                this.ctx.fillStyle = node.inputHover ? "red" : "green";
                this.ctx.beginPath();
                this.ctx.arc(node.x + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
                this.ctx.fill();
                this.ctx.beginPath();
                this.ctx.fillStyle = node.outputHover ? "red" : "green";
                this.ctx.moveTo(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY);
                this.ctx.arc(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
                this.ctx.fill();
            }
        }
        catch (e_12_1) { e_12 = { error: e_12_1 }; }
        finally {
            try {
                if (_k && !_k.done && (_b = _j["return"])) _b.call(_j);
            }
            finally { if (e_12) throw e_12.error; }
        }
        if (this.contextMenu.active) {
            this.contextMenu.draw();
        }
    };
    Diagrams.prototype.addNode = function (id, x, y, label, meta) {
        if (meta === void 0) { meta = {}; }
        this.nodes.set(id, new DiagramNode(id, x, y, label, meta));
    };
    Diagrams.prototype.addConnection = function (A, B) {
        this.connections.push([A, B]);
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
        this.connections.push([A, B]);
    };
    Diagrams.prototype.removeConnection = function (A, B) {
        var e_13, _a;
        var index = 0;
        try {
            for (var _b = __values(this.connections), _c = _b.next(); !_c.done; _c = _b.next()) {
                var _d = __read(_c.value, 2), output = _d[0], input = _d[1];
                if (output.id == A.id && input.id == B.id) {
                    this.connections.splice(index, 1);
                }
                index++;
            }
        }
        catch (e_13_1) { e_13 = { error: e_13_1 }; }
        finally {
            try {
                if (_c && !_c.done && (_a = _b["return"])) _a.call(_b);
            }
            finally { if (e_13) throw e_13.error; }
        }
    };
    Diagrams.prototype.onresize = function () {
        this.fillParent();
    };
    Diagrams.prototype.fillParent = function () {
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        this.draw();
    };
    return Diagrams;
}());
// http://www.independent-software.com/determining-coordinates-on-a-html-canvas-bezier-curve.html
function getBezierXY(t, sx, sy, cp1x, cp1y, cp2x, cp2y, ex, ey) {
    return {
        x: Math.pow(1 - t, 3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x
            + 3 * t * t * (1 - t) * cp2x + t * t * t * ex,
        y: Math.pow(1 - t, 3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y
            + 3 * t * t * (1 - t) * cp2y + t * t * t * ey
    };
}
