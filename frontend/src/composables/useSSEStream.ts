import { computed, onBeforeUnmount, ref } from 'vue';

export type SSEEventType = 'delta' | 'state' | 'reference' | 'done';

export interface SSEBaseEvent<TType extends string = SSEEventType, TPayload = unknown> {
  type: TType;
  payload: TPayload;
}

export type SSEStreamEvent =
  | SSEBaseEvent<'delta', unknown>
  | SSEBaseEvent<'state', unknown>
  | SSEBaseEvent<'reference', unknown>
  | SSEBaseEvent<'done', unknown>;

export interface UseSSEStreamOptions {
  withCredentials?: boolean;
  onDelta?: (event: SSEBaseEvent<'delta', unknown>) => void;
  onState?: (event: SSEBaseEvent<'state', unknown>) => void;
  onReference?: (event: SSEBaseEvent<'reference', unknown>) => void;
  onDone?: (event: SSEBaseEvent<'done', unknown>) => void;
  onError?: (message: string, event?: Event) => void;
  eventTypes?: SSEEventType[];
}

const DEFAULT_EVENT_TYPES: SSEEventType[] = ['delta', 'state', 'reference', 'done'];

function isSupportedEventType(type: string): type is SSEEventType {
  return DEFAULT_EVENT_TYPES.includes(type as SSEEventType);
}

function normalizeIncomingEvent(type: string, rawData: string): SSEStreamEvent {
  let payload: unknown = rawData;

  if (rawData) {
    try {
      payload = JSON.parse(rawData) as unknown;
    } catch {
      throw new Error(`Failed to parse SSE event \"${type}\" as JSON.`);
    }
  }

  if (!isSupportedEventType(type)) {
    throw new Error(`Unsupported SSE event type: ${type}`);
  }

  return {
    type,
    payload,
  };
}

export function useSSEStream(options: UseSSEStreamOptions = {}) {
  const source = ref<EventSource | null>(null);
  const currentUrl = ref('');
  const isStreaming = ref(false);
  const latestEvent = ref<SSEStreamEvent | null>(null);
  const error = ref<string | null>(null);

  const eventTypes = options.eventTypes?.length ? options.eventTypes : DEFAULT_EVENT_TYPES;

  const clearSource = () => {
    if (!source.value) {
      currentUrl.value = '';
      return;
    }

    source.value.close();
    source.value = null;
    currentUrl.value = '';
  };

  const stop = () => {
    clearSource();
    isStreaming.value = false;
  };

  const dispatchEvent = (event: SSEStreamEvent) => {
    latestEvent.value = event;

    switch (event.type) {
      case 'delta':
        options.onDelta?.(event);
        break;
      case 'state':
        options.onState?.(event);
        break;
      case 'reference':
        options.onReference?.(event);
        break;
      case 'done':
        options.onDone?.(event);
        stop();
        break;
    }
  };

  const bindNamedEvent = (eventType: SSEEventType) => {
    source.value?.addEventListener(eventType, (event) => {
      const messageEvent = event as MessageEvent<string>;

      try {
        error.value = null;
        dispatchEvent(normalizeIncomingEvent(eventType, messageEvent.data));
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to handle SSE event.';
        error.value = message;
        options.onError?.(message, event);
        stop();
      }
    });
  };

  const start = (url: string): boolean => {
    if (source.value && isStreaming.value && currentUrl.value === url) {
      return false;
    }

    stop();
    error.value = null;
    latestEvent.value = null;

    const eventSource = new EventSource(url, {
      withCredentials: options.withCredentials,
    });

    source.value = eventSource;
    currentUrl.value = url;
    isStreaming.value = true;

    eventSource.onopen = () => {
      error.value = null;
      isStreaming.value = true;
    };

    eventSource.onerror = (event) => {
      const message = 'Live stream disconnected. Please retry or stop the current round.';
      error.value = message;
      options.onError?.(message, event);
      stop();
    };

    eventTypes.forEach(bindNamedEvent);
    return true;
  };

  onBeforeUnmount(() => {
    stop();
  });

  return {
    start,
    stop,
    isStreaming: computed(() => isStreaming.value),
    latestEvent: computed(() => latestEvent.value),
    error: computed(() => error.value),
  };
}