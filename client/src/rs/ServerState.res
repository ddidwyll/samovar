let getFieldUrl = "http://localhost:4000/store/fields"
let feedSourceUrl = "http://localhost:4000/store/feed"

type errorHandler = Err.t => unit

module Fields = {
  type field = {
    key: string,
    name: string,
    unit: option<string>,
  }

  type t = array<field>

  let parseField = (json: JSON.t): option<field> => {
    JSON.Decode.object(json)->Option.flatMap(obj => {
      let options = ["key", "name", "unit"]->Array.map(k => {
        Dict.get(obj, k)->Option.flatMap(JSON.Decode.string)
      })

      switch options {
      | [Some(key), Some(name), u] => Some({key, name, unit: u})
      | _ => None
      }
    })
  }

  let parseFields = (json: JSON.t): option<t> => {
    JSON.Decode.array(json)->Option.flatMap(arr => {
      Array.map(arr, parseField)->Option.all
    })
  }

  let fetch = async (): result<t, Err.t> => {
    switch await Api.Get.request(getFieldUrl) {
    | Error(err) => Error(err)
    | Ok(json) =>
      switch parseFields(json) {
      | Some(fields) => Ok(fields)
      | None => Err.make("Failed to parse fields")->Error
      }
    }
  }
}

module Values = {
  type value = string
  type t = dict<value>
  type patch = {key: string, value: string}
  type patches = array<patch>
  type patcher = patches => unit

  let parsePatch = (json): option<patch> => {
    JSON.Decode.object(json)->Option.flatMap(obj => {
      let options = ["key", "value"]->Array.map(k => {
        Dict.get(obj, k)->Option.flatMap(JSON.Decode.string)
      })

      switch options {
      | [Some(key), Some(value)] => Some({key, value})
      | _ => None
      }
    })
  }

  let parsePatches = (jsonString): option<patches> => {
    switch JSON.parseOrThrow(jsonString) {
    | exception _ => None
    // | JSON.Array(arr) => arr->Array.map(parsePatch)->Option.all # null-value-patch error
    | JSON.Array(arr) => arr->Array.filterMap(parsePatch)->Some
    | json => parsePatch(json)->Option.map(p => [p])
    }
  }

  let applyPatches = (values, patches) => {
    Array.forEach(patches, patch => {
      Dict.set(values, patch.key, patch.value)
    })
  }
}

module Feed = {
  type source
  type event

  @new external makeEventSource: string => source = "EventSource"
  @send external closeEventSource: source => unit = "close"
  @send external addEventListener: (source, string, event => unit) => unit = "addEventListener"
  @get external eventData: event => string = "data"

  type t = {source: source}

  let subscribe = (~onPatch: Values.patcher, ~onError: errorHandler): t => {
    let source = makeEventSource(feedSourceUrl)

    addEventListener(source, "change", evt => {
      switch eventData(evt)->Values.parsePatches {
      | Some(patches) => onPatch(patches)
      | None => Err.make("Failed to parse patch")->onError
      }
    })

    addEventListener(source, "init", evt => {
      switch eventData(evt)->Values.parsePatches {
      | Some(patches) => onPatch(patches)
      | None => Err.make("Failed to parse patches")->onError
      }
    })

    {source: source}
  }

  let unsubscribe = (feed: t) => {
    closeEventSource(feed.source)
  }
}

module State = {
  type status = {
    mutable ack: [#ok | #err],
    mutable err: option<string>,
  }

  type t = {
    fields: Fields.t,
    values: Values.t,
    status: status,
  }

  let setStatus = (state, status: [#ok | #err(Err.t)]): t => {
    let (ack, err) = switch status {
    | #ok => (#ok, None)
    | #err(e) => (#err, Err.string(e)->Some)
    }

    state.status.ack = ack
    state.status.err = err
    state
  }

  let patchValues = (state, patches): t => {
    state.values->Values.applyPatches(patches)
    state
  }

  let default = (~fields=[]) => {
    fields,
    values: Dict.make(),
    status: {ack: #ok, err: None},
  }

  let make = async (): t => {
    switch await Fields.fetch() {
    | Error(err) => default()->setStatus(#err(err))
    | Ok(fields) => default(~fields)
    }
  }
}

module Store = SvelteStore.Readable

type state = State.t
type store = Store.t<state>

type t = {store: store}

let make = async (): store => {
  let state = await State.make()
  let store = Store.make(state, (_set, update) => {
    let onPatch = patches => update(State.patchValues(_, patches))
    let onError = err => update(State.setStatus(_, #err(err)))
    let feedClient = Feed.subscribe(~onPatch, ~onError)

    () => Feed.unsubscribe(feedClient)
  })

  store
}

let store = await make()
