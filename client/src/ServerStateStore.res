module Map = Belt.Map.String

type item<'value> = FeedClient.patch<'value>
type state<'value> = Map.t<item<'value>>

type t<'value> = {
  subscribe: (state<'value> => unit) => (unit => unit)
}

let make = (
  ~url: string,
  ~decodeValue: JSON.t => option<'value>,
): t<'value> => {
  let state = ref(Map.empty)
  let connection = ref(None)
  let publishRef = ref(() => ())
  let onError = Some(Console.error)

  let onPatch = (patch: item<'value>) => {
    state :=
      state.contents
      -> Map.set(patch.key, patch)

    publishRef.contents()
  }

  let start = () => {
    switch connection.contents {
    | Some(_) => ()
    | None =>
      connection :=
        FeedClient.start(~url, ~decodeValue, ~onPatch, ~onError)
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

  let getState = () => state.contents
  let store = ReadableStore.make(~start, ~stop, ~getState)

  publishRef := store.publish

  {subscribe: store.subscribe}
}
