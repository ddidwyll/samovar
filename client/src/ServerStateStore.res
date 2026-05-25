module Map = Belt.Map.String

type item<'value> = FeedClient.patch<'value>
type state<'value> = Map.t<item<'value>>

type t<'value> = {
  subscribe: (state<'value> => unit) => (unit => unit),
}

let make = (
  ~url: string,
  ~decodeValue: JSON.t => option<'value>,
): t<'value> => {
  let state = ref(Map.empty)
  let notifyRef = ref((_state: state<'value>) => ())
  let connection = ref(None)
  let publish = () => notifyRef.contents(state.contents)

  let onPatch = (patch: item<'value>) => {
    state :=
      state.contents
      -> Map.set(patch.key, patch)

    publish()
  }

  let start = () => {
    switch connection.contents {
    | Some(_) => ()
    | None =>
      connection :=
        FeedClient.start(
          ~url,
          ~decodeValue,
          ~onPatch,
          ~onError=Some(Console.error),
        )
        -> Some
    }
  }

  let stop = () => {
    switch connection.contents {
    | None => ()
    | Some(client) =>
      FeedClient.stop(client)
      connection := None
    }
  }

  let subscribers = ref([])

  let subscribe = run => {
    subscribers := [run, ...subscribers.contents]

    notifyRef := state =>
      subscribers.contents
      -> Belt.Array.forEach(cb => cb(state))

    if Belt.Array.length(subscribers.contents) == 1 {
      start()
    }

    run(state.contents)

    () => {
      subscribers :=
        subscribers.contents
        -> Belt.Array.keep(cb => cb != run)

      if Belt.Array.length(subscribers.contents) == 0 {
        stop()
      }
    }
  }

  {subscribe: subscribe}
}
