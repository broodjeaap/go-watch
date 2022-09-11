var DiagramNode = /** @class */ (function () {
    function DiagramNode(x, y, width, height, label) {
        this.hover = false;
        this.inputHover = false;
        this.outputHover = false;
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
        this.ctx.font = "30px Helvetica";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
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
                    else if (node.pointInOutputCircle(this.worldX, this.worldY)) {
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
                    console.log("Making connection");
                }
            }
            this.makingConnectionNode = null;
        }
        this.draw();
    };
    Diagrams.prototype.drawBackground = function () {
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    };
    Diagrams.prototype.draw = function () {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
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
        for (var _i = 0, _b = this.nodes; _i < _b.length; _i++) {
            var node = _b[_i];
            this.ctx.fillStyle = node.hover ? "#303030" : "#161616";
            this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
            this.ctx.fillStyle = "#D3D3D3";
            this.ctx.font = "30px Helvetica";
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
    Diagrams.prototype.addNode = function (x, y, label) {
        var textSize = this.ctx.measureText(label);
        var textHeight = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.nodes.push(new DiagramNode(x, y, textSize.width + textHeight, textHeight, label));
    };
    Diagrams.prototype.addConnection = function (A, B) {
        this.connections.push([A, B]);
    };
    Diagrams.prototype.fillParent = function () {
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        this.draw();
    };
    return Diagrams;
}());