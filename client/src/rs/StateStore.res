let getFieldUrl = "http://localhost:4000/store/fields"
let feedSourceUrl = "http://localhost:4000/store/feed"

type errorHandler = (Err.t) => unit

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
    | JSON.Array(arr) => arr->Array.map(parsePatch)->Option.all
    | json => parsePatch(json)->Option.map(p => [p])
    }
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

  let start = (
    ~onPatch: Values.patcher,
    ~onError: errorHandler,
  ): t => {
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

  let stop = (client: t) => {
    closeEventSource(client.source)
  }
}

module State = {
  type ack = Ready | Failed

  type status = {
    ack: ack,
    err: option<string>,
  }

  type t = {
    fields: Fields.t,
    values: Values.t,
    status: status,
  }

  let onError = (err): t => {
    let err = Err.string(err)->Some
    let status = {err, ack: Failed}

    {status, fields: [], values: Dict.make()}
  }

  let onSuccess = (fields, values): t => {
    let status = {ack: Ready, err: None}

    {fields, status, values}
  }

  let make = async (): t => {
    switch await Fields.fetch() {
    | Error(err) => onError(err)
    | Ok(fields) => onSuccess(fields, Dict.make())
    }
  }
}
