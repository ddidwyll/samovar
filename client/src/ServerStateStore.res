type patches<'value> = FeedClient.patches<'value>
type decoder<'value> = FeedClient.valueDecoder<'value>
type state<'value> = dict<'value>

type t<'value> = {
  subscribe: (state<'value> => unit) => unit => unit,
}

let make = (~url: string, ~decodeValue: decoder<'value>): t<'value> => {
  let state = ref(Dict.make())
  let connection = ref(None)
  let publishRef = ref(() => ())
  let onError = Some(Console.error)

  let onPatch = (patches: patches<'value>) => {
    Console.log2("onPatch.patches", patches)

    patches->Array.forEach(patch => {
      state.contents->Dict.set(patch.key, patch.value)
    })

    Console.log2("onPatch.state", state.contents)

    publishRef.contents()
  }

  let start = () => {
    let conn = FeedClient.start(~url, ~decodeValue, ~onPatch, ~onError)

    switch connection.contents {
    | Some(_) => ()
    | None => connection := Some(conn)
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
