type subscriber<'a> = 'a => unit

type unsubscribe = unit => unit

type readable<'a> = {
  subscribe: subscriber<'a> => unsubscribe,
}

type entry<'a> = {
  id: int,
  run: subscriber<'a>,
}

type state = Dict.t<JSON.t>

type patch = {
  key: string,
  value: JSON.t,
}

type eventSource
type messageEvent

@new external makeEventSource: string => eventSource = "EventSource"
@get external eventData: messageEvent => string = "data"
@send
external addEventListener: (eventSource, string, messageEvent => unit) => unit = "addEventListener"
@send external closeEventSource: eventSource => unit = "close"

let notify = (subs: ref<array<entry<'a>>>, value: 'a) => {
  subs.contents->Belt.Array.forEach(entry => entry.run(value))
}

let parsePatch = (string): option<patch> => {
  try {
    switch JSON.parseOrThrow(string)->JSON.Decode.object {
    | None => None
    | Some(obj) =>
      switch (Dict.get(obj, "key"), Dict.get(obj, "value")) {
      | (Some(keyJson), Some(valueJson)) =>
        switch JSON.Decode.string(keyJson) {
        | Some(key) => Some({key, value: valueJson})
        | None => None
        }
      | _ => None
      }
    }
  } catch {
  | _ => None
  }
}

let make = (~url): readable<state> => {
  let state = Dict.make()
  let subscribers = ref([])
  let nextId = ref(0)
  let connection = ref(None)

  let emit = () => notify(subscribers, state)

  let start = () => {
    switch connection.contents {
    | Some(_) => ()
    | None => {
        let es = makeEventSource(url)

        addEventListener(es, "message", evt => {
          switch eventData(evt)->parsePatch {
          | None => ()
          | Some({key, value}) => {
              Dict.set(state, key, value)
              emit()
            }
          }
        })

        connection := Some(es)
      }
    }
  }

  let stop = () => {
    switch connection.contents {
    | None => ()
    | Some(es) => {
        closeEventSource(es)
        connection := None
      }
    }
  }

  let subscribe = (run: subscriber<state>): unsubscribe => {
    let id = nextId.contents
    nextId := id + 1

    subscribers := [{id, run}, ...subscribers.contents]

    if Belt.Array.length(subscribers.contents) == 1 {
      start()
    }

    run(state)

    () => {
      subscribers := subscribers.contents->Belt.Array.keep(entry => entry.id != id)

      if Belt.Array.length(subscribers.contents) == 0 {
        stop()
      }
    }
  }

  {subscribe: subscribe}
}
