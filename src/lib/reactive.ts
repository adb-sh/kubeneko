let currentContext: ReactiveContext | null = null; // To hold the current reactive context

class ReactiveContext {
  private dependencies: Set<Signal<any>> = new Set(); // To track dependencies
  public effect: () => void;

  // Register a dependency
  track(signal: Signal<any>) {
    if (currentContext) {
      this.dependencies.add(signal);
      signal.subscribe(currentContext.effect);
    }
  }

  // Notify all dependencies when a signal changes
  notify() {
    this.dependencies.forEach((signal) => {
      signal.notify();
    });
  }
}

export function createSignal<T>(initialValue: T) {
  let _value = initialValue;
  const subscribers = new Set() as Set<() => void>;

  function notify() {
    subscribers.forEach((subscriber) => {
      subscriber();
    });
  }

  return {
    oldValue: null as T | null,
    get value(): T {
      if (currentContext) {
        currentContext.track(this); // Track this signal in the current context
      }
      return _value;
    },
    set value(v: T) {
      this.oldValue = _value;
      _value = v;
      notify();
    },
    subscribe: (subscriber) => {
      subscribers.add(subscriber);
    },
  };
}

class Signal<T> {
  private _value: T;
  public subscribers: Set<() => void> = new Set(); // To hold subscribers (reactive functions)

  constructor(initialValue: T) {
    this._value = initialValue;
  }

  // Getter for the signal value
  get(): T {
    if (currentContext) {
      currentContext.track(this); // Track this signal in the current context
    }
    return this._value;
  }

  // Setter for the signal value
  set(newValue: T) {
    this._value = newValue;
    this.notify(); // Notify subscribers when the value changes
  }

  // Notify all subscribers
  notify() {
    this.subscribers.forEach((subscriber) => subscriber());
  }

  subscribe(subscriber) {
    this.subscribers.add(subscriber);
  }

  get value() {
    if (currentContext) {
      currentContext.track(this); // Track this signal in the current context
    }
    return this._value;
  }
  set value(newValue: T) {
    this._value = newValue;
    this.notify(); // Notify subscribers when the value changes
  }
}

// Create a reactive effect
export function effect(fn: () => void) {
  const context = new ReactiveContext();
  currentContext = context; // Set the current context
  context.effect = fn; // Store the effect function in the context
  fn(); // Run the function to track dependencies
  currentContext = null; // Clear the current context
}
