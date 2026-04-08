type subscriber<'a> = 'a => unit
type unsubscribe = unit => unit

type t<'a> = {
  subscribe: subscriber<'a> => unsubscribe,
  update: ('a => 'a) => unit,
  set: 'a => unit,
}

let make = (initial: 'a): t<'a> => {
  let value = ref(initial)
  let subscribers: ref<array<subscriber<'a>>> = ref([])

  let notify = () => {
    subscribers.contents->Belt.Array.forEach(s => s(value.contents))
  }

  let subscribe = (run: subscriber<'a>): unsubscribe => {
    subscribers.contents = Belt.Array.concat(subscribers.contents, [run])

    run(value.contents)

    () => {
      subscribers.contents = subscribers.contents->Belt.Array.keep(s => s !== run)
    }
  }

  let set = (v: 'a) => {
    value.contents = v
    notify()
  }

  let update = (fn: 'a => 'a) => {
    value.contents = fn(value.contents)
    notify()
  }

  {
    subscribe,
    update,
    set,
  }
}
