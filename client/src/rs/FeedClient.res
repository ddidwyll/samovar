type source
type event

@new external makeEventSource: string => source = "EventSource"
@send external closeEventSource: source => unit = "close"
@send external addEventListener: (source, string, event => unit) => unit = "addEventListener"
@get external eventData: event => string = "data"

type patch<'value> = {key: string, value: 'value}
type patches<'value> = array<patch<'value>>

type json = JSON.t
type valueDecoder<'value> = json => option<'value>
type statePatcher<'value> = patches<'value> => unit
type errorHandler = string => unit

type patchDecoder<'value> = (json, valueDecoder<'value>) => option<patch<'value>>
type patchParser<'value> = (event, valueDecoder<'value>, errorHandler) => option<patches<'value>>

type config<'value> = {
  url: string,
  decodeValue: valueDecoder<'value>,
  onPatch: statePatcher<'value>,
  onError: errorHandler,
}

type t<'value> = {
  source: source,
  _phantom: option<'value>,
}

let decodeKey = JSON.Decode.string

let decodePatch: patchDecoder<'value> = (json, decodeValue) => {
  switch json {
  | JSON.Object(obj) =>
    switch (
      Dict.get(obj, "key")->Option.flatMap(decodeKey),
      Dict.get(obj, "value")->Option.flatMap(decodeValue),
    ) {
    | (Some(key), Some(value)) => Some({key, value})
    | _ => None
    }
  | _ => None
  }
}

let parsePatches: patchParser<'value> = (event, decodeValue, onError) => {
  switch eventData(event)->JSON.parseOrThrow {
  | exception _ =>
    onError("Failed to parse patches")
    None
  | json =>
    switch json {
    | JSON.Object(_) => decodePatch(json, decodeValue)->Option.map(p => [p])
    | JSON.Array(patches) => Array.filterMap(patches, p => decodePatch(p, decodeValue))->Some
    | _ => None
    }
  }
}

let start = (
  ~url: string,
  ~decodeValue: valueDecoder<'value>,
  ~onPatch: statePatcher<'value>,
  ~onError: errorHandler,
): t<'value> => {
  let source = makeEventSource(url)

  addEventListener(source, "change", event => {
    switch parsePatches(event, decodeValue, onError) {
    | Some(patches) => onPatch(patches)
    | None => onError("Unexpected patch format")
    }
  })

  addEventListener(source, "init", event => {
    switch parsePatches(event, decodeValue, onError) {
    | Some(patches) => onPatch(patches)
    | None => onError("Unexpected patch format")
    }
  })

  {source, _phantom: None}
}

let stop = (client: t<'value>) => {
  closeEventSource(client.source)
}
