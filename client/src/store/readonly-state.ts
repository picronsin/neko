/** Protect JSON-like state without proxying browser media objects or Dates. */
const views = new WeakMap<object, object>()

export function readonlyState<T>(value: T): T {
  if (!value || typeof value !== 'object') return value
  const prototype = Object.getPrototypeOf(value)
  if (!Array.isArray(value) && prototype !== Object.prototype && prototype !== null) return value
  const cached = views.get(value)
  if (cached) return cached as T
  const reject = () => {
    throw new TypeError('State is read-only; use a store command')
  }
  const view = new Proxy(value, {
    get(target, key, receiver) {
      return readonlyState(Reflect.get(target, key, receiver))
    },
    set: reject,
    deleteProperty: reject,
    defineProperty: reject,
    setPrototypeOf: reject,
    preventExtensions: reject,
  })
  views.set(value, view)
  views.set(view, view)
  return view
}
