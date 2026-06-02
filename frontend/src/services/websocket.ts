'use client';

import { io, Socket } from 'socket.io-client';
import { WS_URL } from '@/lib/constants';

let socket: Socket | null = null;
let connecting = false;
let connectAttempts = 0;
const MAX_CONNECT_ATTEMPTS = 3;

type ProgressCallback = (data: {
  percent: number;
  speed: string;
  eta: string;
  status: string;
  filename: string;
}) => void;

type StatusCallback = (status: string) => void;

interface Subscription {
  onProgress: ProgressCallback;
  onStatus: StatusCallback;
}

const subscriptions = new Map<string, Subscription>();

/**
 * Get or initialize the Socket.IO connection.
 * Lazy-initialized on first subscribe call.
 * Uses a health-check before connecting to avoid hanging when no backend is running.
 */
async function getSocket(): Promise<Socket | null> {
  // Already have a connected socket
  if (socket?.connected) return socket;

  // Socket exists but disconnected — return it if under max attempts
  if (socket) {
    if (connectAttempts < MAX_CONNECT_ATTEMPTS) {
      connectAttempts++;
      socket.connect();
      return socket;
    }
    // Max attempts reached — disconnect and return null
    socket.removeAllListeners();
    socket.disconnect();
    socket = null;
    return null;
  }

  // Prevent concurrent init
  if (connecting) return null;
  connecting = true;

  try {
    // Quick health-check — if the WS URL is unreachable, skip Socket.IO entirely
    const wsBase = WS_URL.replace(/\/+$/, '');
    const healthUrl = `${wsBase}/api/health`;
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 3000);
    const res = await fetch(healthUrl, { signal: controller.signal, method: 'HEAD' });
    clearTimeout(timeout);
    if (!res.ok) {
      connecting = false;
      return null;
    }
  } catch {
    // Backend not reachable — skip socket
    connecting = false;
    return null;
  }

  socket = io(`${WS_URL}/progress`, {
    transports: ['websocket', 'polling'],
    autoConnect: true,
    reconnection: true,
    reconnectionAttempts: 3,
    reconnectionDelay: 2000,
    reconnectionDelayMax: 5000,
    timeout: 5000,
  });

  socket.on('connect', () => {
    connectAttempts = 0;
  });

  socket.on('disconnect', (_reason) => {
    // no-op
  });

  socket.on('connect_error', () => {
    connectAttempts++;
    if (connectAttempts >= MAX_CONNECT_ATTEMPTS) {
      socket?.removeAllListeners();
      socket?.disconnect();
      socket = null;
    }
  });

  // Listen for progress events for all subscribed tasks
  socket.on('progress', (data: { taskId: string } & Record<string, unknown>) => {
    const sub = subscriptions.get(data.taskId);
    if (sub) {
      sub.onProgress({
        percent: (data.percent as number) || 0,
        speed: (data.speed as string) || '',
        eta: (data.eta as string) || '',
        status: (data.status as string) || '',
        filename: (data.filename as string) || '',
      });
    }
  });

  // Listen for status change events
  socket.on('status', (data: { taskId: string; status: string }) => {
    const sub = subscriptions.get(data.taskId);
    if (sub) {
      sub.onStatus(data.status);
    }
  });

  connecting = false;
  return socket;
}

/**
 * Subscribe to real-time progress for a task.
 * Returns a no-op unsubscribe if the WebSocket connection is unavailable.
 */
export async function subscribeToProgress(
  taskId: string,
  onProgress: ProgressCallback,
  onStatus: StatusCallback
): Promise<() => void> {
  const s = await getSocket();

  if (!s) {
    // WebSocket unavailable — return no-op unsubscribe
    return () => {};
  }

  subscriptions.set(taskId, { onProgress, onStatus });

  // Tell the server we want updates for this task
  s.emit('subscribe', { taskId });

  // Return unsubscribe function
  return () => {
    subscriptions.delete(taskId);
    s.emit('unsubscribe', { taskId });
  };
}

/**
 * Disconnect the WebSocket entirely.
 */
export function disconnectSocket(): void {
  if (socket) {
    subscriptions.clear();
    socket.disconnect();
    socket = null;
  }
}

/**
 * Check if socket is connected.
 */
export function isConnected(): boolean {
  return socket?.connected ?? false;
}
