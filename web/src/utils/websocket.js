class WebSocketClient {
  constructor(url) {
    this.url = url;
    this.ws = null;
    this.handlers = {
      progress: () => {},
      result: () => {},
      complete: () => {},
    };
  }

  connect() {
    this.ws = new WebSocket(this.url);
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.handleMessage(data);
    };

    this.ws.onclose = () => {
      console.log('WebSocket 连接已关闭');
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket 错误:', error);
    };
  }

  handleMessage(data) {
    if (this.handlers[data.type]) {
      this.handlers[data.type](data);
    }
  }

  on(type, handler) {
    this.handlers[type] = handler;
  }

  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    }
  }

  close() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

export const createWebSocket = (url) => new WebSocketClient(url); 