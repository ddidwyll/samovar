type subscribe = unit => unit
type subscriber<'state> = 'state => unit
type unsubscribe = unit => unit

type t<'state> = {
  subscribe: subscriber<'state> => unsubscribe,
  publish: unit => unit,
}

type subscription<'state> = {
  id: int,
  run: subscriber<'state>,
}

let make = (~start, ~stop, ~getState): t<'state> => {
  let subscribers = ref([])
  let nextId = ref(0)

  let notify = (state: 'state) => {
    Console.log2("notify", state)

    subscribers.contents->Belt.Array.forEach(sub => sub.run(state))
  }

  let publish = () => getState()->notify

  let countSubscribers = () => Belt.Array.length(subscribers.contents)

  let subscribe = (run: subscriber<'state>) => {
    let id = nextId.contents
    nextId := id + 1

    subscribers := [{id, run}, ...subscribers.contents]

    if countSubscribers() == 1 {
      start()
    }

    run(getState())

    () => {
      subscribers := subscribers.contents->Belt.Array.keep(sub => sub.id != id)

      if countSubscribers() == 0 {
        stop()
      }
    }
  }

  {subscribe, publish}
}
