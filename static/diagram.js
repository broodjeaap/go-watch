var DiagramNode = /** @class */ (function () {
    function DiagramNode(id, x, y, width, height, label) {
        this.hover = false;
        this.inputHover = false;
        this.outputHover = false;
        this.id = id;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.label = label;
    }
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
        return [this.x, this.y + this.height / 3];
    };
    DiagramNode.prototype.getOutputCircleXY = function () {
        return [this.x + this.width, this.y + this.height / 3];
    };
    DiagramNode.prototype.pointInInputCircle = function (x, y) {
        var _a = this.getInputCircleXY(), circleX = _a[0], circleY = _a[1];
        var radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    };
    DiagramNode.prototype.pointInOutputCircle = function (x, y) {
        var _a = this.getOutputCircleXY(), circleX = _a[0], circleY = _a[1];
        var radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    };
    return DiagramNode;
}());
var _diagram;
function diargramOnResize() {
    _diagram.fillParent();
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
var Diagrams = /** @class */ (function () {
    function Diagrams(canvasId) {
        this.nodes = new Array();
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
        this.ctx.font = "20px Helvetica";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        this.canvas.onwheel = diagramOnWheel;
        window.onresize = diargramOnResize;
    }
    Diagrams.prototype.onmousemove = function (ev) {
        var canvasRect = this.canvas.getBoundingClientRect();
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;
        if (this.nodeHover != null) {
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
            for (var _i = 0, _a = this.nodes; _i < _a.length; _i++) {
                var node = _a[_i];
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
        this.draw();
    };
    Diagrams.prototype.onmousedown = function (ev) {
        if (ev.button != 0) {
            return;
        }
        var canvasRect = this.canvas.getBoundingClientRect();
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;
        for (var _i = 0, _a = this.nodes; _i < _a.length; _i++) {
            var node = _a[_i];
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
        this.panning = true;
    };
    Diagrams.prototype.onmouseup = function (ev) {
        this.panning = false;
        this.nodeDragging = null;
        if (this.makingConnectionNode !== null) {
            for (var _i = 0, _a = this.nodes; _i < _a.length; _i++) {
                var node = _a[_i];
                if (node == this.makingConnectionNode) {
                    continue;
                }
                if (node.pointInInputCircle(this.worldX, this.worldY)) {
                    this.addConnection(this.makingConnectionNode, node);
                }
            }
            this.makingConnectionNode = null;
        }
        for (var _b = 0, _c = this.connections; _b < _c.length; _b++) {
            var _d = _c[_b], output = _d[0], input = _d[1];
            var _e = output.getOutputCircleXY(), outputX = _e[0], outputY = _e[1];
            outputX += this.cameraX;
            outputY += this.cameraY;
            var _f = input.getInputCircleXY(), inputX = _f[0], inputY = _f[1];
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
                this.connections.splice(this.connections.indexOf([output, input]), 1);
            }
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
        var scale = 1 / this.scale;
        this.ctx.clearRect(0, 0, this.canvas.width * scale, this.canvas.height * scale);
        this.drawBackground();
        var fullCircleRadians = Math.PI + (Math.PI * 3);
        if (this.makingConnectionNode != null) {
            var _a = this.makingConnectionNode.getOutputCircleXY(), circleX = _a[0], circleY = _a[1];
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
        for (var _i = 0, _b = this.connections; _i < _b.length; _i++) {
            var _c = _b[_i], output = _c[0], input = _c[1];
            var _d = output.getOutputCircleXY(), outputX = _d[0], outputY = _d[1];
            outputX += this.cameraX;
            outputY += this.cameraY;
            var _e = input.getInputCircleXY(), inputX = _e[0], inputY = _e[1];
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
        for (var _f = 0, _g = this.nodes; _f < _g.length; _f++) {
            var node = _g[_f];
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
    };
    Diagrams.prototype.addNode = function (id, x, y, label) {
        var textSize = this.ctx.measureText(label);
        var textHeight = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.nodes.push(new DiagramNode(id, x, y, textSize.width + textHeight, textHeight, label));
    };
    Diagrams.prototype.addConnection = function (A, B) {
        this.connections.push([A, B]);
    };
    Diagrams.prototype.drawDiagramNode = function (node) {
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
