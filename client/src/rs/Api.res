module Get = {
  let request = async (url): result<JSON.t, Err.t> => {
    let defaultErr = "Error processing API get-request"

    try {
      let resp = await Fetch.fetch(url)
      let json = await Response.json(resp)

      resp.ok ? Ok(json) : json->Err.parse(~defaultErr)->Error
    } catch {
    | JsExn(exn) => exn->Err.catch(~defaultErr)->Error
    }
  }
}

module Patch = {
  let request = async (url, key, value): result<JSON.t, Err.t> => {
    let defaultErr = "Error processing API patch-request"
    let body = BodyInit.fromString(`{"${key}":"${value}"}`)
    let headersDict = dict{"Content-Type": "application/json"}
    let headers = HeadersInit.fromDict(headersDict)
    let init: Request.requestInit = {method: "PATCH", body, headers}

    try {
      let resp = await Fetch.fetch(url, ~init)
      let json = await Response.json(resp)

      resp.ok ? Ok(json) : json->Err.parse(~defaultErr)->Error
    } catch {
    | JsExn(exn) => exn->Err.catch(~defaultErr)->Error
    }
  }
}
